// handler/handler.go
package handler

import (
	"fmt"
	"net/url"
	"time"

	"github.com/ffaann02/cosplace-server/api/helper"
	"github.com/ffaann02/cosplace-server/internal/config"
	m "github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
)

func GetAllProducts(c *fiber.Ctx) error {
	db := config.MysqlDB()

	var products []m.Product
	if err := db.Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get products"})
	}

	var productResponses []m.ProductResponse
	for _, product := range products {
		var productImages []m.ProductImage
		if err := db.Where("product_id = ?", product.ProductID).Find(&productImages).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product images"})
		}
		productResponse := m.ProductResponse{
			Product:       product,
			ProductImages: productImages,
		}
		productResponses = append(productResponses, productResponse)
	}

	return c.JSON(productResponses)
}

func GetSellerProducts(c *fiber.Ctx) error {
	db := config.MysqlDB()
	sellerID := c.Query("seller_id") // Get seller_id from query parameters

	if sellerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "seller_id is required"})
	}

	var products []m.Product
	if err := db.Where("created_by = ?", sellerID).Find(&products).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get products"})
	}

	var productResponses []m.ProductResponse
	for _, product := range products {
		var productImages []m.ProductImage
		if err := db.Where("product_id = ?", product.ProductID).Find(&productImages).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product images"})
		}
		productResponse := m.ProductResponse{
			Product:       product,
			ProductImages: productImages,
		}
		productResponses = append(productResponses, productResponse)
	}

	return c.JSON(productResponses)
}

func GetProductByID(c *fiber.Ctx) error {
	db := config.MysqlDB()
	productID := c.Params("product_id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "product_id is required"})
	}

	var product m.Product
	if err := db.Where("product_id = ?", productID).First(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product"})
	}

	var productImages []m.ProductImage
	if err := db.Where("product_id = ?", productID).Find(&productImages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product images"})
	}

	productResponse := m.ProductResponse{
		Product:       product,
		ProductImages: productImages,
	}

	return c.JSON(productResponse)
}

func GetProductByName(c *fiber.Ctx) error {
	db := config.MysqlDB()
	ProductName := c.Params("product_name")

	if ProductName == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "product_name is required"})
	}

	ProductName, err := url.QueryUnescape(ProductName)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to decode product name"})
	}

	var product m.Product
	if err := db.Where("name = ?", ProductName).First(&product).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product"})
	}

	var productImages []m.ProductImage
	if err := db.Where("product_id = ?", product.ProductID).Find(&productImages).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product images"})
	}

	productResponse := m.ProductResponse{
		Product:       product,
		ProductImages: productImages,
		SellerID:      product.CreatedBy,
	}

	return c.JSON(productResponse)
}

func UpdateProduct(c *fiber.Ctx) error {
	db := config.MysqlDB()
	productID := c.Params("product_id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "product_id is required"})
	}

	// Parse the request body to get the product data
	var product m.Product
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	// Set the updated_at field
	product.UpdatedAt = time.Now()

	// Start a database transaction
	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}

	// Update the product in the database
	if err := tx.Where("product_id = ?", productID).Updates(&product).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update product"})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.JSON(fiber.Map{"message": "Product updated successfully"})
}

func CreateProduct(c *fiber.Ctx) error {
	// Parse the request body to get the product data
	var product m.Product
	db := config.MysqlDB()
	if err := c.BodyParser(&product); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	fmt.Println(product)
	fmt.Printf("by: %s\n", product.CreatedBy)

	// Set the created_at and updated_at fields
	product.CreatedAt = time.Now()
	product.UpdatedAt = time.Now()

	// Start a database transaction
	tx := db.Begin()
	productID, err := helper.GenerateNewProductID(tx)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate new product ID"})
	}
	product.ProductID = productID

	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}

	// Create the product in the database
	if err := tx.Create(&product).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	fmt.Print("CreateProduct")
	return c.JSON(fiber.Map{"message": "Product created successfully", "product_id": product.ProductID})
}

func DeleteProduct(c *fiber.Ctx) error {
	db := config.MysqlDB()
	productID := c.Params("product_id")

	if productID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "product_id is required"})
	}

	tx := db.Begin()
	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start transaction"})
	}

	if err := tx.Where("product_id = ?", productID).Delete(&m.ProductImage{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product images"})
	}

	if err := tx.Where("product_id = ?", productID).Delete(&m.Product{}).Error; err != nil {
		tx.Rollback()
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete product"})
	}

	if err := tx.Commit().Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transaction"})
	}

	return c.JSON(fiber.Map{"message": "Product deleted successfully"})
}

func UploadProducImages(c *fiber.Ctx) error {
	fmt.Print("UploadProducImages")
	return c.JSON(fiber.Map{"message": "UploadProducImages"})
}

func GetAllProductBySellerID(c *fiber.Ctx) error {
	db := config.MysqlDB()
	sellerID := c.Query("seller_id")

	var products []m.Product
	if err := db.Find(&products).Where("seller_id = ?", sellerID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get products"})
	}

	var productResponses []m.ProductResponse
	for _, product := range products {
		var productImages []m.ProductImage
		if err := db.Where("product_id = ?", product.ProductID).Find(&productImages).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product images"})
		}
		productResponse := m.ProductResponse{
			Product:       product,
			ProductImages: productImages,
		}
		productResponses = append(productResponses, productResponse)
	}

	return c.JSON(productResponses)
}
