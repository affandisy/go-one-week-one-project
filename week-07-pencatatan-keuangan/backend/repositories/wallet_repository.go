package repositories

import (
	"github.com/affandisy/financial-app/models"
	"gorm.io/gorm"
)

type WalletRepository interface {
	FindAll() ([]models.Wallet, error)
	FindByID(id string) (*models.Wallet, error)
	Create(wallet *models.Wallet) error
	Update(wallet *models.Wallet) error
	Delete(id string) error
	Count() (int64, error)
}

type walletRepository struct {
	db *gorm.DB
}

func NewWalletRepository(db *gorm.DB) WalletRepository {
	return &walletRepository{db}
}

func (r *walletRepository) FindAll() ([]models.Wallet, error) {
	var wallets []models.Wallet
	err := r.db.Order("created_at asc").Find(&wallets).Error
	return wallets, err
}

func (r *walletRepository) FindByID(id string) (*models.Wallet, error) {
	var wallet models.Wallet
	err := r.db.First(&wallet, "id = ?", id).Error
	return &wallet, err
}

func (r *walletRepository) Create(wallet *models.Wallet) error {
	return r.db.Create(wallet).Error
}

func (r *walletRepository) Update(wallet *models.Wallet) error {
	return r.db.Save(wallet).Error
}

func (r *walletRepository) Delete(id string) error {
	return r.db.Delete(&models.Wallet{}, "id = ?", id).Error
}

func (r *walletRepository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&models.Wallet{}).Count(&count).Error
	return count, err
}
