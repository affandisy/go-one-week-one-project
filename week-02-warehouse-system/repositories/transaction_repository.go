package repositories

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	ExecuteTransaction(txData *models.Transaction, stockUpdates map[uint]int) error
}

type transactionRepository struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) ExecuteTransaction(txData *models.Transaction, stockUpdates map[uint]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(txData).Error; err != nil {
			return err
		}

		for productID, newStock := range stockUpdates {
			if err := tx.Model(&models.Product{}).Where("id = ?", productID).Update("current_stock", newStock).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
