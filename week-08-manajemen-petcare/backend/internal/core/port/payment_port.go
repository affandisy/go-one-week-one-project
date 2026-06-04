package port

import "github.com/affandisy/petcare-system/internal/core/domain"

// Driven Port (Database)
type PaymentRepository interface {
	SavePaymentAndUpdateInvoice(payment *domain.Payment, invoiceID string) error
	GetInvoiceTotal(invoiceID string) (float64, error)
	CheckInvoiceStatus(invoiceID string) (string, error)
}

// Driving Port (API)
type PaymentUseCase interface {
	ProcessPayment(userID, invoiceID, method string, amount float64, reference string) (*domain.Payment, error)
}
