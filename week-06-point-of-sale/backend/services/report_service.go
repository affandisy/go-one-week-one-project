package services

import (
	"time"

	"github.com/affandisy/pos-system/models"
	"github.com/affandisy/pos-system/repositories"
)

type ProductSoldDetail struct {
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Omzet    float64 `json:"omzet"`
}

type DailySalesSummary struct {
	Date              string              `json:"date"`
	TotalTransactions int                 `json:"total_transactions"`
	TotalOmzet        float64             `json:"total_omzet"`
	TotalProfit       float64             `json:"total_profit"`  // Tambahan: Margin Kotor
	CashTotal         float64             `json:"cash_total"`    // Tambahan: Rincian Pembayaran
	QRISTotal         float64             `json:"qris_total"`    // Tambahan: Rincian Pembayaran
	ProductsSold      []ProductSoldDetail `json:"products_sold"` // Tambahan: Daftar Barang Terjual
}
type ReportService interface {
	GetDailySalesSummary(dateString string) (*DailySalesSummary, error)
	GetLowStockAlerts() ([]models.Product, error)
}

type reportService struct {
	repo     repositories.ReportRepository
	prodRepo repositories.ProductRepository
}

func NewReportService(repo repositories.ReportRepository, prodRepo repositories.ProductRepository) ReportService {
	return &reportService{repo, prodRepo}
}

func (s *reportService) GetDailySalesSummary(dateString string) (*DailySalesSummary, error) {
	var targetDate time.Time
	var err error

	if dateString == "" {
		targetDate = time.Now()
	} else {
		targetDate, err = time.Parse("2006-01-02", dateString)
		if err != nil {
			return nil, err
		}
	}

	startDate := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, targetDate.Location())
	endDate := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 23, 59, 59, 0, targetDate.Location())

	transactions, err := s.repo.GetDailySales(startDate, endDate)
	if err != nil {
		return nil, err
	}

	summary := &DailySalesSummary{
		Date:              startDate.Format("2006-01-02"),
		TotalTransactions: len(transactions),
		ProductsSold:      []ProductSoldDetail{},
	}

	// Map untuk merekap barang yang terjual
	productRecap := make(map[string]*ProductSoldDetail)

	for _, trx := range transactions {
		summary.TotalOmzet += trx.FinalAmount

		// Pisahkan berdasarkan metode pembayaran
		if trx.PaymentMethod == "cash" {
			summary.CashTotal += trx.FinalAmount
		} else {
			summary.QRISTotal += trx.FinalAmount
		}

		// Hitung profit dan rekap barang terjual
		for _, detail := range trx.Details {
			// Cari data produk untuk mengetahui Harga Beli (Modal)
			product, _ := s.prodRepo.FindByID(detail.ProductID.String())

			// Margin Kotor = (Harga Jual saat transaksi - Harga Beli) * Kuantitas
			if product != nil {
				profitPerItem := (detail.Price - product.PurchasePrice) * float64(detail.Quantity)
				summary.TotalProfit += profitPerItem

				// Masukkan ke rekap barang terjual
				if _, exists := productRecap[product.Name]; !exists {
					productRecap[product.Name] = &ProductSoldDetail{Name: product.Name}
				}
				productRecap[product.Name].Quantity += detail.Quantity
				productRecap[product.Name].Omzet += detail.Subtotal
			}
		}
	}

	// Ubah Map menjadi Array untuk JSON
	for _, item := range productRecap {
		summary.ProductsSold = append(summary.ProductsSold, *item)
	}

	return summary, nil
}

func (s *reportService) GetLowStockAlerts() ([]models.Product, error) {
	return s.repo.GetLowStockProducts()
}
