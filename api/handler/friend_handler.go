package handler

import (
	"fmt"

	"github.com/ffaann02/cosplace-server/api/helper"
	"github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
)

func GetFriendRequests(c *fiber.Ctx) error {
	db := config.MysqlDB()

	type RequestBody struct {
		UserID string `json:"user_id"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	fmt.Println(body)

	var friends []m.Friendship
	if err := db.Where("user_id = ? AND status = ?", body.UserID, "request").Find(&friends).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}
	fmt.Println(friends)
	return c.JSON(friends)
}

func SendFriendRequest(c *fiber.Ctx) error {
	db := config.MysqlDB()

	type RequestBody struct {
		UserID         string `json:"user_id"`
		FriendUsername string `json:"friend_username"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	fmt.Println(body)

	// Fetch the user by friend_username to get the friend_user_id
	var friendUser m.User
	if err := db.Where("username = ?", body.FriendUsername).First(&friendUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find friend user",
		})
	}

	newFriendshipID, err := helper.GenerateNewFriendshipID(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate new friendship ID",
		})
	}

	// Create the friendship
	friendship := m.Friendship{
		FriendshipID: newFriendshipID,
		UserID:       body.UserID,
		FriendID:     friendUser.UserID, // Use the friend_user_id
		Status:       "request",
	}

	if err := db.Create(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create friendship",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Friend request sent successfully",
	})
}

func GetSuggestions(c *fiber.Ctx) error {
	db := config.MysqlDB()

	type RequestBody struct {
		UserID string `json:"user_id"`
	}

	var body RequestBody
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	var friends []m.Friendship
	if err := db.Where("user_id = ?", body.UserID).Find(&friends).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}
	fmt.Println(friends)
	return c.JSON(friends)
}
