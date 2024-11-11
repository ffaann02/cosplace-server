package handler

import (
	"fmt"

	"github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	db := config.MysqlDB()
	username := c.Query("username")
	var user = m.User{Username: username}
	if err := db.Find(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}
	fmt.Println(user)
	return c.JSON(user)
}

func GetUser(c *fiber.Ctx) error {
	db := config.MysqlDB()
	username := c.Query("username")
	var user = m.User{Username: username}
	if err := db.Find(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query database",
		})
	}
	fmt.Println(user)
	return c.JSON(user)
}

func CreateUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Create user",
	})
}

func UpdateUser(c *fiber.Ctx) error {
	db := config.MysqlDB()

	var user m.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	username := user.Username
	fmt.Print(username)

	// Begin transaction
	tx := db.Begin()
	if err := tx.Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to begin transaction",
		})
	}

	// Update user
	if err := tx.Model(&m.User{}).Where("username = ?", username).Updates(user).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update user",
		})
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "User updated successfully",
		"user":    user,
	})
}

func DeleteUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Delete user",
	})
}
