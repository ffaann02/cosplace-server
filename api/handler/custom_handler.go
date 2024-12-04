package handler

import (
	"github.com/ffaann02/cosplace-server/api/helper"
	"github.com/ffaann02/cosplace-server/internal/config"
	"github.com/ffaann02/cosplace-server/internal/model"
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
	// Define a struct to parse the request body
	type CreateCommisionRequest struct {
		Title           string  `json:"title"`
		Description     string  `json:"description"`
		PriceRangeStart float64 `json:"price_range_start"`
		PriceRangeEnd   float64 `json:"price_range_end"`
		AnimeName       string  `json:"anime_name"`
		CreatedBy       string  `json:"created_by"`
	}

	// Parse the request body
	var req CreateCommisionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Cannot parse JSON",
		})
	}

	// Start a transaction
	db := config.MysqlDB()
	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to start transaction",
		})
	}

	postID, err := helper.GenerateNewCustomPostID(tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate new post ID",
		})
	}

	// Create a new commission
	commission := model.CustomPost{
		PostID:          postID,
		Title:           req.Title,
		Description:     req.Description,
		PriceRangeStart: req.PriceRangeStart,
		PriceRangeEnd:   req.PriceRangeEnd,
		AnimeName:       req.AnimeName,
		CreatedBy:       req.CreatedBy,
	}

	// Save the commission to the database
	if err := tx.Create(&commission).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create commission",
		})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to commit transaction",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Commission created successfully",
		"post_id": commission.PostID,
	})
}

func GetCommisionByUserID(c *fiber.Ctx) error {
	userID := c.Params("user_id")

	var commissions []model.CustomPostResponse
	db := config.MysqlDB()

	// Fetch CustomPost along with their CustomPostRefImages
	if err := db.Model(&model.CustomPost{}).
		Where("created_by = ?", userID).
		Preload("CustomRefImages").
		Find(&commissions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve commissions",
		})
	}

	return c.JSON(fiber.Map{
		"commissions": commissions,
	})
}
