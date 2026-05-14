package domain

import "time"

// Wallet murni Go struct, tanpa tag JSON/GORM
type Wallet struct {
	ID        string
	Name      string
	Balance   float64
	Currency  string
	CreatedAt time.Time
}

// Perilaku Domain (Domain Behavior)
func (w *Wallet) AddBalance(amount float64) {
	w.Balance += amount
}

func (w *Wallet) DeductBalance(amount float64) {
	w.Balance -= amount
}
