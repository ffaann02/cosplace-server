package api

import (
	"github.com/ffaann02/cosplace-server/api/routes"
	"github.com/ffaann02/cosplace-server/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func SetupRoutes(app *fiber.App) {

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "http://localhost:3000",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	apiGroup := app.Group("/api")

	routes.UserRoutes(apiGroup)
	routes.AuthenRoutes(apiGroup)

	protectedGroup := apiGroup.Group("/protected")
	protectedGroup.Use(middleware.JWTProtected())

	routes.CommisionRoutes(protectedGroup)
}
