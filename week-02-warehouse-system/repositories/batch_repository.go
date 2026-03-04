package repositories

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"gorm.io/gorm"
)

type BatchRepository interface {
	FindByProductID(productID uint) ([]models.ProductBatch, error)
	UpdateBatch(batch *models.ProductBatch) error
}

type batchRepository struct {
	db *gorm.DB
}

func NewBatchRepository(db *gorm.DB) BatchRepository {
	return &batchRepository{db: db}
}

func (r *batchRepository) FindByProductID(productID uint) ([]models.ProductBatch, error) {
	var batches []models.ProductBatch

	err := r.db.Where("product_id = ? AND stock > 0", productID).Find(&batches).Error

	return batches, err
}

func (r *batchRepository) UpdateBatch(batch *models.ProductBatch) error {
	return r.db.Save(batch).Error
}
