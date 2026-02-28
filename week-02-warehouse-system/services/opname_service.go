package services

import (
	"errors"
	"fmt"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type OpnameItemReq struct {
	ProductID uint `json:"product_id"`
	ActualQty int  `json:"actual_qty"`
}

type ProcessOpnameReq struct {
	Notes string          `json:"notes"`
	Items []OpnameItemReq `json:"items"`
}

type OpnameService interface {
	ProcessMassOpname(req ProcessOpnameReq, userID uint) (int, error)
}

type opnameService struct {
	adjRepo     repositories.AdjustmentRepository
	productRepo repositories.ProductRepository
}

func NewOpnameService(adjRepo repositories.AdjustmentRepository, productRepo repositories.ProductRepository) OpnameService {
	return &opnameService{adjRepo: adjRepo, productRepo: productRepo}
}

func (s *opnameService) ProcessMassOpname(req ProcessOpnameReq, userID uint) (int, error) {
	if len(req.Items) == 0 {
		return 0, errors.New("Daftar opname tidak boleh kosong")
	}

	stockUpdates := make(map[uint]int)
	var adjustments []models.StockAdjustment
	discrepancyCount := 0

	for _, item := range req.Items {
		if item.ActualQty < 0 {
			return 0, errors.New("Quantity fisik tidak boleh negatif")
		}

		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return 0, fmt.Errorf("Produk dengan ID %d tidak ditemukan", item.ProductID)
		}

		diff := item.ActualQty - product.CurrentStock

		if diff != 0 {
			discrepancyCount++
			stockUpdates[product.ID] = item.ActualQty

			adjustments = append(adjustments, models.StockAdjustment{
				ProductID: product.ID,
				UserID:    userID,
				Qty:       diff,
				Reason:    fmt.Sprintf("Opname Massal: %s", req.Notes),
				Status:    "approved",
			})
		}
	}

	if discrepancyCount == 0 {
		return 0, nil
	}

	err := s.adjRepo.ExecuteMassAdjustment(adjustments, stockUpdates)
	if err != nil {
		return 0, errors.New("Gagal memproses opname massal")
	}

	return discrepancyCount, nil
}
