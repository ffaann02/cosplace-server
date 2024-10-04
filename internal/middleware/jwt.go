package middleware

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWTProtected() fiber.Handler {
	return func(c *fiber.Ctx) error {
		cookie := c.Cookies("access_token")

		fmt.Println("Cookie:", cookie)
		if cookie == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Unauthorized, token missing",
			})
		}

		secret := os.Getenv("JWT_SECRET")
		token, err := jwt.Parse(cookie, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid or expired token",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			c.Locals("user_id", claims["user_id"])
			c.Locals("username", claims["username"])

			return c.Next()
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "Invalid token claims",
			})
		}
	}
}
