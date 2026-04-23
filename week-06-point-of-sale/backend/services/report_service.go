package services

import (
	"time"

	"github.com/affandisy/pos-system/models"
	"github.com/affandisy/pos-system/repositories"
)

// DTO (Data Transfer Object) untuk respons laporan
type DailySalesSummary struct {
	Date              string  `json:"date"`
	TotalTransactions int     `json:"total_transactions"`
	TotalOmzet        float64 `json:"total_omzet"`
}

type ReportService interface {
	GetDailySalesSummary(dateString string) (*DailySalesSummary, error)
	GetLowStockAlerts() ([]models.Product, error)
}

type reportService struct {
	repo repositories.ReportRepository
}

func NewReportService(repo repositories.ReportRepository) ReportService {
	return &reportService{repo}
}

func (s *reportService) GetDailySalesSummary(dateString string) (*DailySalesSummary, error) {
	var targetDate time.Time
	var err error

	// Jika tanggal tidak dikirim, gunakan hari ini
	if dateString == "" {
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("2006-01-02", dateString)
		if err != nil {
			return nil, err
		}
	}

	// Atur jam ke 00:00:00 dan 23:59:59
	startDate := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endDate := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 23, 59, 59, 0, targetDate.Location())

	transactions, err := s.repo.GetDailySales(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var totalOmzet float64
	for _, trx := range transactions {
		totalOmzet += trx.FinalAmount
	}

	summary := &DailySalesSummary{
		Date:              startDate.Format("2006-01-02"),
		TotalTransactions: len(transactions),
		TotalOmzet:        totalOmzet,
	}

	return summary, nil
}

func (s *reportService) GetLowStockAlerts() ([]models.Product, error) {
	return s.repo.GetLowStockProducts()
}
