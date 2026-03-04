package models

import (
	"time"

	"gorm.io/gorm"
)

type Warehouse struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"` // Contoh: JKT-01
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`            // Contoh: Gudang Utama Jakarta
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
