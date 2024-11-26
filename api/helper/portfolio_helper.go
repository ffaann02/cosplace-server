package helper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewPortfolioID(db *gorm.DB) (string, error) {
	var portfolios []model.Portfolio
	if err := db.Find(&portfolios).Error; err != nil {
		return "", err
	}

	if len(portfolios) == 0 {
		return "PF-1", nil
	}

	// Extract numeric parts and sort them
	var nums []int
	for _, portfolio := range portfolios {
		parts := strings.Split(portfolio.PortfolioID, "-")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid portfolio_id format")
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
