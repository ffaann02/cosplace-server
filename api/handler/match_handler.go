package handler

import (
	"github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
)

func GetCosplayerList(c *fiber.Ctx) error {
	db := config.MysqlDB()

	// Extract the `user_id` query parameter
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid user_id",
		})
	}

	var users []struct {
		UserID          string `json:"user_id"`
		Username        string `json:"username"`
		DisplayName     string `json:"display_name"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		ProfileImageUrl string `json:"profile_image_url"`
	}

	// Exclude users who are already friends, and users who sent incoming friend requests
	subQuery1 := db.Table("friendships").
		Select("user_id").
		Where("friend_id = ?", userID)

	// Exclude users who are already friends, and users who waiting for friend requests to be accepted
	subQuery2 := db.Table("friendships").
		Select("friend_id").
		Where("user_id = ?", userID)

	if err := db.Model(&m.User{}).
		Select("users.user_id, users.username, users.display_name, users.first_name, users.last_name, profiles.profile_image_url").
		Joins("left join profiles on profiles.user_id = users.user_id").
		Where("users.user_id != ?", userID).
		Where("users.user_id NOT IN (?) AND users.user_id NOT IN (?)", subQuery1, subQuery2).
		Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}

	return c.JSON(users)
}
