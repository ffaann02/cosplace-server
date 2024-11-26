package helper

import (
	"fmt"
	"strconv"
	"strings"

	m "github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewCustomPostID(db *gorm.DB) (string, error) {
	var lastPost m.CustomPost
	if err := db.Order("post_id desc").First(&lastPost).Error; err != nil {
		// If no previous user, set the first user ID
		if err == gorm.ErrRecordNotFound {
			return "CP-1", nil
		}
		return "", err
	}

	parts := strings.Split(lastPost.PostID, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid post_id format")
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}
	newID := fmt.Sprintf("CP-%d", num+1)
	return newID, nil
}
