package www

import (
	"encoding/json"
	"fmt"
	"image/png"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"

	"get.cutie.cafe/levelup/conf"
	"get.cutie.cafe/levelup/draw"

	"get.cutie.cafe/levelup/db/dynamodb"
	"get.cutie.cafe/levelup/fetch"
	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/gofiber/fiber/v2"
)

var (
	app                *fiber.App
	wasHelpfulRecently map[string]int
)

type post struct {
	AppID     *int
	Recaptcha *string
	Review    *string
}

type wasHelpful struct {
	AppID      *int
	WasHelpful *bool
}

type adminResponse struct {
	UnapprovedApps map[string]*types.App
}

func handleStatus(ctx *fiber.Ctx, code int, message string) {
	ctx.SendStatus(code)
	ctx.SendString(message)
}

// findIP finds the most likely IP of the user, checking X-Forwarded-For first.
func findIP(ctx *fiber.Ctx) string {
	if len(ctx.IPs()) > 0 {
		return ctx.IPs()[0]
	}

	return ctx.IP()
}

// Start the web server.
func Start() {
	ResetHelpfulRateLimit()

	ok, rcSiteKey, _ := InitRecaptcha()
	if !ok {
		panic("No reCAPTCHA site key or server key found. Check your LU_RECAPTCHA_SITE and LU_RECAPTCHA_SERVER environment variables.")
	}

	admin, ok := os.LookupEnv("LU_ADMIN")
	if !ok {
		panic("No admin key found. Check your LU_ADMIN environment variable.")
	}

	app = fiber.New()

	app.Get("/api/suggestions", func(ctx *fiber.Ctx) error {
		bytes, err := json.Marshal(dynamodb.GetFullApps(false))
		if err != nil {
			handleStatus(ctx, 500, "Internal server error")
			return nil
		}

		handleStatus(ctx, 200, string(bytes))

		return nil
	})

	app.Post("/api/suggestions", func(ctx *fiber.Ctx) error {
		input := post{}

		if err := ctx.BodyParser(&input); err != nil {
			handleStatus(ctx, 400, "Bad request: could not parse")
			return nil
		}

		if input.Recaptcha == nil || *input.Recaptcha == "" || input.AppID == nil || *input.AppID < 10 || input.Review == nil || len(*input.Review) < 10 || len(*input.Review) > 300 {
			handleStatus(ctx, 400, "Bad request")
			return nil
		}

		realIP := ctx.IP()
		if len(ctx.IPs()) > 0 {
			realIP = ctx.IPs()[len(ctx.IPs())-1]
		}

		if human, err := recaptcha.Confirm(realIP, *input.Recaptcha); !human || err != nil {
			handleStatus(ctx, 400, "CAPTCHA failed")
			return nil
		}

		appx := dynamodb.GetApp(*input.AppID)

		if appx.AppID != 0 {
			handleStatus(ctx, 400, "App already suggested")
			return nil
		}

		app := types.App{
			AppID:         *input.AppID,
			RecommendedAt: time.Now().Unix(),
			IsPending:     true,
			Review:        *input.Review,
		}

		dynamodb.PutApp(app)

		if err := fetch.PostDiscordPreapprove(app.AppID, app.Review); err != nil {
			util.Warn(fmt.Sprintf("Error: %v", err))
		}

		return nil
	})

	app.Post("/api/helpful", func(ctx *fiber.Ctx) error {
		ip := findIP(ctx)
		if _, ok := wasHelpfulRecently[ip]; !ok {
			wasHelpfulRecently[ip] = 0
		}

		if wasHelpfulRecently[ip] > 4 {
			util.Warn(fmt.Sprintf("%s hit 'was helpful' rate limit", ip))
			handleStatus(ctx, 400, "You've rated too much as (un)helpful recently. Try again later.")
			return nil
		}

		wasHelpfulRecently[ip]++

		input := &wasHelpful{}
		if err := ctx.BodyParser(input); err != nil {
			handleStatus(ctx, 400, "Could not parse request")
			return nil
		}

		if input.AppID == nil || input.WasHelpful == nil {
			handleStatus(ctx, 400, "Not all fields filled")
			return nil
		}

		app := dynamodb.GetApp(*input.AppID)
		if app == nil || app.IsPending {
			handleStatus(ctx, 404, "Target app not found")
			return nil
		}

		if *input.WasHelpful {
			app.HelpfulPositive++
		}

		app.HelpfulTotal++

		dynamodb.PutApp(*app)

		handleStatus(ctx, 200, "Thanks!")
		return nil
	})

	app.Get("/api/admin", func(ctx *fiber.Ctx) error {
		if ctx.Query("key") != admin {
			handleStatus(ctx, 400, "Invalid password")
			return nil
		}

		switch ctx.Query("action") {
		case "approve":
			appid, err := strconv.Atoi(ctx.Query("appid"))

			if err != nil || appid < 10 {
				handleStatus(ctx, 400, "Bad AppID")
				return nil
			}

			app := dynamodb.GetApp(appid)

			if app.AppID == 0 {
				app.AppID = appid
			} else {
				app.IsPending = false
			}

			app.RecommendedAt = time.Now().Unix()

			if err := fetch.AllRegions(app); err != nil {
				handleStatus(ctx, 500, "Could not update app")
				return nil
			}

			// HUGE HACK: if we're in serve-only mode mark the last fetch as ZERO so we can make sure
			// we fetch the right prices when the fetch bot runs. Fetch bot is also configured to
			// post webhooks when last update was 0
			if !conf.Fetch {
				app.LastUpdate = 0
			} else {
				if err := fetch.PostDiscord(app.AppID); err != nil {
					util.Warn(fmt.Sprintf("Error: %v", err))
				}
			}

			if err := dynamodb.PutApp(*app); err != nil {
				handleStatus(ctx, 500, "Could not store app")
				return nil
			}

			break
		case "delete":
			appid, err := strconv.Atoi(ctx.Query("appid"))

			if err != nil || appid < 10 {
				handleStatus(ctx, 400, "Bad AppID")
				return nil
			}

			err = dynamodb.DeleteApp(appid)

			if err != nil {
				handleStatus(ctx, 500, "Could not delete app")
				return nil
			}

			break
		default:
		}

		adminx := &adminResponse{
			UnapprovedApps: dynamodb.GetFullApps(true),
		}

		bytes, err := json.Marshal(adminx)
		if err != nil {
			handleStatus(ctx, 500, "Internal server error")
			return nil
		}

		handleStatus(ctx, 200, string(bytes))

		return nil
	})

	app.Get("/", func(ctx *fiber.Ctx) error {
		byt, err := ioutil.ReadFile("static/index.html")

		if err != nil {
			handleStatus(ctx, 500, "Internal server error")
			return nil
		}

		ctx.Set("Content-Type", "text/html")
		ctx.SendString(strings.Replace(string(byt), "{{recaptcha_site_key}}", *rcSiteKey, -1))

		return nil
	})

	app.Get("/api/image/:app.png", func(ctx *fiber.Ctx) error {
		var (
			appid int
			app   *types.App
			err   error
		)

		if appid, err = strconv.Atoi(ctx.Params("app")); err != nil {
			handleStatus(ctx, 404, "Could not find app")
			return nil
		}

		if app = dynamodb.GetApp(appid); app == nil {
			handleStatus(ctx, 404, "Could not find app")
			return nil
		}

		widget, err := draw.Draw(app)

		if err != nil {
			util.Warn(fmt.Sprintf("Error: %v", err))
			handleStatus(ctx, 500, "Could not render app widget")
			return nil
		}

		//ctx.SendStatus(200)
		ctx.Set("Content-Type", "image/png")

		if err := png.Encode(ctx, widget); err != nil {
			util.Warn(fmt.Sprintf("Error: %v", err))
			handleStatus(ctx, 500, "Could not render app widget")
		}

		return nil
	})

	app.Static("/", "./static")

	if os.Getenv("PORT") != "" {
		util.Info("App server starting on port " + os.Getenv("PORT"))
		app.Listen(":" + os.Getenv("PORT"))
	} else {
		util.Info("App server starting on port 4000")
		app.Listen(":4000")
	}
}

// ResetHelpfulRateLimit resets the "was this helpful?" rate limit
func ResetHelpfulRateLimit() {
	wasHelpfulRecently = make(map[string]int)
}

// InitRecaptcha initializes the reCAPTCHA stuff.
func InitRecaptcha() (bool, *string, *string) {
	siteKey, exists := os.LookupEnv("LU_RECAPTCHA_SITE")
	if siteKey == "" || !exists {
		util.Warn("Could not find site key")
		return false, nil, nil
	}

	serverKey, exists := os.LookupEnv("LU_RECAPTCHA_SERVER")
	if serverKey == "" || !exists {
		util.Warn("Could not find server key")
		return false, nil, nil
	}

	recaptcha.Init(serverKey)

	return true, &siteKey, &serverKey
}
