package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func FriendRoutes(app fiber.Router) {
	// commision := app.Group("/commission")

	// commision.Get("/", handler.GetCommisions)
	// commision.Get("/:id", handler.GetCommision)
}

func ProtectedFriendRoutes(app fiber.Router) {
	commision := app.Group("/friend")
	commision.Get("/requests", handler.GetFriendRequests)
	commision.Get("/suggests", handler.GetSuggestions)
	commision.Post("/add", handler.SendFriendRequest)
}
