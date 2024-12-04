package helper

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/ffaann02/cosplace-server/internal/model"
	"gorm.io/gorm"
)

func GenerateNewCustomPostID(db *gorm.DB) (string, error) {
	var posts []model.CustomPost
	if err := db.Find(&posts).Error; err != nil {
		return "", err
	}

	if len(posts) == 0 {
		return "CP-1", nil
	}

	// Extract numeric parts and sort them
	var nums []int
	for _, post := range posts {
		parts := strings.Split(post.PostID, "-")
		if len(parts) != 2 {
			return "", fmt.Errorf("invalid post_id format")
		}
		num, err := strconv.Atoi(parts[1])
		if err != nil {
			return "", err
		}
		nums = append(nums, num)
	}

	sort.Ints(nums)
	newNum := nums[len(nums)-1] + 1
	newID := fmt.Sprintf("CP-%d", newNum)
	return newID, nil
}
