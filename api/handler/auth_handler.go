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
			accessClaims := jwt.MapClaims{
				"user_id":  user.ID,
				"username": user.Username,
				"exp":      time.Now().Add(time.Minute * 15).Unix(), // Token expires in 72 hours
			}

			refreshClaims := jwt.MapClaims{
				"user_id": user.ID,
				"exp":     time.Now().Add(time.Hour * 720).Unix(), // Refresh token expires in 72 hours
			}

			accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
			refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

			secret := os.Getenv("JWT_SECRET")
			if secret == "" {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Secret key not found",
				})
			}

			// Generate encoded token and send it as response.
			accessT, err := accessToken.SignedString([]byte(secret))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Could not generate access token",
				})
			}

			refreshT, err := refreshToken.SignedString([]byte(secret))
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"message": "Could not generate refresh token",
				})
			}

			fmt.Println("Token:", accessT)
			fmt.Println("Refresh Token:", refreshT)

			c.Cookie(&fiber.Cookie{
				Name:     "access_token",
				Value:    accessT,
				Expires:  time.Now().Add(time.Minute * 15),
				HTTPOnly: true,
				Secure:   true,
				SameSite: "None",
			})

			c.Cookie(&fiber.Cookie{
				Name:     "refresh_token",
				Value:    refreshT,
				Expires:  time.Now().Add(time.Hour * 720),
				HTTPOnly: true,
				Secure:   true,
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

func Refresh(c *fiber.Ctx) error {
	refreshToken := c.Cookies("refresh_token")
	if refreshToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Refresh token not provided",
		})
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Secret key not found",
		})
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC and matches what we expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid refresh token",
		})
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "Invalid token claims",
		})
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Unix(int64(exp), 0).Before(time.Now()) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Refresh token has expired",
			})
		}
	}

	userId := claims["user_id"].(float64)

	var user m.User
	for _, u := range users {
		if u.ID == int64(userId) {
			user = u
			break
		}
	}

	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message": "User not found",
		})
	}

	accessClaims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessT, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Could not generate access token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "access_token",
		Value:    accessT,
		Expires:  time.Now().Add(time.Minute * 15),
		HTTPOnly: true,
		Secure:   true,
		SameSite: "None",
	})

	return c.JSON(fiber.Map{
		"message":     "Token refreshed",
		"accessToken": accessT,
	})
}
