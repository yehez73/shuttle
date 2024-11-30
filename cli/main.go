package main

import (
	"shuttle/databases"
	"shuttle/routes"
	"shuttle/utils"
	zerolog "shuttle/logger"

	"github.com/gofiber/contrib/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

func main() {
	zerolog.InitLogger()

	app := fiber.New()

	app.Use(cors.New())

	app.Use("/ws", func(c *fiber.Ctx) error {
		if websocket.IsWebSocketUpgrade(c) {
			c.Locals("allowed", true)
			return c.Next()
		}
		return fiber.ErrUpgradeRequired
	})
	
	app.Get("/ws/:id", websocket.New(utils.HandleWebSocketConnection))

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} ${path} [${status}] ${latency}\n",
	}))

	routes.Route(app)
	database.MongoConnection()

	if err := app.Listen(viper.GetString("BASE_URL")); err != nil {
        panic(err)
    }
}