package usecase

import (
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/google/uuid"
)

type auditUseCase struct {
	repo port.AuditRepository
}

func NewAuditUseCase(repo port.AuditRepository) port.AuditUseCase {
	return &auditUseCase{repo}
}

func (uc *auditUseCase) RecordAction(userID, action, entityName, entityID, payload string) error {
	// Jika sistem digunakan tanpa login (misal via sistem internal/cron), gunakan "SYSTEM"
	if userID == "" {
		userID = "SYSTEM"
	}

	log := &domain.AuditLog{
		ID:         uuid.NewString(),
		UserID:     userID,
		Action:     action,
		EntityName: entityName,
		EntityID:   entityID,
		Payload:    payload,
		CreatedAt:  time.Now(),
	}

	return uc.repo.SaveLog(log)
}
