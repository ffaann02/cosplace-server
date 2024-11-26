package model

import (
	"time"
)

type Seller struct {
	SellerID          string     `json:"seller_id" gorm:"type:varchar(10);primaryKey"`
	UserID            string     `json:"user_id" gorm:"type:varchar(10)"`
	ShopType          string     `json:"shop_type" gorm:"type:varchar(50)"`
	ShopName          string     `json:"shop_name" gorm:"type:varchar(100)"`
	ShopDesc          string     `json:"shop_desc" gorm:"type:varchar(300)"`
	ProfileImageURL   string     `json:"profile_image_url" gorm:"type:varchar(300)"`
	Verify            bool       `json:"verify" gorm:"type:tinyint(1)"`
	AcceptCreditCard  bool       `json:"accept_credit_card" gorm:"type:tinyint(1)"`
	AcceptQrPromptpay bool       `json:"accept_qr_prompt_pay" gorm:"type:tinyint(1)"`
	ExternalLink      string     `json:"external_link" gorm:"type:varchar(300)"`
	BankName          string     `json:"bank_name" gorm:"type:varchar(50)"`
	BankAccountNumber string     `json:"bank_account_number" gorm:"type:varchar(12)"`
	CreatedAt         time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt         time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt         *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type SellerResponse struct {
	Seller
	SellerUsername string `json:"seller_username"`
}

func (Seller) TableName() string {
	return "sellers"
}
