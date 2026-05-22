package postgres

import (
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"gorm.io/gorm"
)

// DTO GORM
type NutritionLogModel struct {
	ID          string    `gorm:"type:uuid;primaryKey"`
	PetID       string    `gorm:"type:uuid;index"`
	LogDate     time.Time `gorm:"index"`
	FoodBrand   string    // Pencatatan spesifik diet/makanan berkualitas tinggi
	Calories    int
	HealthNotes string // Pencatatan detail kondisi fisik
}

func (NutritionLogModel) TableName() string { return "nutrition_logs" }

type nutritionRepository struct {
	db *gorm.DB
}

func NewNutritionRepository(db *gorm.DB) port.NutritionRepository {
	db.AutoMigrate(&NutritionLogModel{})
	return &nutritionRepository{db}
}

func (r *nutritionRepository) SaveLog(log *domain.NutritionLog) error {
	model := NutritionLogModel{
		ID:          log.ID,
		PetID:       log.PetID,
		LogDate:     log.LogDate,
		FoodBrand:   log.FoodBrand,
		Calories:    log.Calories,
		HealthNotes: log.HealthNotes,
	}
	return r.db.Create(&model).Error
}

func (r *nutritionRepository) GetLogsByPet(petID string) ([]domain.NutritionLog, error) {
	var models []NutritionLogModel
	// Mengurutkan dari catatan paling baru ke lama
	if err := r.db.Where("pet_id = ?", petID).Order("log_date desc").Find(&models).Error; err != nil {
		return nil, err
	}

	var logs []domain.NutritionLog
	for _, m := range models {
		logs = append(logs, domain.NutritionLog{
			ID:          m.ID,
			PetID:       m.PetID,
			LogDate:     m.LogDate,
			FoodBrand:   m.FoodBrand,
			Calories:    m.Calories,
			HealthNotes: m.HealthNotes,
		})
	}
	return logs, nil
}
