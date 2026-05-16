package sqlite

import (
	"time"

	"github.com/affandisy/financial-app/internal/core/domain"
	"github.com/affandisy/financial-app/internal/core/port"
	"gorm.io/gorm"
)

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) port.ReportRepository {
	return &reportRepository{db}
}

func (r *reportRepository) GetTransactionsForReport(walletID string, start, end time.Time) ([]domain.Transaction, error) {
	var models []TransactionModel // TransactionModel ini sudah ada dari langkah sebelumnya

	// Preload Category dan tembak indeks secara presisi
	err := r.db.Preload("Category").
		Where("wallet_id = ? AND date_time >= ? AND date_time < ?", walletID, start, end).
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	// Mapping dari Model DB ke Entitas Domain
	var transactions []domain.Transaction
	for _, m := range models {
		transactions = append(transactions, domain.Transaction{
			ID:         m.ID,
			WalletID:   m.WalletID,
			Type:       m.Type,
			CategoryID: m.CategoryID,
			Amount:     m.Amount,
			DateTime:   m.DateTime,
			Category: domain.Category{
				Name:  m.Category.Name,
				Color: m.Category.Color,
			},
		})
	}

	return transactions, nil
}
