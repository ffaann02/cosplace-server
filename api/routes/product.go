package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func ProductRoutes(app fiber.Router) {
	product := app.Group("/product")

	product.Get("/", handler.GetAllProducts)
	product.Get("/id/:product_id", handler.GetProductByID)
	product.Get("/name/:product_name", handler.GetProductByName)
	product.Get("/seller/:seller_id", handler.GetAllProductBySellerID)
}

func ProtectedProductRoutes(app fiber.Router) {
	product := app.Group("/product")

	product.Get("/", handler.GetSellerProducts)
	product.Get("/id/:product_id", handler.GetProductByID)
	product.Get("/name/:product_name", handler.GetProductByName)
	product.Put("/:product_id", handler.UpdateProduct)
	product.Delete("/:product_id", handler.DeleteProduct)

	// profile.Put("/upload/", handler.UpdateUser)
	product.Post("/create", handler.CreateProduct)
	product.Post("/upload-images", handler.UploadProducImages)
	// product.Post("/bio", handler.EditBio)
}
