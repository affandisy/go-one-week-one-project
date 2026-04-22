package repositories

import (
	"github.com/affandisy/pos-system/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateTransaction(trx *models.Transaction, movements []models.StockMovement) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

func (r *transactionRepository) CreateTransaction(trx *models.Transaction, movements []models.StockMovement) error {
	// Memulai Database Transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan Header Transaksi
		if err := tx.Create(trx).Error; err != nil {
			return err
		}

		// 2. Simpan Detail Transaksi & Kurangi Stok Produk
		for _, detail := range trx.Details {
			// Kurangi stok produk
			if err := tx.Model(&models.Product{}).
				Where("id = ?", detail.ProductID).
				UpdateColumn("stock", gorm.Expr("stock - ?", detail.Quantity)).Error; err != nil {
				return err
			}
		}

		// 3. Catat Riwayat Pergerakan Stok (Barang Keluar)
		for _, movement := range movements {
			if err := tx.Create(&movement).Error; err != nil {
				return err
			}
		}

		return nil // Commit transaksi
	})
}
