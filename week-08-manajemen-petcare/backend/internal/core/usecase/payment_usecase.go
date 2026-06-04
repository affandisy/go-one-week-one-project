package usecase

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/google/uuid"
)

type paymentUseCase struct {
	repo    port.PaymentRepository
	auditUC port.AuditUseCase
}

func NewPaymentUseCase(repo port.PaymentRepository, auditUC port.AuditUseCase) port.PaymentUseCase {
	return &paymentUseCase{repo, auditUC}
}

func (uc *paymentUseCase) ProcessPayment(userID, invoiceID, method string, amount float64, reference string) (*domain.Payment, error) {
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

	payment := &domain.Payment{
		ID:        uuid.NewString(),
		InvoiceID: invoiceID,
		Method:    method,
		Amount:    amount,
		PaidAt:    time.Now(),
		Reference: reference,
	}

	if err := uc.repo.SavePaymentAndUpdateInvoice(payment, invoiceID); err != nil {
		return nil, errors.New("gagal memproses pembayaran")
	}

	// 4. Catat Audit Trail! (Jalankan secara Asinkron menggunakan Goroutine agar tidak memperlambat HTTP Response)
	go func() {
		// Serialize detail pembayaran ke JSON
		payloadBytes, _ := json.Marshal(payment)

		// Rekam jejak
		_ = uc.auditUC.RecordAction(
			userID,
			"PROCESS_PAYMENT",
			"Payment",
			payment.ID,
			string(payloadBytes),
		)
	}()

	return payment, nil
}
