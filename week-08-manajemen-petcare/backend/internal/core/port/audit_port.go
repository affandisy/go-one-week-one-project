package port

import "github.com/affandisy/petcare-system/internal/core/domain"

// Driven Port (Database)
type AuditRepository interface {
	SaveLog(log *domain.AuditLog) error
}

// Driving Port (Dipanggil oleh UseCase lain)
type AuditUseCase interface {
	RecordAction(userID, action, entityName, entityID, payload string) error
}
