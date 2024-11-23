package helper

import (
	"fmt"
	"strconv"
	"strings"

	m "github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewSellerID(db *gorm.DB) (string, error) {
	var lastSeller m.Seller
	if err := db.Order("seller_id").First(&lastSeller).Error; err != nil {
		// If no previous user, set the first user ID
		if err == gorm.ErrRecordNotFound {
			return "S-1", nil
		}
		return "", err
	}

	parts := strings.Split(lastSeller.SellerID, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid seller_id format")
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}
	newID := fmt.Sprintf("S-%d", num+1)
	return newID, nil
}
