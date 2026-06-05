package usecase

import (
	"errors"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
)

type analyticsUseCase struct {
	repo port.AnalyticsRepository
}

func NewAnalyticsUseCase(repo port.AnalyticsRepository) port.AnalyticsUseCase {
	return &analyticsUseCase{repo}
}

func (uc *analyticsUseCase) GeneratePetNutritionReport(petID string) (*domain.PetNutritionSummary, error) {
	if petID == "" {
		return nil, errors.New("ID hewan peliharaan tidak boleh kosong")
	}

	summary, err := uc.repo.GetNutritionSummaryByPet(petID)
	if err != nil {
		return nil, errors.New("gagal menghasilkan laporan analitik gizi")
	}

	return summary, nil
}
