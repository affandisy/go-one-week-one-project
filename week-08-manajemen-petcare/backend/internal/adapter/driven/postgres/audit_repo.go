package postgres

import (
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"gorm.io/gorm"
)

type AuditLogModel struct {
	ID         string    `gorm:"type:uuid;primaryKey"`
	UserID     string    `gorm:"type:uuid;index"`
	Action     string    `gorm:"type:varchar(50);index"`
	EntityName string    `gorm:"type:varchar(50)"`
	EntityID   string    `gorm:"type:uuid;index"`
	Payload    string    `gorm:"type:jsonb"` // JSONB untuk optimasi pencarian detail di Postgres
	CreatedAt  time.Time `gorm:"index"`
}

func (AuditLogModel) TableName() string { return "audit_logs" }

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) port.AuditRepository {
	db.AutoMigrate(&AuditLogModel{})
	return &auditRepository{db}
}

func (r *auditRepository) SaveLog(log *domain.AuditLog) error {
	model := AuditLogModel{
		ID:         log.ID,
		UserID:     log.UserID,
		Action:     log.Action,
		EntityName: log.EntityName,
		EntityID:   log.EntityID,
		Payload:    log.Payload,
		CreatedAt:  log.CreatedAt,
	}
	return r.db.Create(&model).Error
}
