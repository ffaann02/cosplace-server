package helper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewSellerID(db *gorm.DB) (string, error) {
	var sellers []model.Seller
	if err := db.Find(&sellers).Error; err != nil {
		return "", err
	}

	if len(sellers) == 0 {
		return "S-1", nil
	}

	// Extract numeric parts and sort them
	var nums []int
	for _, seller := range sellers {
		parts := strings.Split(seller.SellerID, "-")
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
	newID := fmt.Sprintf("S-%d", newNum)
	return newID, nil
}
