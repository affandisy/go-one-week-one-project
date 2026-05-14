package usecase

import (
	"errors"
	"strings"

	"github.com/affandisy/financial-app/internal/core/domain"
	"github.com/affandisy/financial-app/internal/core/port"
	"github.com/google/uuid"
)

type walletUseCase struct {
	repo port.WalletRepository
}

func NewWalletUseCase(repo port.WalletRepository) port.WalletUseCase {
	return &walletUseCase{repo}
}

func (uc *walletUseCase) CreateWallet(name string, initialBalance float64) error {
	cleanName := strings.TrimSpace(name)
	if cleanName == "" {
		return errors.New("nama dompet tidak boleh kosong")
	}

	wallet := &domain.Wallet{
		ID:       uuid.NewString(),
		Name:     cleanName,
		Balance:  initialBalance,
		Currency: "IDR",
	}

	return uc.repo.Save(wallet)
}

func (uc *walletUseCase) GetWallet(id string) (*domain.Wallet, error) {
	return uc.repo.FindByID(id)
}

func (uc *walletUseCase) ListWallets() ([]domain.Wallet, error) {
	return uc.repo.FindAll()
}
