package postgres

import (
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"gorm.io/gorm"
)

type PaymentModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	InvoiceID string `gorm:"type:uuid;index"`
	Method    string
	Amount    float64
	PaidAt    time.Time
	Reference string
}

func (PaymentModel) TableName() string { return "payments" }

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) port.PaymentRepository {
	db.AutoMigrate(&PaymentModel{})
	return &paymentRepository{db}
}

func (r *paymentRepository) CheckInvoiceStatus(invoiceID string) (string, error) {
	var status string
	err := r.db.Model(&InvoiceModel{}).Select("status").Where("id = ?", invoiceID).Scan(&status).Error
	return status, err
}

func (r *paymentRepository) GetInvoiceTotal(invoiceID string) (float64, error) {
	var total float64
	err := r.db.Model(&InvoiceModel{}).Select("total_amount").Where("id = ?", invoiceID).Scan(&total).Error
	return total, err
}

// Implementasi Database Transaction (Atomicity)
func (r *paymentRepository) SavePaymentAndUpdateInvoice(payment *domain.Payment, invoiceID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan rekam pembayaran
		model := PaymentModel{
			ID:        payment.ID,
			InvoiceID: payment.InvoiceID,
			Method:    payment.Method,
			Amount:    payment.Amount,
			PaidAt:    payment.PaidAt,
			Reference: payment.Reference,
		}
		if err := tx.Create(&model).Error; err != nil {
			return err // Rollback
		}

		// 2. Update status tagihan menjadi 'Paid'
		if err := tx.Model(&InvoiceModel{}).Where("id = ?", invoiceID).Update("status", "Paid").Error; err != nil {
			return err // Rollback
		}

		return nil // Commit
	})
}
