package port

import "github.com/affandisy/financial-app/internal/core/domain"

// Driven Port (Outbound)
type TransactionRepository interface {
	// Menyimpan transaksi dan mengupdate dompet secara atomik (Database Transaction)
	SaveWithWalletUpdate(trx *domain.Transaction, wallet *domain.Wallet) error
	GetRecentByWallet(walletID string, limit int) ([]domain.Transaction, error)
}

// Untuk Kategori (Karena PRD mengharuskan fallback ke "Lainnya")
type CategoryRepository interface {
	GetCategoriesByType(txType string) ([]domain.Category, error)
}

// Driving Port (Inbound)
type TransactionUseCase interface {
	RecordTransaction(walletID, txType, categoryID string, amount float64, note string) error
	GetRecentHistory(walletID string) ([]domain.Transaction, error)
}
