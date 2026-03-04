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
	ReferenceNo       string               `json:"reference_no"`
	Type              string               `json:"type"`
	WarehouseID       uint                 `json:"warehouse_id"`
	TargetWarehouseID *uint                `json:"target_warehouse_id"`
	PartnerID         *uint                `json:"partner_id"`
	Notes             string               `json:"notes"`
	Items             []TransactionItemDTO `json:"items"`
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
	whRepo      repositories.WarehouseRepository
	whStockRepo repositories.WarehouseStockRepository
}

func NewTransactionService(txRepo repositories.TransactionRepository, productRepo repositories.ProductRepository, partnerRepo repositories.PartnerRepository, batchRepo repositories.BatchRepository, whRepo repositories.WarehouseRepository, whStockRepo repositories.WarehouseStockRepository) TransactionService {
	return &transactionService{txRepo: txRepo, productRepo: productRepo, partnerRepo: partnerRepo, batchRepo: batchRepo, whRepo: whRepo, whStockRepo: whStockRepo}
}

func (s *transactionService) ProcessTransaction(req CreateTransactionRequest, userID uint) (*models.Transaction, error) {
	if len(req.Items) == 0 {
		return nil, errors.New("Transaksi harus memiliki minimal 1 barang")
	}

	if req.Type != "INBOUND" && req.Type != "OUTBOUND" && req.Type != "TRANSFER" {
		return nil, errors.New("Tipe transaksi harus INBOUND, OUTBOUND, atau TRANSFER")
	}

	if _, err := s.whRepo.FindByID(req.WarehouseID); err != nil {
		return nil, errors.New("Gudang asal tidak valid")
	}

	if req.Type == "TRANSFER" {
		if req.TargetWarehouseID == nil {
			return nil, errors.New("Gudang tujuan wajib diisi untuk transaksi TRANSFER")
		}
		if req.WarehouseID == *req.TargetWarehouseID {
			return nil, errors.New("Gudang asal dan tujuan tidak boleh sama")
		}
		if _, err := s.whRepo.FindByID(*req.TargetWarehouseID); err != nil {
			return nil, errors.New("Gudang tujuan tidak valid")
		}
	} else {
		if req.PartnerID == nil {
			return nil, errors.New("Partner (Supplier/Customer) wajib diisi untuk INBOUND/OUTBOUND")
		}
		partner, err := s.partnerRepo.FindByID(*req.PartnerID)
		if err != nil {
			return nil, errors.New("Partner (Supplier/Customer) tidak ditemukan")
		}
		if req.Type == "INBOUND" && partner.Type != "SUPPLIER" {
			return nil, errors.New("Transaksi INBOUND hanya dapat dilakukan dari partner bertipe SUPPLIER")
		}
		if req.Type == "OUTBOUND" && partner.Type != "CUSTOMER" {
			return nil, errors.New("Transaksi OUTBOUND hanya dapat dilakukan ke partner bertipe CUSTOMER")
		}
	}

	txData := &models.Transaction{
		ReferenceNo:       req.ReferenceNo,
		TransactionDate:   time.Now(),
		Type:              req.Type,
		Status:            "draft",
		WarehouseID:       req.WarehouseID,
		TargetWarehouseID: req.TargetWarehouseID,
		PartnerID:         req.PartnerID,
		Notes:             req.Notes,
		CreatedByID:       userID,
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

	err := s.txRepo.ExecuteTransaction(txData, map[uint]int{})
	if err != nil {
		return nil, errors.New("Gagal menyimpan draft transaksi")
	}

	return txData, nil
}

func (s *transactionService) ApproveTransaction(txID uint, approverID uint) error {
	tx, err := s.txRepo.FindByID(txID)
	if err != nil || tx.Status == "approved" {
		return errors.New("Transaksi tidak ditemukan atau sudah disetujui")
	}

	var lockKeys []string
	uniqueKeys := make(map[string]bool)

	for _, item := range tx.Items {
		keyOrigin := fmt.Sprintf("lock:warehouse:%d:product:%d", tx.WarehouseID, item.ProductID)
		if !uniqueKeys[keyOrigin] {
			uniqueKeys[keyOrigin] = true
			lockKeys = append(lockKeys, keyOrigin)
		}
		if tx.Type == "TRANSFER" && tx.TargetWarehouseID != nil {
			keyTarget := fmt.Sprintf("lock:warehouse:%d:product:%d", *tx.TargetWarehouseID, item.ProductID)
			if !uniqueKeys[keyTarget] {
				uniqueKeys[keyTarget] = true
				lockKeys = append(lockKeys, keyTarget)
			}
		}
	}

	sort.Strings(lockKeys)

	var acquiredLocks []string
	defer func() {
		for _, key := range acquiredLocks {
			config.RedisClient.Del(config.Ctx, key)
		}
	}()

	for _, key := range lockKeys {
		isLocked, err := config.RedisClient.SetNX(config.Ctx, key, "locked", 10*time.Second).Result()
		if err != nil || !isLocked {
			return errors.New("Sistem sedang memproses transaksi untuk barang ini di gudang terkait. Silakan coba beberapa saat lagi")
		}
		acquiredLocks = append(acquiredLocks, key)
	}

	stockMutations := make(map[string]int)
	var batchUpdates []models.ProductBatch

	for _, item := range tx.Items {
		product, err := s.productRepo.FindByID(item.ProductID)
		if err != nil {
			return errors.New("Ada barang yang tidak valid dalam transaksi ini")
		}

		if tx.Type == "OUTBOUND" || tx.Type == "TRANSFER" {
			stockData, err := s.whStockRepo.GetStock(tx.WarehouseID, item.ProductID)
			if err != nil || stockData.Stock < item.Quantity {
				return fmt.Errorf("Persetujuan ditolak! Stok %s di Gudang Asal tidak mencukupi", product.Name)
			}
			originKey := fmt.Sprintf("%d_%d", tx.WarehouseID, item.ProductID)
			stockMutations[originKey] -= item.Quantity

			// Logika FEFO (Hanya saat barang keluar)
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
				return errors.New("Sistem anomali: Stok gudang cukup, tetapi stok batch FEFO tidak cukup")
			}
		}

		if tx.Type == "INBOUND" {
			targetKey := fmt.Sprintf("%d_%d", tx.WarehouseID, item.ProductID)
			stockMutations[targetKey] += item.Quantity
		}
		if tx.Type == "TRANSFER" && tx.TargetWarehouseID != nil {
			targetKey := fmt.Sprintf("%d_%d", *tx.TargetWarehouseID, item.ProductID)
			stockMutations[targetKey] += item.Quantity
		}
	}

	tx.Status = "approved"
	tx.ApprovedByID = &approverID

	err = s.txRepo.ExecuteMultiWarehouseMutation(tx, stockMutations)
	if err != nil {
		return errors.New("Gagal menyimpan persetujuan mutasi gudang ke database")
	}

	for _, b := range batchUpdates {
		err := s.batchRepo.UpdateBatch(&b)
		if err != nil {
			config.Log.Error().Err(err).Msgf("Gagal update stok batch ID %d", b.ID)
		}
	}

	return nil
}
