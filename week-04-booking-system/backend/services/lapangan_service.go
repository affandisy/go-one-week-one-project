package services

import (
	"errors"

	"github.com/affandisy/padel-booking-system/models"
	"github.com/affandisy/padel-booking-system/repositories"
)

type CourtService interface {
	CreateCourt(req CreateCourtRequest) error
	GetAllCourts() ([]models.Court, error)
	GetCourtByID(id string) (*models.Court, error)
	UpdateCourt(id string, req UpdateCourtRequest) error
	DeleteCourt(id string) error
}

type courtService struct {
	repo repositories.CourtRepository
}

func NewCourtService(repo repositories.CourtRepository) CourtService {
	return &courtService{repo}
}

// ================= DTO (Data Transfer Objects) =================

type CreateCourtRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`   // indoor, outdoor
	Status      string `json:"status"` // active, inactive, maintenance
	Description string `json:"description"`
}

type UpdateCourtRequest struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Status      string `json:"status"`
	Description string `json:"description"`
}

// ================= IMPLEMENTASI LOGIKA =================

func (s *courtService) CreateCourt(req CreateCourtRequest) error {
	if req.Name == "" || req.Type == "" {
		return errors.New("nama dan tipe lapangan wajib diisi")
	}

	if req.Status == "" {
		req.Status = "active" // Default value
	}

	court := &models.Court{
		Name:        req.Name,
		Type:        req.Type,
		Status:      req.Status,
		Description: req.Description,
	}

	return s.repo.Create(court)
}

func (s *courtService) GetAllCourts() ([]models.Court, error) {
	return s.repo.FindAll()
}

func (s *courtService) GetCourtByID(id string) (*models.Court, error) {
	if id == "" {
		return nil, errors.New("ID lapangan tidak valid")
	}
	return s.repo.FindByID(id)
}

func (s *courtService) UpdateCourt(id string, req UpdateCourtRequest) error {
	// 1. Cari data lama
	court, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("lapangan tidak ditemukan")
	}

	// 2. Update field yang diubah
	if req.Name != "" {
		court.Name = req.Name
	}
	if req.Type != "" {
		court.Type = req.Type
	}
	if req.Status != "" {
		court.Status = req.Status
	}
	if req.Description != "" {
		court.Description = req.Description
	}

	// 3. Simpan perubahan
	return s.repo.Update(court)
}

func (s *courtService) DeleteCourt(id string) error {
	// Pastikan data ada sebelum dihapus
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("lapangan tidak ditemukan")
	}

	return s.repo.Delete(id)
}
