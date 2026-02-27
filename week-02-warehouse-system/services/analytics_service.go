package services

import (
	"sort"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type BestSellerResult struct {
	ProductName string `json:"product_name"`
	TotalOut    int    `json:"total_out"`
}

type AnalyticsService interface {
	GetTopBestSellers(days int, limit int) ([]BestSellerResult, error)
}

type analyticsService struct {
	txRepo repositories.TransactionRepository
}

func NewAnalyticsService(txRepo repositories.TransactionRepository) AnalyticsService {
	return &analyticsService{txRepo: txRepo}
}

func (s *analyticsService) GetTopBestSellers(days int, limit int) ([]BestSellerResult, error) {
	endDate := time.Now()
	startDate := endDate.AddDate(0, 0, -days)

	items, err := s.txRepo.FindOutboundItemsByDate(startDate, endDate)
	if err != nil {
		return nil, err
	}

	movementMap := make(map[string]int)

	for _, item := range items {
		movementMap[item.Product.Name] += item.Quantity
	}

	var results []BestSellerResult
	for name, total := range movementMap {
		results = append(results, BestSellerResult{
			ProductName: name,
			TotalOut:    total,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].TotalOut > results[j].TotalOut
	})

	if len(results) > limit {
		results = results[:limit]
	}

	return results, nil
}
