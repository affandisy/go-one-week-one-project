package repositories

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"gorm.io/gorm"
)

type PartnerRepository interface {
	Create(partner *models.Partner) error
	FindAll() ([]models.Partner, error)
	FindByID(id uint) (*models.Partner, error)
}

type partnerRepository struct {
	db *gorm.DB
}

func NewPartnerRepository(db *gorm.DB) PartnerRepository {
	return &partnerRepository{db: db}
}

func (r *partnerRepository) Create(partner *models.Partner) error {
	return r.db.Create(partner).Error
}

func (r *partnerRepository) FindAll() ([]models.Partner, error) {
	var partners []models.Partner

	err := r.db.Find(&partners).Error

	return partners, err
}

func (r *partnerRepository) FindByID(id uint) (*models.Partner, error) {
	var partner models.Partner

	err := r.db.First(&partner, id).Error

	return &partner, err
}
