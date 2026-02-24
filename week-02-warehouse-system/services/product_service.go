package services

import (
	"errors"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type CreateProductRequest struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Category    string  `json:"category"`
	Unit        string  `json:"unit"`
	Price       float64 `json:"price"`
	MinStock    int     `json:"min_stock"`
	MaxStock    int     `json:"max_stock"`
}

type ProductService interface {
	CreateProduct(req CreateProductRequest) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) CreateProduct(req CreateProductRequest) (*models.Product, error) {
	if req.SKU == "" || req.Name == "" || req.Unit == "" {
		return nil, errors.New("SKU, Nama, dan Unit tidak boleh kosong")
	}

	if req.Price < 0 {
		return nil, errors.New("Harga tidak boleh negatif")
	}

	product := models.Product{
		SKU:          req.SKU,
		Name:         req.Name,
		Description:  req.Description,
		Category:     req.Category,
		Unit:         req.Unit,
		Price:        req.Price,
		MinStock:     req.MinStock,
		MaxStock:     req.MaxStock,
		CurrentStock: 0,    // Aturan Bisnis: Barang baru stoknya selalu 0
		IsActive:     true, // Aturan Bisnis: Barang baru otomatis aktif
	}

	err := s.repo.Create(&product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

func (s *productService) GetAllProducts() ([]models.Product, error) {
	return s.repo.FindAll()
}
