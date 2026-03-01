package repositories

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"gorm.io/gorm"
)

type LocationRepository interface {
	Create(location *models.Location) error
	FindAll() ([]models.Location, error)
	FindByCode(code string) (*models.Location, error)
}

type locationRepository struct {
	db *gorm.DB
}

func NewLocationRepository(db *gorm.DB) LocationRepository {
	return &locationRepository{db: db}
}

func (r *locationRepository) Create(location *models.Location) error {
	return r.db.Create(location).Error
}

func (r *locationRepository) FindAll() ([]models.Location, error) {
	var locations []models.Location

	err := r.db.Find(&locations).Error

	return locations, err
}

func (r *locationRepository) FindByCode(code string) (*models.Location, error) {
	var location models.Location

	err := r.db.Where("code = ?", code).First(&location).Error

	return &location, err
}
