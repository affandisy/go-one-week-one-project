package services

import (
	"errors"
	"time"

	"github.com/affandisy/financial-app/models"
	"github.com/affandisy/financial-app/repositories"
	"github.com/google/uuid"
)

type TransactionRequest struct {
	Type       string  `json:"type"` // "income" atau "expense"[cite: 1]
	CategoryID string  `json:"category_id"`
	Amount     float64 `json:"amount"`
	Note       string  `json:"note"`
}

type TransactionService interface {
	RecordTransaction(req TransactionRequest) error
	GetDashboardHistory() ([]models.Transaction, error)
	DeleteTransaction(id string) error
	UpdateTransaction(id string, req TransactionRequest) error
}

type transactionService struct {
	trxRepo    repositories.TransactionRepository
	masterRepo repositories.MasterRepository
}

func NewTransactionService(trxRepo repositories.TransactionRepository, masterRepo repositories.MasterRepository) TransactionService {
	return &transactionService{trxRepo, masterRepo}
}

func (s *transactionService) RecordTransaction(req TransactionRequest) error {
	// Aturan Bisnis 1: Nominal > 0[cite: 1]
	if req.Amount <= 0 {
		return errors.New("nominal harus lebih dari 0")
	}
	if req.Type != "income" && req.Type != "expense" {
		return errors.New("tipe transaksi tidak valid")
	}

	wallet, err := s.masterRepo.GetDefaultWallet()
	if err != nil {
		return errors.New("dompet utama tidak ditemukan")
	}

	// Aturan Bisnis 2: Auto-default kategori "Lainnya" jika kosong[cite: 1]
	if req.CategoryID == "" {
		categories, _ := s.masterRepo.GetCategoriesByType(req.Type)
		for _, cat := range categories {
			if cat.Name == "Lainnya" {
				req.CategoryID = cat.ID
				break
			}
		}
	}

	trx := &models.Transaction{
		ID:         uuid.NewString(),
		WalletID:   wallet.ID,
		Type:       req.Type,
		CategoryID: req.CategoryID,
		Amount:     req.Amount,
		Note:       req.Note,
		DateTime:   time.Now(),
	}

	// Kalkulasi saldo dompet
	if req.Type == "expense" {
		wallet.Balance -= req.Amount
	} else {
		wallet.Balance += req.Amount
	}

	return s.trxRepo.CreateWithWalletUpdate(trx, wallet)
}

func (s *transactionService) GetDashboardHistory() ([]models.Transaction, error) {
	return s.trxRepo.GetRecentTransactions(10) // Menampilkan 10 transaksi terakhir
}

func (s *transactionService) DeleteTransaction(id string) error {
	trx, err := s.trxRepo.FindByID(id)
	if err != nil {
		return errors.New("transaksi tidak ditemukan")
	}

	wallet, _ := s.masterRepo.GetDefaultWallet()

	// Kembalikan saldo seperti sebelum transaksi terjadi
	if trx.Type == "expense" {
		wallet.Balance += trx.Amount
	} else {
		wallet.Balance -= trx.Amount
	}

	return s.trxRepo.DeleteWithWalletUpdate(trx, wallet)
}

func (s *transactionService) UpdateTransaction(id string, req TransactionRequest) error {
	trx, err := s.trxRepo.FindByID(id)
	if err != nil {
		return errors.New("transaksi tidak ditemukan")
	}
	wallet, _ := s.masterRepo.GetDefaultWallet()

	// 1. Batalkan efek transaksi lama terhadap saldo
	if trx.Type == "expense" {
		wallet.Balance += trx.Amount
	} else {
		wallet.Balance -= trx.Amount
	}

	// 2. Terapkan efek transaksi baru
	trx.CategoryID = req.CategoryID
	trx.Amount = req.Amount
	trx.Note = req.Note
	// (Asumsi: Tipe 'income'/'expense' tidak boleh diubah untuk mencegah bug kompleks. Jika salah pilih, user disarankan hapus & buat baru)

	if trx.Type == "expense" {
		wallet.Balance -= trx.Amount
	} else {
		wallet.Balance += trx.Amount
	}

	return s.trxRepo.UpdateWithWalletUpdate(trx, trx, wallet)
}
