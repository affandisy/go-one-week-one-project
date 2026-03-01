package services

import (
	"errors"
	"strings"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type CreateLocationRequest struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type LocationService interface {
	CreateLocation(req CreateLocationRequest) (*models.Location, error)
	GetAllLocations() ([]models.Location, error)
}

type locationService struct {
	repo repositories.LocationRepository
}

func NewLocationService(repo repositories.LocationRepository) LocationService {
	return &locationService{repo: repo}
}

func (s *locationService) CreateLocation(req CreateLocationRequest) (*models.Location, error) {
	if req.Code == "" || req.Name == "" {
		return nil, errors.New("Kode dan Nama lokasi wajib diisi")
	}

	req.Code = strings.ToUpper(strings.TrimSpace(req.Code))
	req.Code = strings.ReplaceAll(req.Code, " ", "-")

	existing, _ := s.repo.FindByCode(req.Code)
	if existing != nil && existing.ID != 0 {
		return nil, errors.New("Kode lokasi ini sudah digunakan")
	}

	location := models.Location{
		Code:        req.Code,
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.repo.Create(&location); err != nil {
		return nil, errors.New("Gagal menyimpan data lokasi")
	}

	return &location, nil
}

func (s *locationService) GetAllLocations() ([]models.Location, error) {
	return s.repo.FindAll()
}
