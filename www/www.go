package www

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"get.cutie.cafe/levelup/db/dynamodb"
	"get.cutie.cafe/levelup/search"
	"get.cutie.cafe/levelup/types"
	"get.cutie.cafe/levelup/util"

	"github.com/dpapathanasiou/go-recaptcha"
	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

type post struct {
	AppID     *int
	Recaptcha *string
}

func handleStatus(ctx *fiber.Ctx, code int, message string) {
	ctx.SendStatus(code)
	ctx.SendString(message)
	return
}

// Start the web server.
func Start() {
	ok, rcSiteKey, _ := InitRecaptcha()
	if !ok {
		panic("No reCAPTCHA site key or server key found. Check your LU_RECAPTCHA_SITE and LU_RECAPTCHA_SERVER environment variables.")
	}

	app = fiber.New()

	app.Get("/api/search", func(ctx *fiber.Ctx) error {
		q := ctx.Query("q")
		if q == "" {
			handleStatus(ctx, 400, "Bad request")
			return nil
		}

		apps, err := search.Query(q)
		if err != nil {
			handleStatus(ctx, 500, "Internal server error")
			return nil
		}

		j, err := json.Marshal(apps)
		if err != nil {
			handleStatus(ctx, 500, "Internal server error")
			return nil
		}

		ctx.Set("Content-Type", "application/json")
		ctx.Send(j)

		return nil
	})

	app.Get("/api/suggestions", func(ctx *fiber.Ctx) error {
		bytes, err := json.Marshal(dynamodb.Cache)
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

		if input.Recaptcha == nil || *input.Recaptcha == "" || input.AppID == nil || *input.AppID < 10 || !search.IsApp(*input.AppID) {
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

		app := &types.App{
			AppID:         *input.AppID,
			RecommendedAt: time.Now().Unix(),
			IsPending:     true,
		}

		dynamodb.PutApp(app)

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

	app.Static("/", "./static")

	if os.Getenv("PORT") != "" {
		util.Info("App server starting on port " + os.Getenv("PORT"))
		app.Listen(":" + os.Getenv("PORT"))
	} else {
		util.Info("App server starting on port 4000")
		app.Listen(":4000")
	}
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