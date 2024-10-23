package model

import "time"

type User struct {
	UserID      string    `gorm:"primaryKey"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	PhoneNumber string    `gorm:"unique" json:"phone_number"`
	Address     string    `json:"address"`
	Pin         string    `json:"pin"`
	Balance     int       `json:"balance"`
	CreatedDate time.Time `json:"created_date"`
	UpdatedDate time.Time `json:"updated_date"`
}

type Transaction struct {
	TransactionID string `gorm:"primaryKey"`
	// TopUpID       string    `json:"top_up_id,omitempty"`
	UserID        string    `json:"user_id"`
	Type          string    `json:"type"`
	Amount        int       `json:"amount"`
	Remarks       string    `json:"remarks"`
	BalanceBefore int       `json:"balance_before"`
	BalanceAfter  int       `json:"balance_after"`
	CreatedDate   time.Time `json:"created_date"`
}
