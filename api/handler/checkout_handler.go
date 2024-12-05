package handler

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/ffaann02/cosplace-server/internal/config"
	"github.com/ffaann02/cosplace-server/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"gorm.io/gorm"
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
			"image_url": productImages[0].ImageURL,
			"price":     product.Price,
			"quantity":  quantity,
		})
	}

	return c.JSON(fiber.Map{
		"products":    productList,
		"totalAmount": totalAmount,
	})
}

func CreateOrder(c *fiber.Ctx) error {
	type OrderRequest struct {
		UserID   string `json:"user_id" validate:"required"`
		SellerID string `json:"seller_id" validate:"required"`
		Products []struct {
			ProductID string `json:"product_id" validate:"required"`
			Quantity  int    `json:"quantity" validate:"required"`
		} `json:"products" validate:"required"`
	}

	var req OrderRequest

	// Parse and validate the request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	db := config.MysqlDB()

	responseOrderID := ""
	// Start a database transaction
	err := db.Transaction(func(tx *gorm.DB) error {
		// Generate a UUID for the order_id (also used as order_number)
		orderID := uuid.New().String()
		// orderID = orderID[:8] + orderID[9:] // Remove separators
		responseOrderID = orderID
		// Calculate the total amount
		var totalAmount float64
		for _, product := range req.Products {
			totalAmount += float64(product.Quantity * 10) // Mock calculation
		}

		// Create the order record
		order := model.Order{
			OrderID:   orderID,
			UserID:    req.UserID,
			SellerID:  req.SellerID,
			Amount:    totalAmount,
			Status:    "paid",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
		}

		if err := tx.Create(&order).Error; err != nil {
			return err // Rollback transaction
		}

		// Create the order list records
		for _, product := range req.Products {
			orderList := model.OrderLists{
				OrderID:   orderID,
				ProductID: product.ProductID,
				Quantity:  product.Quantity,
			}
			if err := tx.Create(&orderList).Error; err != nil {
				return err // Rollback transaction
			}
		}

		// Commit transaction by returning nil
		return nil
	})

	// Check for transaction errors
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Transaction failed", "details": err.Error()})
	}

	// Return the response
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "Order created successfully",
		"order_id": responseOrderID,
	})
}

func GetOrderByOrderID(c *fiber.Ctx) error {
	// Get the order_id from the query parameters
	orderID := c.Query("order_id")
	if orderID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "order_id is required",
		})
	}

	db := config.MysqlDB()

	// Fetch the order from the orders table
	var order model.Order
	if err := db.Where("order_id = ?", orderID).First(&order).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": "Order not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve order",
		})
	}

	// Fetch the order items from the order_lists table
	var orderLists []model.OrderLists
	if err := db.Where("order_id = ?", orderID).Find(&orderLists).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve order lists",
		})
	}

	// Fetch product details for each order list item
	var products []model.ProductResponse
	for _, orderList := range orderLists {
		var product model.ProductResponse
		if err := db.Preload("ProductImages").Where("product_id = ?", orderList.ProductID).First(&product).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to retrieve product details for product_id: %s", orderList.ProductID),
			})
		}
		products = append(products, product)
	}

	// Prepare the response
	response := struct {
		Order      model.Order             `json:"order"`
		OrderLists []model.OrderLists      `json:"items"`
		Products   []model.ProductResponse `json:"products"`
	}{
		Order:      order,
		OrderLists: orderLists,
		Products:   products,
	}

	// Return the response in JSON format
	return c.Status(fiber.StatusOK).JSON(response)
}
func GetAllOrdersByUserID(c *fiber.Ctx) error {
	db := config.MysqlDB()

	// Get the user_id from the URL parameters
	userID := c.Params("user_id")
	fmt.Println("UserID:", userID)

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "user_id is required",
		})
	}

	// Fetch the orders from the orders table for the specific user
	var orders []model.Order
	if err := db.Where("user_id = ?", userID).Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve orders",
		})
	}

	// Prepare to hold the order details
	var orderDetails []model.OrderResponse

	// Loop through orders and fetch associated order items and products
	for _, order := range orders {
		// Fetch associated order lists
		var orderLists []model.OrderLists
		if err := db.Where("order_id = ?", order.OrderID).Find(&orderLists).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to retrieve order lists for order_id: %s", order.OrderID),
			})
		}

		// Fetch associated products for each order list
		var products []model.Product
		for _, orderList := range orderLists {
			var product model.Product
			// Preload the ProductImages related to this product
			if err := db.Preload("ProductImages").Where("product_id = ?", orderList.ProductID).First(&product).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to retrieve product details for product_id: %s", orderList.ProductID),
				})
			}
			products = append(products, product)
		}

		// Create an OrderResponse instance
		orderResponse := model.OrderResponse{
			OrderID:    order.OrderID,
			UserID:     order.UserID,
			Amount:     order.Amount,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
			OrderLists: orderLists, // Add the associated order lists
			Products:   products,   // Add the associated products
		}

		// Add the order response to the orderDetails slice
		orderDetails = append(orderDetails, orderResponse)
	}

	// Return the order details in the response
	return c.Status(fiber.StatusOK).JSON(orderDetails)
}

func GetAllOrdersBySellerID(c *fiber.Ctx) error {
	db := config.MysqlDB()

	// Get the user_id from the URL parameters
	sellerID := c.Params("seller_id")
	fmt.Println("SellerID:", sellerID)

	if sellerID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "seller_id is required",
		})
	}

	// Fetch the orders from the orders table for the specific user
	var orders []model.Order
	if err := db.Where("seller_id = ?", sellerID).Find(&orders).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve orders",
		})
	}

	// Prepare to hold the order details
	var orderDetails []model.OrderResponse

	// Loop through orders and fetch associated order items and products
	for _, order := range orders {
		// Fetch associated order lists
		var orderLists []model.OrderLists
		if err := db.Where("order_id = ?", order.OrderID).Find(&orderLists).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": fmt.Sprintf("Failed to retrieve order lists for order_id: %s", order.OrderID),
			})
		}

		// Fetch associated products for each order list
		var products []model.Product
		for _, orderList := range orderLists {
			var product model.Product
			// Preload the ProductImages related to this product
			if err := db.Preload("ProductImages").Where("product_id = ?", orderList.ProductID).First(&product).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": fmt.Sprintf("Failed to retrieve product details for product_id: %s", orderList.ProductID),
				})
			}
			products = append(products, product)
		}

		// Create an OrderResponse instance
		orderResponse := model.OrderResponse{
			OrderID:    order.OrderID,
			UserID:     order.UserID,
			Amount:     order.Amount,
			Status:     order.Status,
			CreatedAt:  order.CreatedAt,
			UpdatedAt:  order.UpdatedAt,
			OrderLists: orderLists, // Add the associated order lists
			Products:   products,   // Add the associated products
		}

		// Add the order response to the orderDetails slice
		orderDetails = append(orderDetails, orderResponse)
	}

	// Return the order details in the response
	return c.Status(fiber.StatusOK).JSON(orderDetails)
}
