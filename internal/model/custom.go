package model

import (
	"time"
)

type CustomPost struct {
	PostID          string     `json:"post_id" gorm:"type:varchar(10);primaryKey"`
	Title           string     `json:"title" gorm:"type:varchar(50);not null"`
	Description     string     `json:"description" gorm:"type:varchar(200);not null"`
	PriceRangeStart float64    `json:"price_range_start" gorm:"not null"`
	PriceRangeEnd   float64    `json:"price_range_end" gorm:"not null"`
	Status          string     `json:"status" gorm:"type:varchar(20);not null"`
	AnimeName       string     `json:"anime_name" gorm:"type:varchar(100);not null"`
	Tag             string     `json:"tag" gorm:"type:varchar(100);not null"`
	CreatedBy       string     `json:"created_by" gorm:"type:varchar(10);not null;index"`
	CreatedAt       time.Time  `json:"created_at" gorm:"autoCreateTime;not null"`
	UpdatedAt       time.Time  `json:"updated_at" gorm:"autoUpdateTime;not null"`
	DeletedAt       *time.Time `json:"deleted_at"`
}

type CustomPostResponse struct {
	CustomPost
	CustomRefImages []CustomPostRefImage `json:"custom_ref_images"`
}

type CustomPostRefImage struct {
	CustomImageID int        `json:"custom_image_id" gorm:"primaryKey;autoIncrement"`
	PostID        string     `json:"post_id" gorm:"type:varchar(10);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ImageURL      string     `json:"image_url" gorm:"type:varchar(255);not null"`
	CustomPost    CustomPost `json:"custom_post" gorm:"foreignKey:PostID;references:PostID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (CustomPost) TableName() string {
	return "custom_posts"
}

func (CustomPostRefImage) TableName() string {
	return "customs_ref_images"
}
