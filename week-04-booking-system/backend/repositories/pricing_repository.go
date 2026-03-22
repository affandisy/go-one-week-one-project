package repositories

import (
	"github.com/affandisy/padel-booking-system/models"
	"gorm.io/gorm"
)

type PricingRepository interface {
	Create(rule *models.PricingRule) error
	FindByCourtID(courtID string) ([]models.PricingRule, error)
	FindByID(id string) (*models.PricingRule, error)
	Update(rule *models.PricingRule) error
	Delete(id string) error
}

type pricingRepository struct {
	db *gorm.DB
}

func NewPricingRepository(db *gorm.DB) PricingRepository {
	return &pricingRepository{db}
}

func (r *pricingRepository) Create(rule *models.PricingRule) error {
	return r.db.Create(rule).Error
}

// Biasanya kita ingin melihat aturan harga spesifik untuk satu lapangan tertentu
func (r *pricingRepository) FindByCourtID(courtID string) ([]models.PricingRule, error) {
	var rules []models.PricingRule
	err := r.db.Where("court_id = ?", courtID).Order("day_type, start_time").Find(&rules).Error
	return rules, err
}

func (r *pricingRepository) FindByID(id string) (*models.PricingRule, error) {
	var rule models.PricingRule
	err := r.db.Where("id = ?", id).First(&rule).Error
	return &rule, err
}

func (r *pricingRepository) Update(rule *models.PricingRule) error {
	return r.db.Save(rule).Error
}

func (r *pricingRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.PricingRule{}).Error
}
