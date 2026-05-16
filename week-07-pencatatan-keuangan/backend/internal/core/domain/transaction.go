package domain

import "time"

type Category struct {
	Name  string
	Color string
}

// Transaction murni Go struct, tanpa tag JSON/GORM
type Transaction struct {
	ID         string
	WalletID   string
	Type       string
	CategoryID string
	Amount     float64
	Note       string
	DateTime   time.Time
	CreatedAt  time.Time
	UpdatedAt  time.Time
	Category   Category
}
