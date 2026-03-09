package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID           uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Email        string    `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string    `gorm:"not null" json:"-"` // Disembunyikan dari JSON
	FullName     string    `gorm:"not null" json:"full_name"`
	Phone        string    `json:"phone"`
	Status       string    `gorm:"default:'active'" json:"status"`
	Role         string    `gorm:"not null" json:"role"` // Sederhanakan Role untuk MVP
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
