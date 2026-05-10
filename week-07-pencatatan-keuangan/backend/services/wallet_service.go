package services

import (
	"errors"
	"strings"

	"github.com/affandisy/financial-app/models"
	"github.com/affandisy/financial-app/repositories"
	"github.com/google/uuid"
)

type WalletRequest struct {
	Name           string  `json:"name"`
	InitialBalance float64 `json:"initial_balance"`
}

type WalletService interface {
	GetAllWallets() ([]models.Wallet, error)
	GetWallet(id string) (*models.Wallet, error)
	CreateWallet(req WalletRequest) error
	UpdateWallet(id string, req WalletRequest) error
	DeleteWallet(id string) error
}

type walletService struct {
	repo repositories.WalletRepository
}

func NewWalletService(repo repositories.WalletRepository) WalletService {
	return &walletService{repo}
}

func (s *walletService) GetAllWallets() ([]models.Wallet, error) {
	return s.repo.FindAll()
}

func (s *walletService) GetWallet(id string) (*models.Wallet, error) {
	return s.repo.FindByID(id)
}

func (s *walletService) CreateWallet(req WalletRequest) error {
	req.Name = strings.TrimSpace(req.Name)
	if req.Name == "" {
		return errors.New("nama dompet tidak boleh kosong")
	}

	wallet := &models.Wallet{
		ID:       uuid.NewString(),
		Name:     req.Name,
		Balance:  req.InitialBalance,
		Currency: "IDR",
	}

	return s.repo.Create(wallet)
}

func (s *walletService) UpdateWallet(id string, req WalletRequest) error {
	wallet, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("dompet tidak ditemukan")
	}

	req.Name = strings.TrimSpace(req.Name)
	if req.Name != "" {
		wallet.Name = req.Name
	}

	// Catatan: Mengubah saldo (balance) secara manual di sini bisa berbahaya untuk integritas data.
	// Sebaiknya perubahan saldo hanya melalui pencatatan transaksi. Untuk MVP v1.1, kita hanya izinkan ubah nama.

	return s.repo.Update(wallet)
}

func (s *walletService) DeleteWallet(id string) error {
	// Aturan Bisnis: Sistem harus menyisakan minimal 1 dompet agar aplikasi tetap berfungsi
	count, err := s.repo.Count()
	if err != nil {
		return errors.New("gagal memverifikasi jumlah dompet")
	}
	if count <= 1 {
		return errors.New("penghapusan ditolak: anda harus memiliki minimal satu dompet yang aktif")
	}

	// Opsional: Cek apakah ada transaksi yang terikat pada dompet ini sebelum dihapus

	return s.repo.Delete(id)
}
