package repositories

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"gorm.io/gorm"
)

type ProductMovementResult struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	TotalIn     int     `json:"total_in"`
	TotalOut    int     `json:"total_out"`
	TotalValue  float64 `json:"total_value"`
}

type TransactionRepository interface {
	ExecuteTransaction(txData *models.Transaction, stockUpdates map[uint]int) error
	FindTransactionsByDateRange(startDate, endDate time.Time) ([]models.Transaction, error)
	FindByID(id uint) (*models.Transaction, error)
	ApproveAndUpdateStock(tx *models.Transaction, stockUpdates map[uint]int) error
	FindOutboundItemsByDate(startDate, endDate time.Time) ([]models.TransactionItem, error)
	FindTransactionsPaginatedByDate(startDate, endDate time.Time, limit int, offset int) ([]models.Transaction, error)
	AnalyzeProductMovement(startDate, endDate time.Time) ([]ProductMovementResult, error)
	ExecuteMultiWarehouseMutation(tx *models.Transaction, stockMutations map[string]int) error
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

func (r *transactionRepository) FindTransactionsPaginatedByDate(startDate, endDate time.Time, limit int, offset int) ([]models.Transaction, error) {
	var transactions []models.Transaction

	err := r.db.Preload("Partner").Preload("Items").Where("transaction_date >= ? AND transaction_date <= ?", startDate, endDate).
		Limit(limit).
		Offset(offset).
		Find(&transactions).Error

	return transactions, err
}

func (r *transactionRepository) AnalyzeProductMovement(startDate, endDate time.Time) ([]ProductMovementResult, error) {
	var results []ProductMovementResult

	err := r.db.Model(&models.TransactionItem{}).Select(`
			products.id as product_id, 
			products.name as product_name, 
			SUM(CASE WHEN transactions.type = 'INBOUND' THEN transaction_items.quantity ELSE 0 END) as total_in,
			SUM(CASE WHEN transactions.type = 'OUTBOUND' THEN transaction_items.quantity ELSE 0 END) as total_out,
			SUM(transaction_items.quantity * transaction_items.unit_price) as total_value`).
		Joins("JOIN transactions ON transactions.id = transaction_items.transaction_id").
		Joins("JOIN products ON products.id = transaction_items.product_id").
		Where("transactions.status = ?", "approved").
		Where("transactions.transaction_date >= ? AND transactions.transaction_date <= ?", startDate, endDate).
		Group("products.id, products.name").
		Order("total_out DESC").
		Scan(&results).Error

	return results, err
}

func (r *transactionRepository) ExecuteMultiWarehouseMutation(tx *models.Transaction, stockMutations map[string]int) error {
	return r.db.Transaction(func(dbTx *gorm.DB) error {
		if err := dbTx.Save(tx).Error; err != nil {
			return err
		}

		for key, qtyChange := range stockMutations {
			parts := strings.Split(key, "-")
			if len(parts) != 2 {
				return fmt.Errorf("format kunci mutasi tidak valid: %s", key)
			}

			whID, _ := strconv.ParseUint(parts[0], 10, 32)
			prodID, _ := strconv.ParseUint(parts[1], 10, 32)

			var whStock models.WarehouseStock
			err := dbTx.Where("warehouse_id = ? AND product_id = ?", whID, prodID).First(&whStock).Error
			if err != nil {
				if err == gorm.ErrRecordNotFound {
					if qtyChange < 0 {
						return fmt.Errorf("anomali: mencoba mengurangi stok barang yang tidak ada di gudang %d", whID)
					}

					whStock = models.WarehouseStock{
						WarehouseID: uint(whID),
						ProductID:   uint(prodID),
						Stock:       qtyChange,
					}
					if err := dbTx.Create(&whStock).Error; err != nil {
						return err
					}
				} else {
					return err
				}
			} else {
				whStock.Stock += qtyChange

				// Validasi lapis terakhir di level database
				if whStock.Stock < 0 {
					return fmt.Errorf("stok negatif tidak diizinkan untuk produk %d di gudang %d", prodID, whID)
				}

				if err := dbTx.Save(&whStock).Error; err != nil {
					return err
				}
			}

			if err := dbTx.Model(&models.Product{}).
				Where("id = ?", prodID).
				Update("current_stock", gorm.Expr("current_stock + ?", qtyChange)).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
