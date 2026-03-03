package services

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/config"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type CreateTransactionRequest struct {
	ReferenceNo string               `json:"reference_no"`
	Type        string               `json:"type"` // "INBOUND" atau "OUTBOUND"
	PartnerID   uint                 `json:"partner_id"`
	Notes       string               `json:"notes"`
	Items       []TransactionItemDTO `json:"items"`
}

type TransactionItemDTO struct {
	ProductID  uint    `json:"product_id"`
	Quantity   int     `json:"quantity"`
	UnitPrice  float64 `json:"unit_price"`
	BatchNo    string  `json:"batch_no"`
	ExpiryDate string  `json:"expiry_date"`
}

type TransactionService interface {
	ProcessTransaction(req CreateTransactionRequest, userID uint) (*models.Transaction, error)
	ApproveTransaction(txID uint, approverID uint) error
}

type transactionService struct {
	txRepo      repositories.TransactionRepository
	productRepo repositories.ProductRepository
	partnerRepo repositories.PartnerRepository
	batchRepo   repositories.BatchRepository
}

func NewTransactionService(txRepo repositories.TransactionRepository, productRepo repositories.ProductRepository, partnerRepo repositories.PartnerRepository, batchRepo repositories.BatchRepository) TransactionService {
	return &transactionService{txRepo: txRepo, productRepo: productRepo, partnerRepo: partnerRepo, batchRepo: batchRepo}
}

func (s *transactionService) ProcessTransaction(req CreateTransactionRequest, userID uint) (*models.Transaction, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("Transaksi harus memiliki minimal 1 barang")
	}

	if req.Type != "INBOUND" && req.Type != "OUTBOUND" {
		return nil, errors.New("Tipe transaksi harus INBOUND atau OUTBOUND")
	}

	partner, err := s.partnerRepo.FindByID(req.PartnerID)
	if err != nil {
		return nil, errors.New("Partner (Supplier/Customer) tidak ditemukan")
	}

	if req.Type == "INBOUND" && partner.Type != "SUPPLIER" {
		return nil, errors.New("Transaksi INBOUND hanya dapat dilakukan dari partner bertipe SUPPLIER")
	}

	if req.Type == "OUTBOUND" && partner.Type != "CUSTOMER" {
		return nil, errors.New("Transaksi OUTBOUND hanya dapat dilakukan ke partner bertipe CUSTOMER")
	}

	txData := &models.Transaction{
		ReferenceNo:     req.ReferenceNo,
		TransactionDate: time.Now(),
		Type:            req.Type,
		Status:          "draft",
		PartnerID:       req.PartnerID,
		Notes:           req.Notes,
		CreatedByID:     userID,
	}

	// stockUpdates := make(map[uint]int)

	for _, itemReq := range req.Items {
		txData.Items = append(txData.Items, models.TransactionItem{
			ProductID: itemReq.ProductID,
			Quantity:  itemReq.Quantity,
			UnitPrice: itemReq.UnitPrice,
			SubTotal:  float64(itemReq.Quantity) * itemReq.UnitPrice,
		})
	}

	err = s.txRepo.ExecuteTransaction(txData, map[uint]int{})
	if err != nil {
		return nil, errors.New("Gagal menyimpan draft transaksi")
	}

	return txData, nil
}

func (s *transactionService) ApproveTransaction(txID uint, approverID uint) error {
	tx, err := s.txRepo.FindByID(txID)
	if err != nil {
		return errors.New("Transaksi tidak ditemukan")
	}

	if tx.Status == "approved" {
		return errors.New("Transaksi ini sudah pernah di-approve")
	}

	productIDsMap := make(map[uint]bool)
	var productIDs []int
	for _, item := range tx.Items {
		if !productIDsMap[item.ProductID] {
			productIDsMap[item.ProductID] = true
			productIDs = append(productIDs, int(item.ProductID))
		}
	}

	sort.Ints(productIDs)

	var lockedKeys []string

	defer func() {
		for _, key := range lockedKeys {
			config.RedisClient.Del(config.Ctx, key)
		}
	}()

	for _, id := range productIDs {
		lockKey := fmt.Sprintf("lock:product:%d", id)

		isLocked, err := config.RedisClient.SetNX(config.Ctx, lockKey, "locked", 10*time.Second).Result()
		if err != nil || !isLocked {
			return errors.New("Sistem sedang memproses transaksi untuk barang ini. Silakan coba beberapa saat lagi (Race Condition dicegah)")
		}

		lockedKeys = append(lockedKeys, lockKey)
	}

	stockUpdates := make(map[uint]int)
	batchUpdates := []models.ProductBatch{}

	for _, item := range tx.Items {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return errors.New("Ada barang yang tidak valid dalam transaksi ini")
		}

		newStock := product.CurrentStock
		if val, exists := stockUpdates[product.ID]; exists {
			newStock = val
		}

		if tx.Type == "INBOUND" {
			newStock += item.Quantity
		} else if tx.Type == "OUTBOUND" {
			if newStock < item.Quantity {
				return errors.New("Persetujuan ditolak! Stok " + product.Name + " tidak mencukupi saat ini")
			}
			newStock -= item.Quantity

			batches, _ := s.batchRepo.FindByProductID(product.ID)

			sort.Slice(batches, func(i, j int) bool {
				return batches[i].ExpiryDate.Before(batches[j].ExpiryDate)
			})

			qtyNeeded := item.Quantity

			for _, batch := range batches {
				if qtyNeeded <= 0 {
					break
				}

				take := batch.Stock
				if qtyNeeded < batch.Stock {
					take = qtyNeeded
				}

				batch.Stock -= take
				qtyNeeded -= take

				batchUpdates = append(batchUpdates, batch)
			}

			if qtyNeeded > 0 {
				return errors.New("Sistem anomali: Stok global cukup, tetapi stok batch tidak cukup")
			}
		}

		stockUpdates[product.ID] = newStock
	}

	tx.Status = "approved"
	tx.ApprovedByID = &approverID

	err = s.txRepo.ApproveAndUpdateStock(tx, stockUpdates)
	if err != nil {
		return errors.New("Gagal menyimpan persetujuan ke database")
	}

	return nil
}
