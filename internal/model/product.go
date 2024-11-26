package model

import (
	"time"
)

type Product struct {
	ProductID      string         `json:"product_id" gorm:"type:varchar(10);primaryKey"`
	Name           string         `json:"name" gorm:"type:varchar(100);not null"`
	Price          float64        `json:"price" gorm:"not null"`
	Quantity       int            `json:"quantity" gorm:"not null"`
	Rent           bool           `json:"rent" gorm:"default:false"`
	RentDeposit    float64        `json:"rent_deposit" gorm:"default:0"`
	RentReturnDate time.Time      `json:"rent_return_date"`
	Description    string         `json:"description" gorm:"type:varchar(300);not null"`
	Category       string         `json:"category" gorm:"type:varchar(50);not null"`
	Condition      string         `json:"condition" gorm:"type:varchar(50);not null"`
	Size           string         `json:"size" gorm:"type:varchar(5);not null"`
	Region         string         `json:"region" gorm:"type:varchar(50);not null"`
	CreatedBy      string         `json:"created_by" gorm:"type:varchar(10);not null;index"`
	CreatedAt      time.Time      `json:"created_at" gorm:"autoCreateTime;not null"`
	UpdatedAt      time.Time      `json:"updated_at" gorm:"autoUpdateTime;not null"`
	DeletedAt      *time.Time     `json:"deleted_at"`
	Seller         Seller         `json:"seller" gorm:"foreignKey:CreatedBy;references:SellerID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ProductImages  []ProductImage `json:"product_images" gorm:"foreignKey:ProductID;references:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ProductResponse struct {
	Product
	ProductImages []ProductImage `json:"product_images"`
	SellerID      string         `json:"seller_id"`
}

type ProductImage struct {
	ProductImageID int     `json:"product_image_id" gorm:"primaryKey;autoIncrement"`
	ProductID      string  `json:"product_id" gorm:"type:varchar(10);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ImageURL       string  `json:"image_url" gorm:"type:varchar(255);not null"`
	Product        Product `json:"product" gorm:"foreignKey:ProductID;references:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

type ProductImageResponse struct {
	ProductImageID int     `json:"product_image_id" gorm:"primaryKey;autoIncrement"`
	ProductID      string  `json:"product_id" gorm:"type:varchar(10);not null;index;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	ImageURL       string  `json:"image_url" gorm:"type:varchar(255);not null"`
	Product        Product `json:"product" gorm:"foreignKey:ProductID;references:ProductID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (Product) TableName() string {
	return "products"
}

func (ProductImage) TableName() string {
	return "product_images"
}
