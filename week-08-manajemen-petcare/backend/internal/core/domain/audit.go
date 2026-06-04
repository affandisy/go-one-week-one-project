package domain

import "time"

type AuditLog struct {
	ID         string
	UserID     string    // Siapa yang melakukan? (Diambil dari JWT)
	Action     string    // Contoh: "CREATE_PAYMENT", "UPDATE_DIET"
	EntityName string    // Contoh: "Payment", "NutritionLog"
	EntityID   string    // ID dari data yang diubah
	Payload    string    // Detail data (Bisa berupa JSON string)
	CreatedAt  time.Time // Kapan dilakukan?
}