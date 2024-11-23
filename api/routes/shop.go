package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func ProtectedShopRoutes(app fiber.Router) {
	profile := app.Group("/shop")

	profile.Post("/seller/info", handler.GetSellerInfo)

	profile.Post("/upload-image", handler.UploadShopImage)
	profile.Post("/create-new", handler.CreateNewShop)
}
