package models

import (
	"time"

	"gorm.io/gorm"
)

type Partner struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(150);not null" json:"name"`
	Type      string         `gorm:"type:varchar(20);not null" json:"type"` // "SUPPLIER" atau "CUSTOMER"
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
