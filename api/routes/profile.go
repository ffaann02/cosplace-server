package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func ProfileRoutes(app fiber.Router) {
	profile := app.Group("/profile")

	profile.Get("/", handler.GetUsers)
	profile.Get("/:user_id", handler.GetProfile)
}

func ProtectedProfileRoutes(app fiber.Router) {
	// profile := app.Group("/profile")

	// profile.Put("/upload/", handler.UpdateUser)
}
