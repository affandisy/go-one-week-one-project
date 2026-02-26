package services

import (
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type DashboardMetrics struct {
	TotalProductItems  int     `json:"total_product_items"`
	TotalStockValueIDR float64 `json:"total_stock_value_idr"`
	CriticalStockItems int     `json:"critical_stock_items"`
	TodayTransactions  int     `json:"today_transactions"`
}

type DashboardService interface {
	GetDashboardSummary() (DashboardMetrics, error)
}

type dashboardService struct {
	productRepo repositories.ProductRepository
	txRepo      repositories.TransactionRepository
}

func NewDashboardService(productRepo repositories.ProductRepository, txRepo repositories.TransactionRepository) DashboardService {
	return &dashboardService{productRepo: productRepo, txRepo: txRepo}
}

func (s *dashboardService) GetDashboardSummary() (DashboardMetrics, error) {
	var metrics DashboardMetrics

	products, err := s.productRepo.FindAll()
	if err != nil {
		return metrics, err
	}

	metrics.TotalProductItems = len(products)

	for _, p := range products {
		metrics.TotalStockValueIDR += float64(p.CurrentStock) * p.Price

		if p.CurrentStock <= p.MinStock {
			metrics.CriticalStockItems++
		}
	}

	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())

	todayTx, err := s.txRepo.FindTransactionsByDateRange(startOfDay, endOfDay)
	if err != nil {
		return metrics, err
	}

	metrics.TodayTransactions = len(todayTx)

	return metrics, nil
}
