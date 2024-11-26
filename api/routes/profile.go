package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func ProfileRoutes(app fiber.Router) {
	profile := app.Group("/profile")

	// profile.Get("/", handler.GetUsers)
	profile.Get("/:user_id", handler.GetProfile)
	profile.Get("/feed/:username", handler.GetFeedProfile)
}

func ProtectedProfileRoutes(app fiber.Router) {
	profile := app.Group("/profile")

	// profile.Put("/upload/", handler.UpdateUser)
	profile.Post("/display-name", handler.EditDisplayName)
	profile.Post("/bio", handler.EditBio)
	profile.Post("/add-interests", handler.AddInterests)
}
