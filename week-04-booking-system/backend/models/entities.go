package models

import (
	"time"

	"github.com/google/uuid"
)

// Users table
type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Email        string    `gorm:"unique;not null" json:"email"`
	Whatsapp     string    `gorm:"unique;not null" json:"whatsapp"` // Untuk notifikasi
	FullName     string    `gorm:"not null" json:"full_name"`
	PasswordHash string    `gorm:"not null" json:"-"`                // Disembunyikan dari JSON
	Role         string    `gorm:"default:'customer'" json:"role"`   // customer, admin, owner
	IsVerified   bool      `gorm:"default:false" json:"is_verified"` // Setelah verifikasi OTP[cite: 2]
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Courts table[cite: 2]
type Court struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name        string    `gorm:"not null" json:"name"`           // e.g., "Court A Indoor"[cite: 2]
	Type        string    `gorm:"not null" json:"type"`           // indoor, outdoor[cite: 2]
	Status      string    `gorm:"default:'active'" json:"status"` // active, inactive, maintenance[cite: 2]
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// Pricing Rules table untuk dynamic pricing[cite: 2]
type PricingRule struct {
	ID         uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	CourtID    uuid.UUID `gorm:"type:uuid;not null;index" json:"court_id"`
	Court      Court     `gorm:"foreignKey:CourtID;constraint:OnDelete:CASCADE;" json:"court"`
	DayType    string    `gorm:"not null" json:"day_type"`                        // weekday, weekend[cite: 2]
	StartTime  string    `gorm:"type:time;not null" json:"start_time"`            // e.g., "08:00"[cite: 2]
	EndTime    string    `gorm:"type:time;not null" json:"end_time"`              // e.g., "10:00"[cite: 2]
	BasePrice  float64   `gorm:"type:decimal(10,2);not null" json:"base_price"`   // harga per jam[cite: 2]
	Multiplier float64   `gorm:"type:decimal(3,2);default:1.0" json:"multiplier"` // surge multiplier[cite: 2]
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
