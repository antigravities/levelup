package www

import (
	"encoding/json"
	"os"

	"get.cutie.cafe/levelup/db/dynamodb"
	"get.cutie.cafe/levelup/search"
	"get.cutie.cafe/levelup/util"

	"github.com/gofiber/fiber/v2"
)

var app *fiber.App

type qt struct {
	Query string `json:"q"`
}

func handleStatus(ctx *fiber.Ctx, code int, message string) {
	ctx.SendStatus(code)
	ctx.SendString(message)
	return
}

// Start the web server.
func Start() {
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

	app.Static("/", "./static")

	if os.Getenv("PORT") != "" {
		util.Info("App server starting on port " + os.Getenv("PORT"))
		app.Listen(":" + os.Getenv("PORT"))
	} else {
		util.Info("App server starting on port 4000")
		app.Listen(":4000")
	}
}
