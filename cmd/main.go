package main

import (
	"fmt"

	"github.com/ffaann02/cosplace-server/api"
	"github.com/ffaann02/cosplace-server/internal/config"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	config.InitDB()
	app := fiber.New()
	api.SetupRoutes(app)
	app.Listen(":3000")
	fmt.Println("Server is running on port 3000")
}
