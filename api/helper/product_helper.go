package helper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewProductID(db *gorm.DB) (string, error) {
	var products []model.Product
	if err := db.Find(&products).Error; err != nil {
		return "", err
	}

	if len(products) == 0 {
		return "P-1", nil
	}

	// Extract numeric parts and sort them
	var nums []int
	for _, product := range products {
		parts := strings.Split(product.ProductID, "-")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid product_id format")
		}
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}
		nums = append(nums, num)
	}

	sort.Ints(nums)
	newNum := nums[len(nums)-1] + 1
	newID := fmt.Sprintf("P-%d", newNum)
	return newID, nil
}
