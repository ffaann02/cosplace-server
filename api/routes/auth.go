package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func AuthenRoutes(app fiber.Router) {
	authen := app.Group("/auth")

	authen.Post("/register", handler.Register)
	authen.Post("/login", handler.Login)
	authen.Post("/logout", handler.Logout)

	authen.Post("/refresh", handler.Refresh)
}
