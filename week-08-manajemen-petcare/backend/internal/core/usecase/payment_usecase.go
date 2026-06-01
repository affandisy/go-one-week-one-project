package usecase

import (
	"errors"
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/google/uuid"
)

type paymentUseCase struct {
	repo port.PaymentRepository
}

func NewPaymentUseCase(repo port.PaymentRepository) port.PaymentUseCase {
	return &paymentUseCase{repo}
}

func (uc *paymentUseCase) ProcessPayment(invoiceID, method string, amount float64, reference string) (*domain.Payment, error) {
	// 1. Validasi Metode
	validMethods := map[string]bool{"Cash": true, "Transfer": true, "QRIS": true}
	if !validMethods[method] {
		return nil, errors.New("metode pembayaran tidak valid")
	}

	// 2. Validasi Status Tagihan
	status, err := uc.repo.CheckInvoiceStatus(invoiceID)
	if err != nil {
		return nil, errors.New("tagihan tidak ditemukan")
	}
	if status == "Paid" {
		return nil, errors.New("tagihan ini sudah lunas")
	}

	// 3. Validasi Nominal Pembayaran
	total, err := uc.repo.GetInvoiceTotal(invoiceID)
	if err != nil {
		return nil, errors.New("gagal mengambil detail tagihan")
	}
	if amount < total {
		return nil, errors.New("nominal pembayaran kurang dari total tagihan")
	}

	// 4. Buat Entitas Pembayaran
	payment := &domain.Payment{
		ID:        uuid.NewString(),
		InvoiceID: invoiceID,
		Method:    method,
		Amount:    amount,
		PaidAt:    time.Now(),
		Reference: reference,
	}

	// 5. Simpan Pembayaran dan Update Status Invoice secara Atomik
	if err := uc.repo.SavePaymentAndUpdateInvoice(payment, invoiceID); err != nil {
		return nil, errors.New("gagal memproses pembayaran")
	}

	return payment, nil
}
