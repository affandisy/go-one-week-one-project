package usecase

import (
	"errors"
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/google/uuid"
)

type billingUseCase struct {
	repo port.BillingRepository
}

func NewBillingUseCase(repo port.BillingRepository) port.BillingUseCase {
	return &billingUseCase{repo}
}

func (uc *billingUseCase) GenerateInvoice(ownerID string, inputItems []domain.InvoiceItem) (*domain.Invoice, error) {
	if len(inputItems) == 0 {
		return nil, errors.New("invoice harus memiliki minimal satu layanan")
	}

	invoiceID := uuid.NewString()
	var total float64
	var finalItems []domain.InvoiceItem

	// Memproses dan mengisolasi kalkulasi secara spesifik per hewan
	for _, item := range inputItems {
		if item.PetID == "" {
			return nil, errors.New("setiap layanan wajib diikat pada ID hewan peliharaan spesifik")
		}

		total += item.Price
		finalItems = append(finalItems, domain.InvoiceItem{
			ID:        uuid.NewString(),
			InvoiceID: invoiceID,
			PetID:     item.PetID,
			ServiceID: item.ServiceID,
			Price:     item.Price,
		})
	}

	invoice := &domain.Invoice{
		ID:          invoiceID,
		OwnerID:     ownerID,
		TotalAmount: total,
		Status:      "Unpaid",
		CreatedAt:   time.Now(),
		Items:       finalItems,
	}

	if err := uc.repo.SaveInvoice(invoice); err != nil {
		return nil, err
	}

	return invoice, nil
}
