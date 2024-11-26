package model

type Order struct {
	OrderID   string  `json:"order_id" gorm:"type:varchar(50);primary_key"`
	UserID    string  `json:"user_id" gorm:"type:varchar(50)"`
	Amount    float64 `json:"amount" gorm:"type:float"`
	Status    string  `json:"status" gorm:"type:varchar(20)"`
	CreatedAt string  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt string  `json:"updated_at" gorm:"autoUpdateTime"`
}

type CheckoutResponse struct {
	OrderID string `json:"order_id"`
	Amount  string `json:"amount"`
}

func (Order) TableName() string {
	return "orders"
}
