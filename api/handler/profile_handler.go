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

func GetFeedProfile(c *fiber.Ctx) error {
	fmt.Println("GetFeedProfile")
	db := config.MysqlDB()
	username := c.Params("username")

	var user m.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query user",
		})
	}
	var profile m.Profile
	if err := db.Where("user_id = ?", user.UserID).First(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query profile",
		})
	}

	var seller m.Seller
	sellerID := ""
	if err := db.Where("user_id = ?", user.UserID).First(&seller).Error; err == nil {
		sellerID = seller.SellerID
	}

	profileResponse := m.ProfileResponse{
		ProfileID:       profile.ProfileID,
		UserID:          profile.UserID,
		SellerID:        sellerID,
		DisplayName:     profile.DisplayName,
		ProfileImageURL: profile.ProfileImageURL,
		CoverImageURL:   profile.CoverImageURL,
		Bio:             profile.Bio,
		InstagramURL:    profile.InstagramURL,
		TwitterURL:      profile.TwitterURL,
		FacebookURL:     profile.FacebookURL,
		CreatedAt:       profile.CreatedAt,
		UpdatedAt:       profile.UpdatedAt,
		Username:        user.Username,
		Gender:          user.Gender,
		DateOfBirth:     user.DateOfBirth,
	}

	return c.JSON(profileResponse)
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
	fmt.Println(requestBody)
	db := config.MysqlDB()
	// Find the profile with the given user_id
	var profile = m.Profile{UserID: requestBody.UserID}
	if err :=
		db.Find(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}

	// Update the bio
	profile.Bio = requestBody.Bio
	if err := db.Save(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update bio",
		})
	}

	// Respond with a success message
	return c.JSON(fiber.Map{
		"message": "Bio updated successfully",
		"bio":     profile.Bio,
	})
}

func EditDisplayName(c *fiber.Ctx) error {
	// Parse the request body to get the user_id and new bio
	var requestBody struct {
		UserID      string `json:"user_id"`
		DisplayName string `json:"display_name"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	db := config.MysqlDB()
	// Find the profile with the given user_id
	var profile = m.Profile{UserID: requestBody.UserID}
	if err :=
		db.Find(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}

	// Update the display name
	profile.DisplayName = requestBody.DisplayName
	if err := db.Save(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update display name",
		})
	}

	// Respond with a success message
	return c.JSON(fiber.Map{
		"message":      "Display name  updated successfully",
		"display_name": profile.DisplayName,
	})
}
