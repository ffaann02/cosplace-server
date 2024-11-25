package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app fiber.Router) {
	product := app.Group("/product")

	product.Get("/", handler.GetAllProducts)
	// product.Get("/:product_id", handler.GetProfile)
	// product.Get("/all", handler.GetFeedProfile)
}

func ProtectedProductRoutes(app fiber.Router) {
	product := app.Group("/product")

	product.Get("/", handler.GetSellerProducts)
	product.Get("/:product_id", handler.GetProduct)
	product.Put("/:product_id", handler.UpdateProduct)
	product.Delete("/:product_id", handler.DeleteProduct)

	// profile.Put("/upload/", handler.UpdateUser)
	product.Post("/create", handler.CreateProduct)
	product.Post("/upload-images", handler.UploadProducImages)
	// product.Post("/bio", handler.EditBio)
}
