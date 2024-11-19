package helper

import (
	"fmt"
	"strconv"
	"strings"

	m "github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewFriendshipID(db *gorm.DB) (string, error) {
	var lastFriendship m.Friendship
	if err := db.Order("friendship_id desc").First(&lastFriendship).Error; err != nil {
		// If no previous user, set the first user ID
		if err == gorm.ErrRecordNotFound {
			return "FS-1", nil
		}
		return "", err
	}

	parts := strings.Split(lastFriendship.FriendshipID, "-")
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid friendship_id format")
	}

	num, err := strconv.Atoi(parts[1])
	if err != nil {
		return "", err
	}
	newID := fmt.Sprintf("FS-%d", num+1)
	return newID, nil
}
