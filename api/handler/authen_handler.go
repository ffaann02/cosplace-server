package handler

import (
	"fmt"

	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
)

var users = []m.User{
	{
		ID:       1,
		Username: "user1",
		Email:    "test@gmail.com",
		Password: "123456",
	},
	{
		ID:       2,
		Username: "user2",
		Email:    "test2@gmail.com",
		Password: "123456",
	},
}

func Register(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Register",
	})
}

func Login(c *fiber.Ctx) error {
	var loginRequest m.LoginRequest
	fmt.Printf("%+v", c.Queries())

	loginRequest.Email = c.Query("email")
	loginRequest.Username = c.Query("username")
	loginRequest.Password = c.Query("password")

	for _, user := range users {
		if (user.Email == loginRequest.Email || user.Username == loginRequest.Username) && user.Password == loginRequest.Password {
			return c.JSON(fiber.Map{
				"message": "Login success",
				"user":    user,
			})
		}
	}

	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
		"message": "Invalid credentials",
	})
}

func Logout(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"message": "Logout",
	})
}
