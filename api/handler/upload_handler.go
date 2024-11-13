// handler/handler.go
package handler

import (
	"github.com/ffaann02/cosplace-server/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func UploadImage(c *fiber.Ctx) error {
	// Parse the request body to get the base64-encoded image string and user_id
	var requestBody struct {
		Image  string `json:"image"`
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Call the utility function to upload the image with the user_id
	imageURL, err := utils.UploadImageToImgBB(requestBody.UserID, requestBody.Image)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload image",
		})
	}

	// Respond with the URL of the uploaded image
	return c.JSON(fiber.Map{
		"message":   "Image uploaded successfully",
		"image_url": imageURL,
	})
}
