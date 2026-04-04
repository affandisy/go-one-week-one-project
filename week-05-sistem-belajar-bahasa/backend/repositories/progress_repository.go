package repositories

import (
	"github.com/affandi/belajar-bahasa/models"
	"gorm.io/gorm"
)

type ProgressRepository interface {
	GetUserProgress(userID, moduleID string) (*models.UserProgress, error)
	SaveUserProgress(progress *models.UserProgress) error
	GetNextModule(currentLevelOrder int) (*models.Module, error)
	GetAllByUser(userID string) ([]models.UserProgress, error)
	DeleteByUser(userID string) error
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

// Mengambil semua riwayat progres milik satu user
func (r *progressRepository) GetAllByUser(userID string) ([]models.UserProgress, error) {
	var progresses []models.UserProgress
	// Kita tidak perlu memuat seluruh data modul, cukup ID dan skornya
	err := r.db.Where("user_id = ?", userID).Find(&progresses).Error
	return progresses, err
}

// Menghapus seluruh riwayat progres milik satu user (Reset)
func (r *progressRepository) DeleteByUser(userID string) error {
	// Menghapus secara permanen (Hard Delete) untuk tabel progres
	return r.db.Where("user_id = ?", userID).Unscoped().Delete(&models.UserProgress{}).Error
}
