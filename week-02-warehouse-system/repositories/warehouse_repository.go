package repositories

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"gorm.io/gorm"
)

type WarehouseRepository interface {
	Create(warehouse *models.Warehouse) error
	FindAll() ([]models.Warehouse, error)
	FindByID(id uint) (*models.Warehouse, error)
}

type warehouseRepository struct {
	db *gorm.DB
}

func NewWarehouseRepository(db *gorm.DB) WarehouseRepository {
	return &warehouseRepository{db: db}
}

func (r *warehouseRepository) Create(warehouse *models.Warehouse) error {
	return r.db.Create(warehouse).Error
}
func (r *warehouseRepository) FindAll() ([]models.Warehouse, error) {
	var warehouses []models.Warehouse

	err := r.db.Find(&warehouses).Error

	return warehouses, err
}

func (r *warehouseRepository) FindByID(id uint) (*models.Warehouse, error) {
	var warehouse models.Warehouse

	err := r.db.First(&warehouse, id).Error

	return &warehouse, err
}
