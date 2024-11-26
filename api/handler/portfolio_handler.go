// handler/handler.go
package handler

import (
	"fmt"
	"time"

	"github.com/ffaann02/cosplace-server/api/helper"
	"github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
)

func CreatePortfolio(c *fiber.Ctx) error {
	// Parse the request body to get the product data
	fmt.Println("Create Portfolio")
	var portfolio m.Portfolio
	db := config.MysqlDB()
	if err := c.BodyParser(&portfolio); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Set the created_at and updated_at fields
	portfolio.CreatedAt = time.Now()
	portfolio.UpdatedAt = time.Now()

	// Start a database transaction
	tx := db.Begin()
	portfolioID, err := helper.GenerateNewPortfolioID(db)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate new product ID"})
	}
	portfolio.PortfolioID = portfolioID

	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}

	// Create the product in the database
	if err := tx.Create(&portfolio).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	fmt.Print("Create Port")
	return c.JSON(fiber.Map{"message": "Product created successfully", "portfolio_id": portfolioID})
}

func GetPortfolioByUserID(c *fiber.Ctx) error {
	db := config.MysqlDB()
	user_id := c.Params("user_id")

	var portfolios []m.Portfolio
	// Preload the portfolio images
	if err := db.Preload("PortfolioImages").Where("created_by = ?", user_id).Find(&portfolios).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query portfolio",
		})
	}

	// Return portfolio with images
	return c.JSON(portfolios)
}
