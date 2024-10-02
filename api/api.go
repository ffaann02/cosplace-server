package api

import (
	"github.com/ffaann02/cosplace-server/api/routes"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
	}))

	app.Use(logger.New())
	apiGroup := app.Group("/api")

	routes.UserRoutes(apiGroup)
	routes.AuthenRoutes(apiGroup)
}
