package repositories

import (
	"time"

	"github.com/affandisy/financial-app/models"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetMonthlyData(startDate, endDate time.Time) ([]models.Transaction, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db}
}

func (r *reportRepository) GetMonthlyData(startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	// Untuk MVP, kita ambil data mentah dan proses di Service.
	// Di masa depan, kueri ini bisa diubah menjadi SQL SUM() dan GROUP BY untuk optimasi.
	err := r.db.Preload("Category").
		Where("date_time >= ? AND date_time < ?", startDate, endDate).
		Find(&transactions).Error
	return transactions, err
}
