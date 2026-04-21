package repositories

import (
	"github.com/affandisy/pos-system/models"
	"gorm.io/gorm"
)

type StockRepository interface {
	AddStockIn(movement *models.StockMovement) error
	GetMovementsByProduct(productID string) ([]models.StockMovement, error)
}

type stockRepository struct {
	db *gorm.DB
}

func NewStockRepository(db *gorm.DB) StockRepository {
	return &stockRepository{db}
}

// AddStockIn mencatat riwayat masuk dan menambah stok produk secara atomik (bersamaan)
func (r *stockRepository) AddStockIn(movement *models.StockMovement) error {
	// Memulai Database Transaction
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Simpan riwayat pergerakan (Stock Movement)
		if err := tx.Create(movement).Error; err != nil {
			return err
		}

		// 2. Tambahkan kuantitas ke stok Produk yang ada saat ini
		if err := tx.Model(&models.Product{}).
			Where("id = ?", movement.ProductID).
			UpdateColumn("stock", gorm.Expr("stock + ?", movement.Quantity)).Error; err != nil {
			return err
		}

		return nil // Jika semua lancar, transaksi di-commit
	})
}

func (r *stockRepository) GetMovementsByProduct(productID string) ([]models.StockMovement, error) {
	var movements []models.StockMovement
	err := r.db.Where("product_id = ?", productID).Order("created_at desc").Find(&movements).Error
	return movements, err
}
