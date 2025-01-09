package main

import (
	"log"
	"runtime/debug"
	"shuttle/databases"
	zerolog "shuttle/logger"
	"shuttle/routes"
	"shuttle/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/spf13/viper"
)

func main() {
	utils.InitFirebase()
	zerolog.InitLogger()

	app := fiber.New()

	app.Use(cors.New())

	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${method} ${path} [${status}] ${latency}\n",
	}))

	app.Use(func(c *fiber.Ctx) error {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("Recovered from panic: %v\n%s", r, debug.Stack())
				log.Printf("Request URL: %s\n", c.OriginalURL())
				c.Status(fiber.StatusInternalServerError).SendString("Internal Server Error, please try again later")
				zerolog.LogInfo("System still continue to run", nil)
			}
		}()
		return c.Next()
	})

	db, err := databases.PostgresConnection()
	if err != nil {
		panic(err)
	}

	routes.Route(app, db)

	if err := app.Listen(viper.GetString("BASE_URL")); err != nil {
        panic(err)
    }
}