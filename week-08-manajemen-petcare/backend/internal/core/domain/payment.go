package domain

import "time"

type Payment struct {
	ID        string
	InvoiceID string
	Method    string // "Cash", "Transfer", "QRIS"
	Amount    float64
	PaidAt    time.Time
	Reference string // Nomor referensi transfer/QR
}
