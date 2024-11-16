package handler

import (
	"fmt"

	"github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {
	db := config.MysqlDB()
	userId := c.Query("user_id")
	var profile = m.Profile{UserID: userId}
	if err := db.Find(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}
	fmt.Println(profile)
	return c.JSON(profile)
}

func EditBio(c *fiber.Ctx) error {
	// Parse the request body to get the user_id and new bio
	var requestBody struct {
		UserID string `json:"user_id"`
		Bio    string `json:"bio"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	// Respond with a success message
	return c.JSON(fiber.Map{
		"message": "Bio updated successfully",
	})
}
