package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func PortfolioRoutes(app fiber.Router) {
	portfolio := app.Group("/portfolio")

	portfolio.Get("/:user_id", handler.GetPortfolioByUserID)
	// commision.Get("/:id", handler.GetCommision)
}

func ProtectedPortfolioRoutes(app fiber.Router) {
	portfolio := app.Group("/portfolio")
	// portfolio.Get("/requests", handler.GetFriendRequests)
	// portfolio.Get("/suggests", handler.GetSuggestions)
	portfolio.Post("/create", handler.CreatePortfolio)
}
