package api

import (
	"github.com/ffaann02/cosplace-server/api/routes"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	routes.UserRoutes(app)
}
