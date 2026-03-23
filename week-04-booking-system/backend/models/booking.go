package models

import (
	"time"

	"github.com/google/uuid"
)

type Booking struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	BookingCode string    `gorm:"unique;not null" json:"booking_code"` // misal: PDL-260322-001

	UserID uuid.UUID `gorm:"type:uuid;not null;index" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`

	CourtID uuid.UUID `gorm:"type:uuid;not null;index" json:"court_id"`
	Court   Court     `gorm:"foreignKey:CourtID" json:"court"`

	BookingDate time.Time `gorm:"type:date;not null;index" json:"booking_date"`
	StartTime   string    `gorm:"type:time;not null" json:"start_time"` // Format HH:MM
	EndTime     string    `gorm:"type:time;not null" json:"end_time"`   // Format HH:MM

	Status     string     `gorm:"not null;index" json:"status"` // locked, pending, paid, cancelled, expired
	LockExpiry *time.Time `json:"lock_expiry"`                  // Waktu batas bayar
	TotalPrice float64    `gorm:"type:decimal(10,2);not null" json:"total_price"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
