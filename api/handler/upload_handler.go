// handler/handler.go
package handler

import (
	"fmt"

	"github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/ffaann02/cosplace-server/internal/utils"
	"github.com/gofiber/fiber/v2"
)

func UploadProfileImage(c *fiber.Ctx) error {
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
	fmt.Println(imageURL)

	db := config.MysqlDB()
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to begin transaction",
		})
	}

	// Update the profile_image_url in the users table
	if err := tx.Model(&m.Profile{}).Where("user_id = ?", requestBody.UserID).Update("profile_image_url", imageURL).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile image URL",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Respond with the URL of the uploaded image
	return c.JSON(fiber.Map{
		"message":   "Image uploaded successfully",
		"image_url": imageURL,
	})
}

func UploadCoverImage(c *fiber.Ctx) error {
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

	db := config.MysqlDB()
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to begin transaction",
		})
	}

	// Update the profile_image_url in the users table
	if err := tx.Model(&m.Profile{}).Where("user_id = ?", requestBody.UserID).Update("cover_image_url", imageURL).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update profile image URL",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Respond with the URL of the uploaded image
	return c.JSON(fiber.Map{
		"message":   "Image uploaded successfully",
		"image_url": imageURL,
	})
}

func UploadShopImage(c *fiber.Ctx) error {
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

func UploadProductImage(c *fiber.Ctx) error {
	// Parse the request body to get the base64-encoded image string and product_id
	var productImage m.ProductImage
	if err := c.BodyParser(&productImage); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Call the utility function to upload the image with the product_id
	imageURL, err := utils.UploadImageToImgBB(productImage.ProductID, productImage.ImageURL)
	if err != nil {
		fmt.Println(err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload image",
		})
	}

	db := config.MysqlDB()
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to begin transaction",
		})
	}

	// Create a new entry in the product_images table
	newProductImage := m.ProductImage{
		ProductID: productImage.ProductID,
		ImageURL:  imageURL,
	}
	if err := tx.Create(&newProductImage).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create product image",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	// Respond with the URL of the uploaded image
	return c.JSON(fiber.Map{
		"message":   "Image uploaded successfully",
		"image_url": imageURL,
	})
}

func TestUploadToAmazonS3(c *fiber.Ctx) error {
	// Parse the request body to get the base64-encoded image string and user_id
	var requestBody struct {
		UserID string `json:"user_id"`
		Prefix string `json:"prefix"`
		Image  string `json:"image"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	imageURL, err := utils.UploadImageToAmazonS3(requestBody.Image, "test", requestBody.UserID)
	if err != nil {
		fmt.Println("Error uploading image:", err)
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
