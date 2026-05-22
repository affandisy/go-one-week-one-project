package usecase

import (
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/google/uuid"
)

type nutritionUseCase struct {
	repo port.NutritionRepository
}

func NewNutritionUseCase(repo port.NutritionRepository) port.NutritionUseCase {
	return &nutritionUseCase{repo}
}

func (uc *nutritionUseCase) RecordDiet(petID string, brand string, calories int, notes string) error {
	log := &domain.NutritionLog{
		ID:          uuid.NewString(),
		PetID:       petID,
		LogDate:     time.Now(),
		FoodBrand:   brand,
		Calories:    calories,
		HealthNotes: notes,
	}

	return uc.repo.SaveLog(log)
}

func (uc *nutritionUseCase) GetLogsByPet(petID string) ([]domain.NutritionLog, error) {
	return uc.repo.GetLogsByPet(petID)
}
