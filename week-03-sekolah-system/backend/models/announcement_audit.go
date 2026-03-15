package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Announcement struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Title       string    `gorm:"not null" json:"title"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	TargetRoles string    `gorm:"type:jsonb;default:'[]'" json:"target_roles"`
	CreatedByID uuid.UUID `gorm:"not null" json:"created_by_id"`
	CreatedBy   User      `gorm:"foreignKey:CreatedByID" json:"created_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AuditLog struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	ActorUserID uuid.UUID `gorm:"not null" json:"actor_user_id"`
	Action      string    `gorm:"not null" json:"action"`      // CREATE, UPDATE, DELETE, LOGIN
	EntityType  string    `gorm:"not null" json:"entity_type"` // e.g., "Student", "Grade", "Payslip"
	EntityID    string    `json:"entity_id"`                   // ID dari entitas yang diubah
	BeforeData  string    `gorm:"type:jsonb" json:"before_data"`
	AfterData   string    `gorm:"type:jsonb" json:"after_data"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	CreatedAt   time.Time `json:"created_at"`
}

func (m *Announcement) BeforeCreate(tx *gorm.DB) (err error) { m.ID = uuid.New(); return }
func (m *AuditLog) BeforeCreate(tx *gorm.DB) (err error)     { m.ID = uuid.New(); return }
