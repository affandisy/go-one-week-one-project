package port

import "github.com/affandisy/petcare-system/internal/core/domain"

// Driven Ports (Outbound ke Database)
type BillingRepository interface {
	SaveInvoice(invoice *domain.Invoice) error
}

type NutritionRepository interface {
	SaveLog(log *domain.NutritionLog) error
	GetLogsByPet(petID string) ([]domain.NutritionLog, error)
}

// Driving Ports (Inbound dari API)
type BillingUseCase interface {
	GenerateInvoice(ownerID string, items []domain.InvoiceItem) (*domain.Invoice, error)
}

type NutritionUseCase interface {
	RecordDiet(petID string, brand string, calories int, notes string) error
	GetLogsByPet(petID string) ([]domain.NutritionLog, error)
}
