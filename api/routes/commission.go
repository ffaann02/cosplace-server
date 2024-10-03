package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func CommisionRoutes(app fiber.Router) {
	user := app.Group("/commission")

	user.Get("/", handler.GetCommisions)
}
