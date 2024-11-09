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
	return c.JSON(fiber.Map{
		"message": "Update user",
	})
}

func DeleteUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Delete user",
	})
}
