package port

import (
	"time"

	"github.com/affandisy/financial-app/internal/core/domain"
)

// Driven Port (Outbound) -> Akan diimplementasikan oleh SQLite
type ReportRepository interface {
	GetTransactionsForReport(walletID string, startDate, endDate time.Time) ([]domain.Transaction, error)
}

// Driving Port (Inbound) -> Akan diimplementasikan oleh UseCase
type ReportUseCase interface {
	GetMonthlySummary(walletID string, year int, month time.Month) (*domain.MonthlyReport, error)
}
