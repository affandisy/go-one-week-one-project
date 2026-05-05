// repositories/master_repository.go
package repositories

import (
	"github.com/affandisy/financial-app/models"
	"gorm.io/gorm"
)

type MasterRepository interface {
	GetDefaultWallet() (*models.Wallet, error)
	GetCategoriesByType(txType string) ([]models.Category, error)
}

type masterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) MasterRepository {
	return &masterRepository{db}
}

func (r *masterRepository) GetDefaultWallet() (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.First(&wallet).Error
	return &wallet, err
}

func (r *masterRepository) GetCategoriesByType(txType string) ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Where("type = ? AND is_active = ?", txType, true).Find(&categories).Error
	return categories, err
}
