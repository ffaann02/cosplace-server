package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func CustomRoutes(app fiber.Router) {
	custom := app.Group("/custom")

	custom.Get("/", handler.GetCommisions)
	// custom.Post("/", handler.CreateCommision)
	custom.Get("/:id", handler.GetCommision)
}

func ProtectedCustomRoutes(app fiber.Router) {
	custom := app.Group("/custom")
	custom.Get("/", handler.GetCommisions)
	custom.Post("/", handler.CreateCommision)
}
