package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func CheckoutRoutes(app fiber.Router) {
	checkout := app.Group("/checkout")

	// checkout.Get("/", handler.GetCheckout)
	// custom.Post("/", handler.CreateCommision)
	checkout.Get("/:id", handler.GetCommision)
}

func ProtectedCheckoutRoutes(app fiber.Router) {
	checkout := app.Group("/checkout")
	checkout.Get("/", handler.GetCheckout)
	// checkout.Post("/", handler.CreateCommision)
	checkout.Post("/", handler.CreateOrder)
	checkout.Get("/:order_id", handler.GetOrderByOrderID)
	checkout.Get("/all/user/:user_id", handler.GetAllOrdersByUserID)
	checkout.Get("/all/seller/:seller_id", handler.GetAllOrdersBySellerID)
}
