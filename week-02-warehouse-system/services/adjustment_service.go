package services

import (
	"errors"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type CreateAdjustmentRequest struct {
	ProductID uint   `json:"product_id"`
	Qty       int    `json:"qty"`
	Reason    string `json:"reason"`
}

type AdjustmentService interface {
	ProcessAdjustment(req CreateAdjustmentRequest, userID uint) (*models.StockAdjustment, error)
}

type adjustmentService struct {
	adjRepo     repositories.AdjustmentRepository
	productRepo repositories.ProductRepository
}

func NewAdjustmentService(adjRepo repositories.AdjustmentRepository, productRepo repositories.ProductRepository) AdjustmentService {
	return &adjustmentService{
		adjRepo:     adjRepo,
		productRepo: productRepo,
	}
}

func (s *adjustmentService) ProcessAdjustment(req CreateAdjustmentRequest, userID uint) (*models.StockAdjustment, error) {
	if req.Qty == 0 {
		return nil, errors.New("Quantity penyesuaian tidak boleh 0")
	}

	if req.Reason == "" {
		return nil, errors.New("Alasan penyesuaian wajid diisi")
	}

	product, err := s.productRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, errors.New("Produk tidak ditemukan")
	}

	newStock := product.CurrentStock + req.Qty

	if newStock < 0 {
		return nil, errors.New("Penyesuaian ditolak: Stok akhir tidak boleh kurang dari 0")
	}

	adjustmentData := &models.StockAdjustment{
		ProductID: req.ProductID,
		UserID:    userID,
		Qty:       req.Qty,
		Reason:    req.Reason,
		Status:    "approved",
	}

	err = s.adjRepo.ExecuteAdjustment(adjustmentData, newStock)
	if err != nil {
		return nil, errors.New("Gagal memproses penyesuaian stok")
	}

	return adjustmentData, nil
}
