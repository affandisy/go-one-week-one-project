package repositories

import (
	"github.com/affandi/belajar-bahasa/models"
	"gorm.io/gorm"
)

type MaterialRepository interface {
	GetByModuleID(moduleID string) ([]models.Material, error)
	Create(material *models.Material) error
}

type materialRepository struct {
	db *gorm.DB
}

func NewMaterialRepository(db *gorm.DB) MaterialRepository {
	return &materialRepository{db}
}

func (r *materialRepository) GetByModuleID(moduleID string) ([]models.Material, error) {
	var materials []models.Material
	// Urutkan berdasarkan display_order
	err := r.db.Where("module_id = ?", moduleID).Order("display_order asc").Find(&materials).Error
	return materials, err
}

func (r *materialRepository) Create(material *models.Material) error {
	return r.db.Create(material).Error
}
