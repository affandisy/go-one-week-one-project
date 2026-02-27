package services

import (
	"errors"
	"strings"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type CreatePartnerRequest struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

type PartnerService interface {
	CreatePartner(req CreatePartnerRequest) (*models.Partner, error)
	GetAllPartners() ([]models.Partner, error)
}

type partnerService struct {
	repo repositories.PartnerRepository
}

func NewPartnerService(repo repositories.PartnerRepository) PartnerService {
	return &partnerService{repo: repo}
}

func (s *partnerService) CreatePartner(req CreatePartnerRequest) (*models.Partner, error) {
	if req.Name == "" {
		return nil, errors.New("Nama partner wajid diisi")
	}

	req.Type = strings.ToUpper(req.Type)
	if req.Type != "SUPPLIER" && req.Type != "CUSTOMER" {
		return nil, errors.New("Tipe partner harus SUPPLIER atau CUSTOMER")
	}

	partner := models.Partner{
		Name:    req.Name,
		Type:    req.Type,
		Phone:   req.Phone,
		Address: req.Address,
	}

	if err := s.repo.Create(&partner); err != nil {
		return nil, errors.New("Gagal menyimpan data partner")
	}

	return &partner, nil
}

func (s *partnerService) GetAllPartners() ([]models.Partner, error) {
	return s.repo.FindAll()
}
