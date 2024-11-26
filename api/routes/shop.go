package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func ShopRoutes(app fiber.Router) {
	profile := app.Group("/shop")

	profile.Get("/:seller_id", handler.GetShopInfo)
	// profile.Get("/seller/info", handler.GetSellerInfo)

	// profile.Post("/upload-image", handler.UploadShopImage)
	// profile.Post("/create-new", handler.CreateNewShop)
}

func ProtectedShopRoutes(app fiber.Router) {
	profile := app.Group("/shop")

	// profile.Get("/seller/info", handler.GetSellerInfo)

	profile.Post("/upload-image", handler.UploadShopImage)
	profile.Post("/create-new", handler.CreateNewShop)
}
