package port

import "github.com/affandisy/financial-app/internal/core/domain"

// Driven Port (Outbound) -> Akan diimplementasikan oleh adapter/sqlite
type WalletRepository interface {
	Save(wallet *domain.Wallet) error
	FindByID(id string) (*domain.Wallet, error)
	FindAll() ([]domain.Wallet, error)
}

// Driving Port (Inbound) -> Akan diimplementasikan oleh usecase
type WalletUseCase interface {
	CreateWallet(name string, initialBalance float64) error
	GetWallet(id string) (*domain.Wallet, error)
	ListWallets() ([]domain.Wallet, error)
}
