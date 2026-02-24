package repositories

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"gorm.io/gorm"
)

type ProductRepository interface {
	Create(product *models.Product) error
	FindAll() ([]models.Product, error)
	FindByID(id uint) (*models.Product, error)
	Update(product *models.Product) error
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) Create(product *models.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) FindAll() ([]models.Product, error) {
	var products []models.Product
	err := r.db.Find(&products).Error

	return products, err
}

func (r *productRepository) FindByID(id uint) (*models.Product, error) {
	var product models.Product
	err := r.db.First(&product, id).Error

	return &product, err
}

func (r *productRepository) Update(product *models.Product) error {
	return r.db.Save(product).Error
}
