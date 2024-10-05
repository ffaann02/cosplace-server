package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func UserRoutes(app fiber.Router) {
	user := app.Group("/user")

	user.Get("/", handler.GetUsers)
	user.Get("/:id", handler.GetUser)
	user.Post("/", handler.CreateUser)
	user.Put("/:id", handler.UpdateUser)
	user.Delete("/:id", handler.DeleteUser)
}

func ProtectedUserRoutes(app fiber.Router) {
	user := app.Group("/user")

	user.Post("/", handler.CreateUser)
	user.Put("/:id", handler.UpdateUser)
	user.Delete("/:id", handler.DeleteUser)
}
