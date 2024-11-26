package model

import (
	"time"
)

type Portfolio struct {
	PortfolioID     string           `gorm:"primaryKey;size:10" json:"portfolio_id"`
	Title           string           `gorm:"size:100" json:"title"`
	Description     string           `gorm:"size:300" json:"description"`
	CreatedBy       string           `gorm:"size:10" json:"created_by"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       time.Time        `json:"deleted_at"`
	PortfolioImages []PortfolioImage `gorm:"foreignKey:PortfolioID;references:PortfolioID" json:"portfolio_images"`
}

type PortfolioImage struct {
	PortfolioImageID int       `gorm:"primaryKey;autoIncrement" json:"portfolio_image_id"`
	PortfolioID      string    `gorm:"size:10;not null" json:"portfolio_id"`
	ImageURL         string    `gorm:"size:255" json:"image_url"`
	CreatedAt        time.Time `json:"created_at"`
	DeletedAt        time.Time `json:"deleted_at"`
}

type PortfolioResponse struct {
	PortfolioID     string           `gorm:"primaryKey;size:10" json:"portfolio_id"`
	Title           string           `gorm:"size:100" json:"title"`
	Description     string           `gorm:"size:300" json:"description"`
	CreatedBy       string           `gorm:"size:10" json:"created_by"`
	CreatedAt       time.Time        `json:"created_at"`
	UpdatedAt       time.Time        `json:"updated_at"`
	DeletedAt       time.Time        `json:"deleted_at"`
	PortfolioImages []PortfolioImage `gorm:"foreignKey:PortfolioID" json:"portfolio_images"`
}

func (Portfolio) TableName() string {
	return "portfolios"
}

func (PortfolioImage) TableName() string {
	return "portfolio_images"
}
