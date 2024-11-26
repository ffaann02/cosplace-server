package main

import (
	"fmt"
	"log"
	"os"

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
	config.InitAmazonS3()
	app := fiber.New()
	api.SetupRoutes(app)
	var port string = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	var isProduction string = os.Getenv("PRODUCTION")
	if isProduction == "false" {
		log.Fatal(app.Listen("localhost:" + port))
	}
	log.Fatal(app.Listen(":" + port))
}
