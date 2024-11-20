package model

import (
	"time"
)

type Seller struct {
	FriendshipID string    `json:"friendship_id" gorm:"type:varchar(10);primaryKey"`
	UserID       string    `json:"user_id" gorm:"type:varchar(10);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	FriendID     string    `json:"friend_id" gorm:"type:varchar(10);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Status       string    `json:"status" gorm:"type:varchar(10)"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
}

func (Seller) TableName() string {
	return "sellers"
}
