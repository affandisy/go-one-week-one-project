package models

import (
	"time"

	"gorm.io/gorm"
)

type Location struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"` // Contoh: RAK-A1-01
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`            // Contoh: Rak Makanan Ringan
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
