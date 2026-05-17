package usecase

import (
	"errors"
	"time"

	"github.com/affandisy/financial-app/internal/core/domain"
	"github.com/affandisy/financial-app/internal/core/port"
	"github.com/google/uuid"
)

type transactionUseCase struct {
	trxRepo      port.TransactionRepository
	walletRepo   port.WalletRepository
	categoryRepo port.CategoryRepository
}

// Injeksi 3 Repositori untuk melakukan orkestrasi
func NewTransactionUseCase(t port.TransactionRepository, w port.WalletRepository, c port.CategoryRepository) port.TransactionUseCase {
	return &transactionUseCase{t, w, c}
}

func (uc *transactionUseCase) RecordTransaction(walletID, txType, categoryID string, amount float64, note string) error {
	// 1. Validasi Aturan Bisnis Dasar
	if amount <= 0 {
		return errors.New("nominal harus lebih dari 0")
	}
	if txType != "income" && txType != "expense" {
		return errors.New("tipe transaksi tidak valid")
	}

	// 2. Ambil Entitas Dompet
	wallet, err := uc.walletRepo.FindByID(walletID)
	if err != nil {
		return errors.New("dompet tidak ditemukan")
	}

	// 3. Aturan Bisnis: Fallback kategori "Lainnya"
	if categoryID == "" {
		categories, _ := uc.categoryRepo.GetCategoriesByType(txType)
		for _, cat := range categories {
			if cat.Name == "Lainnya" {
				categoryID = cat.ID
				break
			}
		}
	}

	// 4. Buat Entitas Transaksi
	trx := &domain.Transaction{
		ID:         uuid.NewString(),
		WalletID:   wallet.ID,
		Type:       txType,
		CategoryID: categoryID,
		Amount:     amount,
		Note:       note,
		DateTime:   time.Now(),
	}

	// 5. Perilaku Domain: Kalkulasi Saldo Dompet
	if txType == "expense" {
		wallet.DeductBalance(amount) // Menggunakan method dari entitas Wallet
	} else {
		wallet.AddBalance(amount)
	}

	// 6. Simpan menggunakan Port (Adapter SQLite nanti akan menangani DB Transaction-nya)
	return uc.trxRepo.SaveWithWalletUpdate(trx, wallet)
}

func (uc *transactionUseCase) GetRecentHistory(walletID string) ([]domain.Transaction, error) {
	return uc.trxRepo.GetRecentByWallet(walletID, 10) // Ambil 10 terakhir
}
