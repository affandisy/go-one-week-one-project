package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username     string         `gorm:"unique;not null" json:"username"`
	PasswordHash string         `gorm:"not null" json:"-"` // Disembunyikan dari response JSON
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	Progress     []UserProgress `gorm:"foreignKey:UserID" json:"progress,omitempty"`
}

type Module struct {
	ID          uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Title       string     `gorm:"not null" json:"title"`
	Description string     `json:"description"`
	LevelOrder  int        `gorm:"unique;not null" json:"level_order"`
	IsLocked    bool       `gorm:"default:true" json:"is_locked"` // Secara default terkunci, kecuali level 1
	Materials   []Material `gorm:"foreignKey:ModuleID" json:"materials,omitempty"`
}

type Material struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ModuleID      uuid.UUID `gorm:"type:uuid;not null;index" json:"module_id"`
	ContentType   string    `gorm:"not null" json:"content_type"` // learn_card, quiz_mcq, quiz_match, quiz_unscramble
	Question      string    `gorm:"not null" json:"question"`
	CorrectAnswer string    `gorm:"not null" json:"correct_answer"`
	Options       string    `gorm:"type:jsonb" json:"options"` // Disimpan sebagai JSON string ["A", "B", "C"]
	ImageURL      string    `json:"image_url"`
	AudioURL      string    `json:"audio_url"`
	DisplayOrder  int       `gorm:"default:0" json:"display_order"`
}

type UserProgress struct {
	ID           uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	UserID       uuid.UUID  `gorm:"type:uuid;not null;index:idx_user_module,unique" json:"user_id"`
	ModuleID     uuid.UUID  `gorm:"type:uuid;not null;index:idx_user_module,unique" json:"module_id"`
	Completed    bool       `gorm:"default:false" json:"completed"`
	BestScore    float64    `gorm:"type:decimal(5,2);default:0" json:"best_score"`
	Attempts     int        `gorm:"default:0" json:"attempts"`
	LastQuizDate *time.Time `json:"last_quiz_date"`
}
