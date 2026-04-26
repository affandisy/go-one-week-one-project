package repositories

import (
	"github.com/affandisy/pos-system/models"
	"gorm.io/gorm"
)

type SupplierRepository interface {
	FindAll() ([]models.Supplier, error)
	FindByID(id string) (*models.Supplier, error)
	Create(supplier *models.Supplier) error
	Update(supplier *models.Supplier) error
	Delete(id string) error
}

type supplierRepository struct {
	db *gorm.DB
}

func NewSupplierRepository(db *gorm.DB) SupplierRepository {
	return &supplierRepository{db}
}

func (r *supplierRepository) FindAll() ([]models.Supplier, error) {
	var suppliers []models.Supplier
	err := r.db.Order("name asc").Find(&suppliers).Error
	return suppliers, err
}

func (r *supplierRepository) FindByID(id string) (*models.Supplier, error) {
	var supplier models.Supplier
	err := r.db.Where("id = ?", id).First(&supplier).Error
	return &supplier, err
}

func (r *supplierRepository) Create(supplier *models.Supplier) error {
	return r.db.Create(supplier).Error
}

func (r *supplierRepository) Update(supplier *models.Supplier) error {
	return r.db.Save(supplier).Error
}

func (r *supplierRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Supplier{}).Error
}
