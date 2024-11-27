package helper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	m "github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewFriendshipID(db *gorm.DB) (string, error) {
	var friendships []m.Friendship
	if err := db.Find(&friendships).Error; err != nil {
		return "", err
	}

	if len(friendships) == 0 {
		return "FS-1", nil
	}

	// Extract numeric parts and sort them
	var nums []int
	for _, friendship := range friendships {
		parts := strings.Split(friendship.FriendshipID, "-")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid friendship_id format")
		}
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}
		nums = append(nums, num)
	}

	sort.Ints(nums)
	newNum := nums[len(nums)-1] + 1
	newID := fmt.Sprintf("FS-%d", newNum)
	return newID, nil
}
