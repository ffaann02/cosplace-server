package helper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewUserID(db *gorm.DB) (string, error) {
	var users []model.User
	if err := db.Find(&users).Error; err != nil {
		return "", err
	}

	if len(users) == 0 {
		return "U-1", nil
	}

	// Extract numeric parts and sort them
	var nums []int
	for _, user := range users {
		parts := strings.Split(user.UserID, "-")
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
	newID := fmt.Sprintf("U-%d", newNum)
	return newID, nil
}
