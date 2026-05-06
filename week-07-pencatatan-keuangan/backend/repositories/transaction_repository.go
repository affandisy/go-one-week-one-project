package repositories

import (
	"github.com/affandisy/financial-app/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	CreateWithWalletUpdate(trx *models.Transaction, wallet *models.Wallet) error
	GetRecentTransactions(limit int) ([]models.Transaction, error)
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db}
}

// CreateWithWalletUpdate menyimpan transaksi dan memperbarui saldo dompet secara aman
func (r *transactionRepository) CreateWithWalletUpdate(trx *models.Transaction, wallet *models.Wallet) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan riwayat transaksi
		if err := tx.Create(trx).Error; err != nil {
			return err
		}
		// 2. Perbarui saldo dompet
		if err := tx.Save(wallet).Error; err != nil {
			return err
		}
		return nil
	})
}

// GetRecentTransactions mengambil riwayat transaksi untuk dasbor (FR-003)
func (r *transactionRepository) GetRecentTransactions(limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	// Preload Category agar nama dan ikon kategori ikut terbawa
	err := r.db.Preload("Category").Order("date_time desc").Limit(limit).Find(&transactions).Error
	return transactions, err
}
