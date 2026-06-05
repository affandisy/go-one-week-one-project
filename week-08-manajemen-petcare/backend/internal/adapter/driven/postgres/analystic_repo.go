package postgres

import (
	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"gorm.io/gorm"
)

type analyticsRepository struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) port.AnalyticsRepository {
	return &analyticsRepository{db}
}

func (r *analyticsRepository) GetNutritionSummaryByPet(petID string) (*domain.PetNutritionSummary, error) {
	var summary domain.PetNutritionSummary

	// 1. Kueri untuk mengambil rata-rata kalori dan total pencatatan
	type AggregationResult struct {
		TotalLogs       int
		AverageCalories float64
	}
	var agg AggregationResult

	err := r.db.Model(&NutritionLogModel{}).
		Select("COUNT(id) as total_logs, COALESCE(AVG(calories), 0) as average_calories").
		Where("pet_id = ?", petID).
		Scan(&agg).Error

	if err != nil {
		return nil, err
	}

	// 2. Kueri terpisah untuk mengambil catatan kesehatan terbaru dan merek makanan
	var latestLog NutritionLogModel
	r.db.Where("pet_id = ?", petID).Order("log_date desc").First(&latestLog)

	var pet PetModel
	r.db.Where("id = ?", petID).First(&pet)

	summary.PetID = petID
	summary.PetName = pet.Name
	summary.TotalLogs = agg.TotalLogs
	summary.AverageCalories = agg.AverageCalories
	summary.MostUsedBrand = latestLog.FoodBrand
	summary.LatestHealthNote = latestLog.HealthNotes

	return &summary, nil
}
