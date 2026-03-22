package repositories

import (
	"github.com/affandisy/padel-booking-system/models"
	"gorm.io/gorm"
)

type CourtRepository interface {
	Create(court *models.Court) error
	FindAll() ([]models.Court, error)
	FindByID(id string) (*models.Court, error)
	Update(court *models.Court) error
	Delete(id string) error
}

type courtRepository struct {
	db *gorm.DB
}

func NewCourtRepository(db *gorm.DB) CourtRepository {
	return &courtRepository{db}
}

func (r *courtRepository) Create(court *models.Court) error {
	return r.db.Create(court).Error
}

func (r *courtRepository) FindAll() ([]models.Court, error) {
	var courts []models.Court
	// Mengurutkan berdasarkan yang terbaru dibuat
	err := r.db.Order("created_at desc").Find(&courts).Error
	return courts, err
}

func (r *courtRepository) FindByID(id string) (*models.Court, error) {
	var court models.Court
	err := r.db.Where("id = ?", id).First(&court).Error
	return &court, err
}

func (r *courtRepository) Update(court *models.Court) error {
	return r.db.Save(court).Error
}

func (r *courtRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Court{}).Error
}
