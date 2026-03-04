package repositories

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"gorm.io/gorm"
)

type WarehouseStockRepository interface {
	GetStock(warehouseID, productID uint) (*models.WarehouseStock, error)
	UpsertStock(warehouseID, productID uint, qty int, tx *gorm.DB) error
}

type warehouseStockRepository struct {
	db *gorm.DB
}

func NewWarehouseStockRepository(db *gorm.DB) WarehouseStockRepository {
	return &warehouseStockRepository{db: db}
}

func (r *warehouseStockRepository) GetStock(warehouseID, productID uint) (*models.WarehouseStock, error) {
	var stock models.WarehouseStock

	err := r.db.Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).First(&stock).Error

	return &stock, err
}

func (r *warehouseStockRepository) UpsertStock(warehouseID, productID uint, qty int, tx *gorm.DB) error {
	var stock models.WarehouseStock

	err := tx.Where("warehouse_id = ? AND product_id = ?", warehouseID, productID).First(&stock).Error

	if err != nil && err == gorm.ErrRecordNotFound {
		stock = models.WarehouseStock{
			WarehouseID: warehouseID,
			ProductID:   productID,
			Stock:       qty,
		}
		return tx.Create(&stock).Error
	} else if err != nil {
		return err
	}

	stock.Stock += qty
	return tx.Save(&stock).Error
}
