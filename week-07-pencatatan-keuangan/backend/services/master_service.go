// services/master_service.go
package services

import (
	"github.com/affandisy/financial-app/models"
	"github.com/affandisy/financial-app/repositories"
)

type MasterService interface {
	GetWalletInfo() (*models.Wallet, error)
	GetActiveCategories(txType string) ([]models.Category, error)
}

type masterService struct {
	repo repositories.MasterRepository
}

func NewMasterService(repo repositories.MasterRepository) MasterService {
	return &masterService{repo}
}

func (s *masterService) GetWalletInfo() (*models.Wallet, error) {
	return s.repo.GetDefaultWallet()
}

func (s *masterService) GetActiveCategories(txType string) ([]models.Category, error) {
	return s.repo.GetCategoriesByType(txType)
}
