package helper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewProfileID(db *gorm.DB) (string, error) {
	var profiles []model.Profile
	if err := db.Find(&profiles).Error; err != nil {
		return "", err
	}

	if len(profiles) == 0 {
		return "U-1", nil
	}

	// Extract numeric parts and sort them
	var nums []int
	for _, profile := range profiles {
		parts := strings.Split(profile.ProfileID, "-")
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
	newID := fmt.Sprintf("PF-%d", newNum)
	return newID, nil
}
