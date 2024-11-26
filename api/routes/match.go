package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func ProtectedMatchRoutes(app fiber.Router) {
	commision := app.Group("/match-cosplayer")
	// Fetch
	commision.Get("/list", handler.GetCosplayerList)
}
