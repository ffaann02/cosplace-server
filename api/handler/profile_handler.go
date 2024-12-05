package handler

import (
	"fmt"

	"github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
)

func GetProfile(c *fiber.Ctx) error {
	db := config.MysqlDB()
	userId := c.Params("user_id")
	var profile = m.Profile{UserID: userId}
	if err := db.Where("user_id = ?", userId).First(&profile).Error; err != nil {
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
	fmt.Println("Username:", username)

	// Fetch user by username
	var user m.User
	if err := db.Where("username = ?", username).First(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query user",
		})
	}
	fmt.Print(user)

	// Fetch profile by user_id
	var profile m.Profile
	if err := db.Where("user_id = ?", user.UserID).First(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query profile",
		})
	}

	// Fetch seller by user_id (optional)
	var seller m.Seller
	sellerID := ""
	if err := db.Where("user_id = ?", user.UserID).First(&seller).Error; err == nil {
		sellerID = seller.SellerID
	}

	// Fetch interests associated with the profile
	var profileInterests []m.ProfileInterest
	if err := db.Where("profile_id = ?", profile.ProfileID).Find(&profileInterests).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query profile interests",
		})
	}

	// Extract labels of interests
	interests := []string{}
	for _, interest := range profileInterests {
		interests = append(interests, interest.Label)
	}

	// Return response
	return c.JSON(fiber.Map{
		"profile":   profile,
		"seller_id": sellerID,
		"interests": interests,
	})
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
	if err := db.Where("user_id = ?", requestBody.UserID).First(&profile).Error; err != nil {
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

	fmt.Println(requestBody)

	db := config.MysqlDB()
	// Find the profile with the given user_id
	var profile = m.Profile{UserID: requestBody.UserID}
	if err := db.Where("user_id = ?", requestBody.UserID).First(&profile).Error; err != nil {
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

func AddInterests(c *fiber.Ctx) error {
	// Parse the request body to get the user_id and interests
	var requestBody struct {
		UserID    string   `json:"user_id"`
		Interests []string `json:"interests"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	db := config.MysqlDB()

	// Find the profile with the given user_id
	var profile m.Profile
	if err := db.Where("user_id = ?", requestBody.UserID).First(&profile).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query profile",
		})
	}

	// Insert interests into profile_interests
	for _, interest := range requestBody.Interests {
		profileInterest := m.ProfileInterest{
			ProfileID: profile.ProfileID,
			Label:     interest,
		}

		if err := db.Create(&profileInterest).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to insert interest: " + interest,
			})
		}
	}

	// Respond with a success message
	return c.JSON(fiber.Map{
		"message":   "Interests updated successfully",
		"interests": requestBody.Interests,
	})
}
