package port

import (
	"time"

	"github.com/affandisy/financial-app/internal/core/domain"
)

// Driven Port (Outbound) -> Akan diimplementasikan oleh adapter/sqlite
type TransactionRepository interface {
	FindByWalletAndDateRange(walletID string, start, end time.Time) ([]domain.Transaction, error)
	Save(transaction *domain.Transaction) error
}
