package services

import (
	"errors"

	"github.com/affandisy/pos-system/models"
	"github.com/affandisy/pos-system/repositories"
)

type SupplierRequest struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
	Notes string `json:"notes"`
}

type SupplierService interface {
	GetAllSuppliers() ([]models.Supplier, error)
	GetSupplierByID(id string) (*models.Supplier, error)
	CreateSupplier(req SupplierRequest) error
	UpdateSupplier(id string, req SupplierRequest) error
	DeleteSupplier(id string) error
}

type supplierService struct {
	repo repositories.SupplierRepository
}

func NewSupplierService(repo repositories.SupplierRepository) SupplierService {
	return &supplierService{repo}
}

func (s *supplierService) GetAllSuppliers() ([]models.Supplier, error) {
	return s.repo.FindAll()
}

func (s *supplierService) GetSupplierByID(id string) (*models.Supplier, error) {
	return s.repo.FindByID(id)
}

func (s *supplierService) CreateSupplier(req SupplierRequest) error {
	if req.Name == "" {
		return errors.New("nama pemasok tidak boleh kosong")
	}

	supplier := &models.Supplier{
		Name:  req.Name,
		Phone: req.Phone,
		Notes: req.Notes,
	}

	return s.repo.Create(supplier)
}

func (s *supplierService) UpdateSupplier(id string, req SupplierRequest) error {
	supplier, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("pemasok tidak ditemukan")
	}

	if req.Name != "" {
		supplier.Name = req.Name
	}
	supplier.Phone = req.Phone
	supplier.Notes = req.Notes

	return s.repo.Update(supplier)
}

func (s *supplierService) DeleteSupplier(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("pemasok tidak ditemukan")
	}
	return s.repo.Delete(id)
}
