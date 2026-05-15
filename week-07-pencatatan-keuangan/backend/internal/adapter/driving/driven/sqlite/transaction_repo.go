package sqlite

import (
	"time"

	"github.com/affandisy/financial-app/internal/core/domain"
	"github.com/affandisy/financial-app/internal/core/port"
	"gorm.io/gorm"
)

// TransactionModel adalah DTO khusus database
type TransactionModel struct {
	ID string `gorm:"type:varchar(36);primaryKey"`

	// Composite Index: Sangat mempercepat pencarian rentang waktu untuk satu dompet spesifik
	WalletID string    `gorm:"type:varchar(36);index:idx_wallet_date"`
	DateTime time.Time `gorm:"column:date_time;index:idx_wallet_date;index:idx_date_only"`

	Type       string `gorm:"type:varchar(10);index"` // Indeks untuk filter pemasukan/pengeluaran
	CategoryID string `gorm:"type:varchar(36)"`
	Amount     float64
	Note       string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (TransactionModel) TableName() string { return "transactions" }

type transactionRepository struct {
	db *gorm.DB
}

// Pastikan mengimplementasikan port.TransactionRepository (Asumsi port ini sudah Anda buat)
func NewTransactionRepository(db *gorm.DB) port.TransactionRepository {
	// GORM otomatis akan membaca tag 'index' dan membuat B-Tree index di SQLite
	db.AutoMigrate(&TransactionModel{})
	return &transactionRepository{db}
}

// Contoh implementasi pencarian laporan bulanan yang super cepat berkat index
func (r *transactionRepository) FindByWalletAndDateRange(walletID string, start, end time.Time) ([]domain.Transaction, error) {
	var models []TransactionModel

	// Query ini akan langsung menembak 'idx_wallet_date' index (O(log N) complexity)
	// Tanpa full table scan, meskipun ada 100.000 transaksi!
	err := r.db.Where("wallet_id = ? AND date_time >= ? AND date_time < ?", walletID, start, end).
		Order("date_time DESC").
		Find(&models).Error

	if err != nil {
		return nil, err
	}

	// Mapper dari DB Model ke Domain Entity
	var transactions []domain.Transaction
	for _, m := range models {
		transactions = append(transactions, domain.Transaction{
			ID:         m.ID,
			WalletID:   m.WalletID,
			Type:       m.Type,
			CategoryID: m.CategoryID,
			Amount:     m.Amount,
			Note:       m.Note,
			DateTime:   m.DateTime,
		})
	}

	return transactions, nil
}

func (r *transactionRepository) Save(t *domain.Transaction) error {
	model := TransactionModel{
		ID:         t.ID,
		WalletID:   t.WalletID,
		Type:       t.Type,
		CategoryID: t.CategoryID,
		Amount:     t.Amount,
		Note:       t.Note,
		DateTime:   t.DateTime,
	}
	return r.db.Save(&model).Error
}
