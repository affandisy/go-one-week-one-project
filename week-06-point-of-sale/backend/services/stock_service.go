package services

import (
	"errors"

	"github.com/affandisy/pos-system/models"
	"github.com/affandisy/pos-system/repositories"
	"github.com/google/uuid"
)

type StockInRequest struct {
	ProductID  string  `json:"product_id"`
	SupplierID *string `json:"supplier_id,omitempty"` // Opsional
	Quantity   int     `json:"quantity"`
	Note       string  `json:"note"`
}

type StockService interface {
	RecordStockIn(req StockInRequest) error
	GetStockHistory(productID string) ([]models.StockMovement, error)
}

type stockService struct {
	repo repositories.StockRepository
}

func NewStockService(repo repositories.StockRepository) StockService {
	return &stockService{repo}
}

func (s *stockService) RecordStockIn(req StockInRequest) error {
	if req.ProductID == "" {
		return errors.New("ID produk tidak boleh kosong")
	}
	if req.Quantity <= 0 {
		return errors.New("kuantitas barang masuk harus lebih dari 0")
	}

	prodID, err := uuid.Parse(req.ProductID)
	if err != nil {
		return errors.New("format ID produk tidak valid")
	}

	movement := &models.StockMovement{
		ProductID: prodID,
		Type:      "IN", // Menandakan barang masuk
		Quantity:  req.Quantity,
		Note:      req.Note,
	}

	// Jika Pemasok (Supplier) dipilih
	if req.SupplierID != nil && *req.SupplierID != "" {
		suppID, err := uuid.Parse(*req.SupplierID)
		if err == nil {
			movement.SupplierID = &suppID
		}
	}

	return s.repo.AddStockIn(movement)
}

func (s *stockService) GetStockHistory(productID string) ([]models.StockMovement, error) {
	return s.repo.GetMovementsByProduct(productID)
}
