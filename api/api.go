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
		AllowOrigins:     "http://localhost:3000, https://cosplace-frontend.pages.dev",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowCredentials: true,
	}))

	app.Use(logger.New())
	apiGroup := app.Group("/api")
	routes.IndexRoutes(apiGroup)
	routes.UserRoutes(apiGroup)
	routes.AuthenRoutes(apiGroup)
	routes.CommisionRoutes(apiGroup)
	routes.ProfileRoutes(apiGroup)
	// routes.UploadRoutes(apiGroup)
	routes.ProductRoutes(apiGroup)

	protectedGroup := apiGroup.Group("/protected")
	protectedGroup.Use(middleware.JWTProtected())
	routes.ProtectedUserRoutes(protectedGroup)
	routes.ProtectedAuthenRoutes(protectedGroup)
	routes.ProtectedCommisionRoutes(protectedGroup)
	routes.ProtectedUploadRoutes(protectedGroup)
	routes.ProtectedProfileRoutes(protectedGroup)
	routes.ProtectedFriendRoutes(protectedGroup)
	routes.ProtectedProductRoutes(protectedGroup)
}
