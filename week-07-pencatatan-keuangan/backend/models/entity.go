// models/entities.go
package models

import "time"

type Wallet struct {
	ID        string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name      string    `json:"name"` // Default: "Dompet Saya"
	Balance   float64   `json:"balance"`
	Currency  string    `json:"currency"` // Default: "IDR"[cite: 1]
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID       string `json:"id" gorm:"type:varchar(36);primaryKey"`
	Name     string `json:"name"`
	Icon     string `json:"icon"`
	Type     string `json:"type"` // "income" atau "expense"[cite: 1]
	Color    string `json:"color"`
	IsActive bool   `json:"is_active"`
}

type Transaction struct {
	ID         string    `json:"id" gorm:"type:varchar(36);primaryKey"`
	WalletID   string    `json:"wallet_id"`
	Type       string    `json:"type"` // "income" atau "expense"[cite: 1]
	CategoryID string    `json:"category_id"`
	Amount     float64   `json:"amount"`
	Note       string    `json:"note"`
	DateTime   time.Time `json:"date_time"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	Category Category `json:"category" gorm:"foreignKey:CategoryID"`
}
