package services

import (
	"errors"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type CreateTransactionRequest struct {
	ReferenceNo string               `json:"reference_no"`
	Type        string               `json:"type"` // "INBOUND" atau "OUTBOUND"
	EntityName  string               `json:"entity_name"`
	Notes       string               `json:"notes"`
	Items       []TransactionItemDTO `json:"items"`
}

type TransactionItemDTO struct {
	ProductID uint    `json:"product_id"`
	Quantity  int     `json:"quantity"`
	UnitPrice float64 `json:"unit_price"`
}

type TransactionService interface {
	ProcessTransaction(req CreateTransactionRequest, userID uint) (*models.Transaction, error)
}

type transactionService struct {
	txRepo      repositories.TransactionRepository
	productRepo repositories.ProductRepository
}

func NewTransactionService(txRepo repositories.TransactionRepository, productRepo repositories.ProductRepository) TransactionService {
	return &transactionService{txRepo: txRepo, productRepo: productRepo}
}

func (s *transactionService) ProcessTransaction(req CreateTransactionRequest, userID uint) (*models.Transaction, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("Transaksi harus memiliki minimal 1 barang")
	}

	if req.Type != "INBOUND" && req.Type != "OUTBOUND" {
		return nil, errors.New("Tipe transaksi harus INBOUND atau OUTBOUND")
	}

	txData := &models.Transaction{
		ReferenceNo:     req.ReferenceNo,
		TransactionDate: time.Now(),
		Type:            req.Type,
		Status:          "approved",
		EntityName:      req.EntityName,
		Notes:           req.Notes,
		CreatedByID:     userID,
	}

	stockUpdates := make(map[uint]int)

	for _, itemReq := range req.Items {
		if itemReq.Quantity <= 0 {
			return nil, errors.New("Quantity barang tidak valid")
		}

		product, err := s.productRepo.FindByID(itemReq.ProductID)
		if err != nil {
			return nil, errors.New("Barang dengan ID tersebut tidak ditemukan")
		}

		newStock := product.CurrentStock

		if req.Type == "INBOUND" {
			newStock += itemReq.Quantity
		} else if req.Type == "OUTBOUND" {
			if product.CurrentStock < itemReq.Quantity {
				return nil, errors.New("Stok tidak mencukupi untuk barang: " + product.Name)
			}
			newStock -= itemReq.Quantity
		}

		stockUpdates[product.ID] = newStock

		txData.Items = append(txData.Items, models.TransactionItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: itemReq.UnitPrice,
			SubTotal:  float64(itemReq.Quantity) * itemReq.UnitPrice,
		})
	}

	err := s.txRepo.ExecuteTransaction(txData, stockUpdates)
	if err != nil {
		return nil, errors.New("Gagal memproses transaksi: " + err.Error())
	}

	return txData, nil
}
