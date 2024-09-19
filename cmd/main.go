package main

import (
	"fmt"

	"github.com/ffaann02/cosplace-server/api"
	"github.com/gofiber/fiber/v2"
)

func main() {
	app := fiber.New()
	api.SetupRoutes(app)
	app.Listen(":3000")
	fmt.Println("Server is running on port 3000")
}
