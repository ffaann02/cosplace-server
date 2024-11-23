package model

import (
	"time"
)

type Seller struct {
	SellerID          string     `json:"seller_id" gorm:"type:varchar(10);primaryKey;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	ShopName          string     `json:"shop_name" gorm:"type:varchar(100);not null"`
	Description       string     `json:"description" gorm:"type:varchar(300)"`
	Verify            bool       `json:"verify" gorm:"default:false"`
	AcceptCreditCard  bool       `json:"accept_credit_card" gorm:"default:false"`
	AcceptQRPromptPay bool       `json:"accept_qr_prompt_pay" gorm:"default:false"`
	JoinedAt          time.Time  `json:"joined_at" gorm:"not null"`
	DeletedAt         *time.Time `json:"deleted_at"`
	User              User       `json:"user" gorm:"foreignKey:SellerID;references:UserID"`
}

func (Seller) TableName() string {
	return "sellers"
}
