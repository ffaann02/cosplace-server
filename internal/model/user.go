package model

import (
	"time"
)

type User struct {
	UserId      string     `json:"user_id" gorm:"type:varchar(20);primaryKey"`
	Username    string     `json:"username" gorm:"type:varchar(100);uniqueIndex"`
	Email       string     `json:"email" gorm:"type:varchar(100);uniqueIndex"`
	Password    string     `json:"password" gorm:"type:varchar(255)"`
	FirstName   string     `json:"first_name" gorm:"type:varchar(100)"`
	LastName    string     `json:"last_name" gorm:"type:varchar(100)"`
	Gender      string     `json:"gender" gorm:"type:varchar(10)"`
	DateOfBirth string     `json:"date_of_birth" gorm:"type:date"`
	PhoneNumber string     `json:"phone_number,omitempty" gorm:"type:varchar(20)"`
	CreatedAt   time.Time  `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time  `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

type RegisterRequest struct {
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	PhoneNumber string `json:"phone_number"`
	DateOfBirth string `json:"date_of_birth"`
	Email       string `json:"email"`
	Gender      string `json:"gender"`
	Username    string `json:"username"`
	Password    string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (u *User) TableName() string {
	return "users"
}
