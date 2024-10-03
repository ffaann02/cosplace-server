package handler

import (
	"fmt"
	"os"
	"time"

	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
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

	loginRequest.Email = c.Query("email")
	loginRequest.Username = c.Query("username")
	loginRequest.Password = c.Query("password")

	if err := c.BodyParser(&loginRequest); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	fmt.Println("Parsed Request:", loginRequest)

	for _, user := range users {
		if (user.Email == loginRequest.Email || user.Username == loginRequest.Username) && user.Password == loginRequest.Password {
			// Create the Claims
			claims := jwt.MapClaims{
				"user_id":  user.ID,
				"username": user.Username,
				"exp":      time.Now().Add(time.Hour * 48).Unix(), // Token expires in 72 hours
			}

			// Create token
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

			// Get secret key from environment variable
			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Secret key not found",
				})
			}

			// Generate encoded token and send it as response.
			t, err := token.SignedString([]byte(secret))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Could not generate token",
				})
			}

			fmt.Println("Token:", t)

			c.Cookie(&fiber.Cookie{
				Name:     "jwt",
				Value:    t,
				Expires:  time.Now().Add(time.Hour * 48),
				HTTPOnly: true,
				Secure:   false,
				SameSite: "None",
			})

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
