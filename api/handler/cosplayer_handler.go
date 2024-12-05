package handler

import (
	"fmt"

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
		UserID          string        `json:"user_id"`
		Username        string        `json:"username"`
		DisplayName     string        `json:"display_name"`
		FirstName       string        `json:"first_name"`
		LastName        string        `json:"last_name"`
		ProfileImageUrl string        `json:"profile_image_url"`
		SellerID        string        `json:"seller_id,omitempty"`
		Interests       []*string     `json:"interests"`
		Portfolios      []m.Portfolio `json:"portfolios"`
	}

	// Exclude users who are already friends, and users who sent incoming friend requests
	subQuery1 := db.Table("friendships").
		Select("user_id").
		Where("friend_id = ?", userID)

	// Exclude users who are already friends, and users who are waiting for friend requests to be accepted
	subQuery2 := db.Table("friendships").
		Select("friend_id").
		Where("user_id = ?", userID)

	if err := db.Model(&m.User{}).
		Select("users.user_id, users.username, users.first_name, users.last_name, profiles.profile_image_url, profiles.display_name").
		Joins("left join profiles on profiles.user_id = users.user_id").
		Where("users.user_id != ?", userID).
		Where("users.user_id NOT IN (?) AND users.user_id NOT IN (?)", subQuery1, subQuery2).
		Find(&users).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}

	for i, user := range users {

		var profile m.Profile
		if err := db.Table("profiles").Where("user_id = ?", user.UserID).First(&profile).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to query profiles table",
			})
		}
		profileID := profile.ProfileID
		// Fetch interests associated with the profile
		var profileInterests []m.ProfileInterest
		if err := db.Where("profile_id = ?", profileID).Find(&profileInterests).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to query profile interests",
			})
		}
		fmt.Println(profileInterests)
		var interests []*string
		for _, interest := range profileInterests {
			interestStr := interest.Label
			interests = append(interests, &interestStr)
		}
		users[i].Interests = interests

		// Fetch portfolios associated with the user
		var portfolios []m.Portfolio
		if err := db.Preload("PortfolioImages").Where("created_by = ?", user.UserID).Find(&portfolios).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to query portfolio",
			})
		}
		users[i].Portfolios = portfolios
	}

	return c.JSON(users)
}
