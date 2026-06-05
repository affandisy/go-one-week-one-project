package port

import "github.com/affandisy/petcare-system/internal/core/domain"

// Driven Port (Database)
type AnalyticsRepository interface {
	GetNutritionSummaryByPet(petID string) (*domain.PetNutritionSummary, error)
}

// Driving Port (API)
type AnalyticsUseCase interface {
	GeneratePetNutritionReport(petID string) (*domain.PetNutritionSummary, error)
}
