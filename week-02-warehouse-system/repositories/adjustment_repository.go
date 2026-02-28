package repositories

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"gorm.io/gorm"
)

type AdjustmentRepository interface {
	ExecuteAdjustment(adjustment *models.StockAdjustment, newStock int) error
	ExecuteMassAdjustment(adjustments []models.StockAdjustment, stockUpdates map[uint]int) error
}

type adjustmentRepository struct {
	db *gorm.DB
}

func NewAdjustmentRepository(db *gorm.DB) AdjustmentRepository {
	return &adjustmentRepository{db: db}
}

func (r *adjustmentRepository) ExecuteAdjustment(adjustment *models.StockAdjustment, newStock int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(adjustment).Error; err != nil {
			return err
		}

		if err := tx.Model(&models.Product{}).Where("id = ?", adjustment.ProductID).Update("current_stock", newStock).Error; err != nil {
			return err
		}

		return nil
	})
}

func (r *adjustmentRepository) ExecuteMassAdjustment(adjustments []models.StockAdjustment, stockUpdates map[uint]int) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if len(adjustments) > 0 {
			if err := tx.Create(&adjustments).Error; err != nil {
				return err
			}
		}

		for productID, newStock := range stockUpdates {
			if err := tx.Model(&models.Product{}).Where("id = ?", productID).Update("current_stock", newStock).Error; err != nil {
				return err
			}
		}

		return nil
	})
}
