package services

import (
	"errors"

	"github.com/affandisy/pos-system/models"
	"github.com/affandisy/pos-system/repositories"
)

type ProductRequest struct {
	Barcode       string  `json:"barcode"`
	Name          string  `json:"name"`
	PurchasePrice float64 `json:"purchase_price"`
	SellingPrice  float64 `json:"selling_price"`
	Unit          string  `json:"unit"`
	MinStock      int     `json:"min_stock"`
}

type ProductService interface {
	GetAllProducts(search string) ([]models.Product, error)
	GetProductByID(id string) (*models.Product, error)
	CreateProduct(req ProductRequest) error
	UpdateProduct(id string, req ProductRequest) error
	DeleteProduct(id string) error
}

type productService struct {
	repo repositories.ProductRepository
}

func NewProductService(repo repositories.ProductRepository) ProductService {
	return &productService{repo}
}

func (s *productService) GetAllProducts(search string) ([]models.Product, error) {
	return s.repo.FindAll(search)
}

func (s *productService) GetProductByID(id string) (*models.Product, error) {
	return s.repo.FindByID(id)
}

func (s *productService) CreateProduct(req ProductRequest) error {
	if req.Name == "" || req.SellingPrice <= 0 || req.Unit == "" {
		return errors.New("nama, harga jual, dan satuan wajib diisi")
	}

	// Cek duplikasi barcode jika barcode diisi
	if req.Barcode != "" {
		existing, _ := s.repo.FindByBarcode(req.Barcode)
		if existing != nil {
			return errors.New("barcode sudah digunakan oleh produk lain")
		}
	}

	product := &models.Product{
		Barcode:       req.Barcode,
		Name:          req.Name,
		PurchasePrice: req.PurchasePrice,
		SellingPrice:  req.SellingPrice,
		Unit:          req.Unit,
		Stock:         0, // Stok awal selalu 0, harus diisi via Barang Masuk
		MinStock:      req.MinStock,
	}

	return s.repo.Create(product)
}

func (s *productService) UpdateProduct(id string, req ProductRequest) error {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}

	product.Barcode = req.Barcode
	product.Name = req.Name
	product.PurchasePrice = req.PurchasePrice
	product.SellingPrice = req.SellingPrice
	product.Unit = req.Unit
	product.MinStock = req.MinStock

	return s.repo.Update(product)
}

func (s *productService) DeleteProduct(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("produk tidak ditemukan")
	}
	return s.repo.Delete(id)
}
