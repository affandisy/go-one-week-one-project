package repositories

import (
	"github.com/affandi/belajar-bahasa/models"
	"gorm.io/gorm"
)

type ProgressRepository interface {
	GetUserProgress(userID, moduleID string) (*models.UserProgress, error)
	SaveUserProgress(progress *models.UserProgress) error
	GetNextModule(currentLevelOrder int) (*models.Module, error)
}

type progressRepository struct {
	db *gorm.DB
}

func NewProgressRepository(db *gorm.DB) ProgressRepository {
	return &progressRepository{db}
}

// Mengambil progres spesifik user pada suatu modul
func (r *progressRepository) GetUserProgress(userID, moduleID string) (*models.UserProgress, error) {
	var progress models.UserProgress
	err := r.db.Where("user_id = ? AND module_id = ?", userID, moduleID).First(&progress).Error
	return &progress, err
}

// Menyimpan atau memperbarui progres
func (r *progressRepository) SaveUserProgress(progress *models.UserProgress) error {
	return r.db.Save(progress).Error
}

// Mencari modul selanjutnya berdasarkan level order
func (r *progressRepository) GetNextModule(currentLevelOrder int) (*models.Module, error) {
	var module models.Module
	err := r.db.Where("level_order = ?", currentLevelOrder+1).First(&module).Error
	return &module, err
}
