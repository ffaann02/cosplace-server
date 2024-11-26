package routes

import (
	"github.com/ffaann02/cosplace-server/api/handler"
	"github.com/gofiber/fiber/v2"
)

func UploadRoutes(app fiber.Router) {
	uploader := app.Group("/upload")
	uploader.Post("/test-s3", handler.TestUploadToAmazonS3)
}

func ProtectedUploadRoutes(app fiber.Router) {
	uploader := app.Group("/upload")
	uploader.Post("/profile-image", handler.UploadProfileImage)
	uploader.Post("/cover-image", handler.UploadCoverImage)
	uploader.Post("/product-image", handler.UploadProductImage)
	uploader.Post("/custom-ref-image", handler.UploadCustomPostReferencesImage)
	uploader.Post("/portfolio-image", handler.UploadPortfolioImage)
}
