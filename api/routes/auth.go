package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func AuthenRoutes(app fiber.Router) {
	authen := app.Group("/auth")

	authen.Get("/check", handler.CheckAuth)
	authen.Post("/refresh", handler.Refresh)
	authen.Post("/register", handler.Register)
	authen.Post("/login", handler.Login)
	// authen.Post("/logout", handler.Logout)
}

func ProtectedAuthenRoutes(app fiber.Router) {
	authen := app.Group("/auth")
	authen.Put("/change-password", handler.ChangePassword)
	authen.Post("/logout", handler.Logout)
}
