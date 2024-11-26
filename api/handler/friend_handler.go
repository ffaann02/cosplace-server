package handler

import (
	"fmt"

	"github.com/ffaann02/cosplace-server/api/helper"
	"github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
)

func GetFriendList(c *fiber.Ctx) error {
	db := config.MysqlDB()

	// Extract the `user_id` query parameter
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid user_id",
		})
	}

	var friends []struct {
		UserID          string `json:"user_id"`
		Username        string `json:"username"`
		DisplayName     string `json:"display_name"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		ProfileImageUrl string `json:"profile_image_url"`
	}

	if err := db.Table("friendships").
		Select("users.user_id, users.username, users.display_name, users.first_name, users.last_name, profiles.profile_image_url").
		Joins("left join users on users.user_id = friendships.friend_id").
		Joins("left join profiles on profiles.user_id = friendships.friend_id").
		Where("friendships.user_id = ? AND friendships.status = ?", userID, "accepted").
		Find(&friends).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}

	return c.JSON(friends)
}

func GetFriendRequests(c *fiber.Ctx) error {
	db := config.MysqlDB()

	// Extract the `user_id` query parameter
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid user_id",
		})
	}

	var requests []struct {
		UserID          string `json:"user_id"`
		Username        string `json:"username"`
		DisplayName     string `json:"display_name"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		ProfileImageUrl string `json:"profile_image_url"`
	}

	if err := db.Table("friendships").
		Select("users.user_id, users.username, users.display_name, users.first_name, users.last_name, profiles.profile_image_url").
		Joins("left join users on users.user_id = friendships.user_id").
		Joins("left join profiles on profiles.user_id = users.user_id").
		Where("friendships.friend_id = ? AND friendships.status = ?", userID, "request").
		Find(&requests).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}

	return c.JSON(requests)
}

func GetFriendWaitingAccept(c *fiber.Ctx) error {
	db := config.MysqlDB()

	// Extract the `user_id` query parameter
	userID := c.Query("user_id")
	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing or invalid user_id",
		})
	}

	var requests []struct {
		UserID          string `json:"user_id"`
		Username        string `json:"username"`
		DisplayName     string `json:"display_name"`
		FirstName       string `json:"first_name"`
		LastName        string `json:"last_name"`
		ProfileImageUrl string `json:"profile_image_url"`
	}

	if err := db.Table("friendships").
		Select("users.user_id, users.username, users.display_name, users.first_name, users.last_name, profiles.profile_image_url").
		Joins("left join users on users.user_id = friendships.friend_id").
		Joins("left join profiles on profiles.user_id = friendships.friend_id").
		Where("friendships.user_id = ? AND friendships.status = ?", userID, "request").
		Find(&requests).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}

	return c.JSON(requests)
}

func GetSuggestions(c *fiber.Ctx) error {
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

	fmt.Println(friendship)

	if err := db.Create(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create friendship",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Friend request sent successfully",
	})
}

func AcceptFriendRequest(c *fiber.Ctx) error {
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

	// Create the friendship
	friendship := m.Friendship{}

	if err := db.Where("user_id = ? AND friend_id = ? AND status = ?", friendUser.UserID, body.UserID, "request").First(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find friend request",
		})
	}

	friendship.Status = "accepted"

	tx := db.Begin()
	if err := tx.Save(&friendship).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to accept friend request",
		})
	}

	// Create the reverse friendship
	newFriendshipID, err := helper.GenerateNewFriendshipID(tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate new friendship ID",
		})
	}

	reverseFriendship := m.Friendship{
		FriendshipID: newFriendshipID,
		UserID:       body.UserID,
		FriendID:     friendUser.UserID, // Use the friend_user_id
		Status:       "accepted",
	}

	if err := tx.Create(&reverseFriendship).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create reverse friendship",
		})
	}

	tx.Commit()

	return c.JSON(fiber.Map{
		"message": "Friend request sent successfully",
	})
}

func RejectFriendRequest(c *fiber.Ctx) error {
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

	// Create the friendship
	friendship := m.Friendship{}

	if err := db.Where("user_id = ? AND friend_id = ? AND status = ?", friendUser.UserID, body.UserID, "request").First(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find friend request",
		})
	}

	tx := db.Begin()
	if err := tx.Delete(&friendship).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to reject friend request",
		})
	}
	tx.Commit()

	return c.JSON(fiber.Map{
		"message": "Friend request rejected successfully",
	})
}

func CancelFriendRequest(c *fiber.Ctx) error {
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

	// Create the friendship
	friendship := m.Friendship{}

	if err := db.Where("user_id = ? AND friend_id = ? AND status = ?", body.UserID, friendUser.UserID, "request").First(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find friend request",
		})
	}

	tx := db.Begin()
	if err := tx.Delete(&friendship).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to cancel friend request",
		})
	}
	tx.Commit()

	return c.JSON(fiber.Map{
		"message": "Friend request canceled successfully",
	})
}

func DeleteFriend(c *fiber.Ctx) error {
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

	// Create the friendship
	friendship := m.Friendship{}

	if err := db.Where("user_id = ? AND friend_id = ? AND status = ?", body.UserID, friendUser.UserID, "accepted").First(&friendship).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find friend request",
		})
	}

	tx := db.Begin()
	if err := tx.Delete(&friendship).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete friend",
		})
	}

	// Create the reverse friendship
	reverseFriendship := m.Friendship{}

	if err := tx.Where("user_id = ? AND friend_id = ? AND status = ?", friendUser.UserID, body.UserID, "accepted").First(&reverseFriendship).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to find reverse friend request",
		})
	}

	if err := tx.Delete(&reverseFriendship).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete reverse friend",
		})
	}

	tx.Commit()

	return c.JSON(fiber.Map{
		"message": "Friend deleted successfully",
	})
}
