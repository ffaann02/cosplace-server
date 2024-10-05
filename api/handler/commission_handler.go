package handler

import (
	"github.com/gofiber/fiber/v2"
)

func GetCommisions(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Get all commisions",
	})
}

func GetCommision(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Get commision",
	})
}

func CreateCommision(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Create commision",
	})
}
