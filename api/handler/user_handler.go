package handler

import (
	"github.com/gofiber/fiber/v2"
)

func GetUsers(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Get all users",
	})
}

func GetUser(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Get user",
	})
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
