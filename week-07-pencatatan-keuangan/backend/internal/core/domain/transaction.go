package domain

import "time"

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
}
