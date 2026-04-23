package repositories

import (
	"time"

	"github.com/affandisy/pos-system/models"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetDailySales(startDate, endDate time.Time) ([]models.Transaction, error)
	GetLowStockProducts() ([]models.Product, error)
}

type reportRepository struct {
	db *gorm.DB
}

func NewReportRepository(db *gorm.DB) ReportRepository {
	return &reportRepository{db}
}

// Mengambil semua transaksi dalam rentang waktu tertentu (misal: hari ini)
func (r *reportRepository) GetDailySales(startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction
	// Preload Details agar kita tahu persis barang apa saja yang terjual
	err := r.db.Preload("Details").
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Find(&transactions).Error
	return transactions, err
}

// Mengambil barang yang stoknya berada di bawah atau sama dengan batas minimum
func (r *reportRepository) GetLowStockProducts() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Where("stock <= min_stock").Order("stock asc").Find(&products).Error
	return products, err
}
