package services

import (
	"errors"

	"github.com/affandi/belajar-bahasa/models"
	"github.com/affandi/belajar-bahasa/repositories"
)

type ModuleService interface {
	GetAllModules() ([]models.Module, error)
	GetModuleByID(id string) (*models.Module, error)
	CreateModule(title, description string, levelOrder int) error
}

type moduleService struct {
	repo repositories.ModuleRepository
}

func NewModuleService(repo repositories.ModuleRepository) ModuleService {
	return &moduleService{repo}
}

func (s *moduleService) GetAllModules() ([]models.Module, error) {
	return s.repo.FindAll()
}

func (s *moduleService) GetModuleByID(id string) (*models.Module, error) {
	return s.repo.FindByID(id)
}

func (s *moduleService) CreateModule(title, description string, levelOrder int) error {
	if title == "" || levelOrder <= 0 {
		return errors.New("title dan level_order tidak valid")
	}

	// Aturan PRD: Hanya Level 1 yang terbuka di awal (IsLocked = false)
	isLocked := true
	if levelOrder == 1 {
		isLocked = false
	}

	module := &models.Module{
		Title:       title,
		Description: description,
		LevelOrder:  levelOrder,
		IsLocked:    isLocked,
	}

	return s.repo.Create(module)
}
