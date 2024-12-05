package model

type Order struct {
	OrderID   string  `json:"order_id" gorm:"type:varchar(50);primaryKey"` // Primary key
	UserID    string  `json:"user_id" gorm:"type:varchar(50)"`
	SellerID  string  `json:"seller_id" gorm:"type:varchar(50)"`
	Amount    float64 `json:"amount" gorm:"type:float"`
	Status    string  `json:"status" gorm:"type:varchar(20)"`
	CreatedAt string  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt string  `json:"updated_at" gorm:"autoUpdateTime"`
}

type OrderResponse struct {
	OrderID    string       `json:"order_id" gorm:"type:varchar(50);primaryKey"` // Primary key
	UserID     string       `json:"user_id" gorm:"type:varchar(50)"`
	Amount     float64      `json:"amount" gorm:"type:float"`
	Status     string       `json:"status" gorm:"type:varchar(20)"`
	CreatedAt  string       `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  string       `json:"updated_at" gorm:"autoUpdateTime"`
	OrderLists []OrderLists `json:"order_lists"` // Add the field to hold order lists
	Products   []Product    `json:"products"`    // Add the field to hold products
}

type OrderLists struct {
	OrderListID int    `json:"order_list_id" gorm:"primaryKey;autoIncrement"` // Primary key with auto increment
	OrderID     string `json:"order_id" gorm:"type:varchar(50);not null"`     // Foreign key to Order
	ProductID   string `json:"product_id" gorm:"type:varchar(50)"`
	Quantity    int    `json:"quantity" gorm:"type:int"`
}

func (Order) TableName() string {
	return "orders"
}

func (OrderLists) TableName() string {
	return "order_lists"
}
