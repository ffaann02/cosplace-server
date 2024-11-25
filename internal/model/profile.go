package model

import (
	"time"
)

type Profile struct {
	ProfileID       string    `json:"profile_id" gorm:"type:varchar(10);primaryKey"`
	UserID          string    `json:"user_id" gorm:"type:varchar(10);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	DisplayName     string    `json:"display_name" gorm:"type:varchar(50)"`
	ProfileImageURL string    `json:"profile_image_url" gorm:"type:varchar(255)"`
	CoverImageURL   string    `json:"cover_image_url" gorm:"type:varchar(255)"`
	Bio             string    `json:"bio" gorm:"type:varchar(200)"`
	InstagramURL    string    `json:"instagram_url" gorm:"type:varchar(100)"`
	TwitterURL      string    `json:"twitter_url" gorm:"type:varchar(100)"`
	FacebookURL     string    `json:"facebook_url" gorm:"type:varchar(100)"`
	CreatedAt       time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	User            User      `json:"user" gorm:"foreignKey:UserID;references:UserID"`
}

type ProfileResponse struct {
	ProfileID       string    `json:"profile_id"`
	UserID          string    `json:"user_id"`
	SellerID        string    `json:"seller_id"`
	DisplayName     string    `json:"display_name"`
	ProfileImageURL string    `json:"profile_image_url"`
	CoverImageURL   string    `json:"cover_image_url"`
	Bio             string    `json:"bio"`
	InstagramURL    string    `json:"instagram_url"`
	TwitterURL      string    `json:"twitter_url"`
	FacebookURL     string    `json:"facebook_url"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
	Username        string    `json:"username"`
	Gender          string    `json:"gender"`
	DateOfBirth     string    `json:"date_of_birth"`
}

func (Profile) TableName() string {
	return "profiles"
}
