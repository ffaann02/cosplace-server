package helper

import (
	"fmt"
	"strconv"
	"strings"

	m "github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewProfileID(db *gorm.DB) (string, error) {
	var lastProfile m.Profile
	if err := db.Order("profile_id desc").First(&lastProfile).Error; err != nil {
		// If no previous user, set the first user ID
		if err == gorm.ErrRecordNotFound {
			return "P-1", nil
		}
		return "", err
	}

	parts := strings.Split(lastProfile.UserID, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid profile_id format")
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}
	newID := fmt.Sprintf("P-%d", num+1)
	return newID, nil
}
