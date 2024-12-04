package handler

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ffaann02/cosplace-server/internal/config"
	"github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
)

func GetCheckout(c *fiber.Ctx) error {
	list := c.Query("list")
	quantity := c.Query("quantity")

	if list == "" || quantity == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Missing list or quantity parameter",
		})
	}

	fmt.Println("Checkout Handler")
	fmt.Println(list)
	fmt.Println(quantity)

	productIDs := strings.Split(list, ",")
	quantities := strings.Split(quantity, ",")

	if len(productIDs) != len(quantities) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Mismatched list and quantity lengths",
		})
	}

	db := config.MysqlDB()
	var products []model.Product
	if err := db.Where("product_id IN ?", productIDs).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to query products",
		})
	}

	var totalAmount float64
	var productList []fiber.Map

	for i, product := range products {
		qty := quantities[i]
		quantity, err := strconv.Atoi(qty)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Invalid quantity value",
			})
		}

		totalAmount += float64(quantity) * product.Price

		var productImages []model.ProductImage
		if err := db.Where("product_id = ?", product.ProductID).First(&productImages).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to query product images",
			})
		}

		productList = append(productList, fiber.Map{
			"seller_id": product.CreatedBy,
			"name":      product.Name,
			"image":     productImages[0].ImageURL,
			"price":     product.Price,
			"quantity":  quantity,
		})
	}

	return c.JSON(fiber.Map{
		"products":    productList,
		"totalAmount": totalAmount,
	})
}
