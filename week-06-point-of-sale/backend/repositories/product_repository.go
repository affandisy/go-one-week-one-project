package repositories

import (
	"github.com/affandisy/pos-system/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	FindAll(search string) ([]models.Product, error)
	FindByID(id string) (*models.Product, error)
	FindByBarcode(barcode string) (*models.Product, error)
	Create(product *models.Product) error
	Update(product *models.Product) error
	Delete(id string) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db}
}

func (r *productRepository) FindAll(search string) ([]models.Product, error) {
	var products []models.Product
	query := r.db.Model(&models.Product{})

	if search != "" {
		// Mencari berdasarkan nama ATAU barcode sesuai PRD
		searchParam := "%" + search + "%"
		query = query.Where("name ILIKE ? OR barcode ILIKE ?", searchParam, searchParam)
	}

	err := query.Order("name asc").Find(&products).Error
	return products, err
}

func (r *productRepository) FindByID(id string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("id = ?", id).First(&product).Error
	return &product, err
}

func (r *productRepository) FindByBarcode(barcode string) (*models.Product, error) {
	var product models.Product
	err := r.db.Where("barcode = ?", barcode).First(&product).Error
	return &product, err
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&models.Product{}).Error
}
