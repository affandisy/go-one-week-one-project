package repositories

import (
	"github.com/affandi/belajar-bahasa/models"
	"gorm.io/gorm"
)

type ModuleRepository interface {
	FindAll() ([]models.Module, error)
	FindByID(id string) (*models.Module, error)
	Create(module *models.Module) error
}

type moduleRepository struct {
	db *gorm.DB
}

func NewModuleRepository(db *gorm.DB) ModuleRepository {
	return &moduleRepository{db}
}

func (r *moduleRepository) FindAll() ([]models.Module, error) {
	var modules []models.Module
	// Urutkan berdasarkan level (1, 2, 3)
	err := r.db.Order("level_order asc").Find(&modules).Error
	return modules, err
}

func (r *moduleRepository) FindByID(id string) (*models.Module, error) {
	var module models.Module
	err := r.db.Where("id = ?", id).First(&module).Error
	return &module, err
}

func (r *moduleRepository) Create(module *models.Module) error {
	return r.db.Create(module).Error
}
