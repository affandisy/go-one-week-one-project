package repository

import (
	"github.com/affandisy/school-system-sma/models"
	"gorm.io/gorm"
)

type AcademicRepository interface {
	// Academic Year
	CreateAcademicYear(ay *models.AcademicYear) error
	FindAllAcademicYears() ([]models.AcademicYear, error)
	UpdateAcademicYearsStatus(tx *gorm.DB, status bool) error

	// Transaction support
	BeginTransaction() *gorm.DB
}

type academicRepository struct {
	db *gorm.DB
}

func NewAcademicRepository(db *gorm.DB) AcademicRepository {
	return &academicRepository{db}
}

func (r *academicRepository) CreateAcademicYear(ay *models.AcademicYear) error {
	return r.db.Create(ay).Error
}

func (r *academicRepository) FindAllAcademicYears() ([]models.AcademicYear, error) {
	var ays []models.AcademicYear
	err := r.db.Order("start_date DESC").Find(&ays).Error
	return ays, err
}

// Untuk mengubah status seluruh tahun ajaran
func (r *academicRepository) UpdateAcademicYearsStatus(tx *gorm.DB, status bool) error {
	return tx.Model(&models.AcademicYear{}).Where("1 = 1").Update("is_active", status).Error
}

func (r *academicRepository) BeginTransaction() *gorm.DB {
	return r.db.Begin()
}
