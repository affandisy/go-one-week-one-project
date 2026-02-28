package repositories

import (
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"gorm.io/gorm"
)

type TransactionRepository interface {
	ExecuteTransaction(txData *models.Transaction, stockUpdates map[uint]int) error
	FindTransactionsByDateRange(startDate, endDate time.Time) ([]models.Transaction, error)
	FindByID(id uint) (*models.Transaction, error)
	ApproveAndUpdateStock(tx *models.Transaction, stockUpdates map[uint]int) error
	FindOutboundItemsByDate(startDate, endDate time.Time) ([]models.TransactionItem, error)
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

func (r *transactionRepository) FindTransactionsByDateRange(startDate, endDate time.Time) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.Preload("Partner").Preload("Items").Where("transaction_date >= ? AND transaction_date <= ?", startDate, endDate).Find(&transactions).Error

	return transactions, err
}

func (r *transactionRepository) FindByID(id uint) (*models.Transaction, error) {
	var tx models.Transaction

	err := r.db.Preload("Items").First(&tx, id).Error

	return &tx, err
}

func (r *transactionRepository) ApproveAndUpdateStock(tx *models.Transaction, stockUpdates map[uint]int) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Save(tx).Error; err != nil {
			return err
		}

		for productID, newStock := range stockUpdates {
			if err := dbTx.Model(&models.Product{}).Where("id = ?", productID).Update("current_stock", newStock).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

func (r *transactionRepository) FindOutboundItemsByDate(startDate, endDate time.Time) ([]models.TransactionItem, error) {
	var items []models.TransactionItem

	err := r.db.Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Where("transactions.type = ? AND transactions.transaction_date >= ? AND transactions.transaction_date <= ?", "OUTBOUND", startDate, endDate).
		Preload("Product").Find(&items).Error

	return items, err
}
