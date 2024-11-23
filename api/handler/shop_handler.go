package handler

import (
	"fmt"
	"os"
	"time"

	"github.com/ffaann02/cosplace-server/api/helper"
	"github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func CreateNewShop(c *fiber.Ctx) error {
	// Parse the request body to get the user_id and new bio
	var requestBody struct {
		UserID            string `json:"user_id"`
		Username          string `json:"username"`
		ShopType          string `json:"shop_type"`
		ShopName          string `json:"shop_name"`
		ShopDesc          string `json:"shop_desc"`
		ProfileImageURL   string `json:"profile_image_url"`
		Verify            bool   `json:"verify"`
		AcceptCreditCard  bool   `json:"accept_credit_card"`
		AcceptQrPromptpay bool   `json:"accept_qr_promptpay"`
		ExternalLink      string `json:"external_link"`
		BankName          string `json:"bank_name"`
		BankAccountNumber string `json:"bank_account_number"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	fmt.Println(requestBody)
	db := config.MysqlDB()
	tx := db.Begin()

	SellerID, err := helper.GenerateNewSellerID(tx)
	if err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "เกิดข้อผิดพลาดในการสร้างร้านค้า",
		})
	}

	// Find the profile with the given user_id
	var seller = m.Seller{
		SellerID:          SellerID,
		UserID:            requestBody.UserID,
		ShopType:          requestBody.ShopType,
		ShopName:          requestBody.ShopName,
		ShopDesc:          requestBody.ShopDesc,
		ProfileImageURL:   requestBody.ProfileImageURL,
		Verify:            requestBody.Verify,
		AcceptCreditCard:  requestBody.AcceptCreditCard,
		AcceptQrPromptpay: requestBody.AcceptQrPromptpay,
		ExternalLink:      requestBody.ExternalLink,
		BankName:          requestBody.BankName,
		BankAccountNumber: requestBody.BankAccountNumber,
	}
	if err := db.Create(&seller).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "เกิดข้อผิดพลาดในการสร้างร้านค้า",
		})
	}

	role := "seller"

	// Create JWT tokens (same as your original code)
	accessClaims := jwt.MapClaims{
		"user_id":   requestBody.UserID,
		"username":  requestBody.Username,
		"role":      role,
		"seller_id": seller.SellerID,
		"exp":       time.Now().Add(time.Minute * 15).Unix(),
	}

	refreshClaims := jwt.MapClaims{
		"user_id":   requestBody.UserID,
		"username":  requestBody.Username,
		"role":      role,
		"seller_id": seller.SellerID,
		"exp":       time.Now().Add(time.Hour * 720).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Secret key not found",
		})
	}

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

	// Set cookies for tokens (same as your original code)
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":       "Create shop successfully",
		"user_id":       requestBody.UserID,
		"username":      requestBody.Username,
		"role":          role,
		"seller_id":     seller.SellerID,
		"access_token":  accessT,
		"refresh_token": refreshT,
	})
}

func GetSellerInfo(c *fiber.Ctx) error {
	// Parse the request body to get the base64-encoded image string and user_id
	var requestBody struct {
		UserID string `json:"user_id"`
	}
	if err := c.BodyParser(&requestBody); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse request body",
		})
	}

	db := config.MysqlDB()
	var seller m.Seller
	if err := db.Where("user_id = ?", requestBody.UserID).First(&seller).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Seller not found",
		})
	}

	return c.JSON(fiber.Map{
		"message": "Seller found",
		"seller":  seller,
	})
}
