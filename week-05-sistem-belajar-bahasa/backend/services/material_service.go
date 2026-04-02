package services

import (
	"errors"

	"github.com/affandi/belajar-bahasa/models"
	"github.com/affandi/belajar-bahasa/repositories"
	"github.com/google/uuid"
)

type MaterialService interface {
	GetMaterialsByModule(moduleID string) ([]models.Material, error)
	CreateMaterial(req CreateMaterialRequest) error
}

type materialService struct {
	repo repositories.MaterialRepository
}

func NewMaterialService(repo repositories.MaterialRepository) MaterialService {
	return &materialService{repo}
}

type CreateMaterialRequest struct {
	ModuleID      string `json:"module_id"`
	ContentType   string `json:"content_type"` // learn_card, quiz_mcq, quiz_match, quiz_unscramble
	Question      string `json:"question"`
	CorrectAnswer string `json:"correct_answer"`
	Options       string `json:"options"` // Array JSON string
	ImageURL      string `json:"image_url"`
	AudioURL      string `json:"audio_url"`
	DisplayOrder  int    `json:"display_order"`
}

func (s *materialService) GetMaterialsByModule(moduleID string) ([]models.Material, error) {
	return s.repo.GetByModuleID(moduleID)
}

func (s *materialService) CreateMaterial(req CreateMaterialRequest) error {
	if req.ModuleID == "" || req.ContentType == "" || req.Question == "" || req.CorrectAnswer == "" {
		return errors.New("data materi tidak lengkap")
	}

	modID, err := uuid.Parse(req.ModuleID)
	if err != nil {
		return errors.New("format Module ID tidak valid")
	}

	material := &models.Material{
		ModuleID:      modID,
		ContentType:   req.ContentType,
		Question:      req.Question,
		CorrectAnswer: req.CorrectAnswer,
		Options:       req.Options,
		ImageURL:      req.ImageURL,
		AudioURL:      req.AudioURL,
		DisplayOrder:  req.DisplayOrder,
	}

	return s.repo.Create(material)
}
