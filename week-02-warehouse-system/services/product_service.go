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

type PaginatedResponse struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	TotalPages int         `json:"total_pages"`
}

type ProductService interface {
	CreateProduct(req CreateProductRequest) (*models.Product, error)
	GetAllProducts() ([]models.Product, error)
	GetAllProductsPaginated(page int, limit int) (PaginatedResponse, error)
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo: repo}
}

func (s *productService) GetAllProductsPaginated(page int, limit int) (PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}

	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	products, total, err := s.repo.FindAllPaginated(limit, offset)
	if err != nil {
		return PaginatedResponse{}, err
	}

	totalPages := int(total) / limit
	if int(total)%limit != 0 {
		totalPages++
	}

	return PaginatedResponse{
		Data:       products,
		Total:      total,
		Page:       page,
		TotalPages: totalPages,
	}, nil
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
