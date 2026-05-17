package domain

import "time"

type Category struct {
	ID       string
	Name     string
	Icon     string
	Type     string
	Color    string
	IsActive bool
}

type Transaction struct {
	ID         string
	WalletID   string
	Type       string // "income" atau "expense"
	CategoryID string
	Amount     float64
	Note       string
	DateTime   time.Time

	// Untuk kebutuhan join/tampilan riwayat
	Category Category
}
