package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func CommisionRoutes(app fiber.Router) {
	commision := app.Group("/commission")

	commision.Get("/", handler.GetCommisions)
	commision.Get("/:id", handler.GetCommision)
}

func ProtectedCommisionRoutes(app fiber.Router) {
	commision := app.Group("/commission")
	commision.Get("/", handler.GetCommisions)

}
