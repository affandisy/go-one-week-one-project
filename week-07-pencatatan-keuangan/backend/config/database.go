// config/database.go
package config

import (
	"log"

	"github.com/affandisy/financial-app/models"
	"github.com/glebarez/sqlite"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("data/keuangan_app.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database SQLite:", err)
	}

	// Migrasi Skema Database[cite: 1]
	err = database.AutoMigrate(&models.Wallet{}, &models.Category{}, &models.Transaction{})
	if err != nil {
		log.Fatal("Gagal migrasi database:", err)
	}

	DB = database
	seedInitialData(database)
	log.Println("✅ Database SQLite siap!")
}

// Menjalankan FR-006 (Default Categories) dan FR-007 (Single Wallet)[cite: 1]
func seedInitialData(db *gorm.DB) {
	var walletCount int64
	db.Model(&models.Wallet{}).Count(&walletCount)

	if walletCount == 0 {
		// Buat Dompet Default[cite: 1]
		db.Create(&models.Wallet{
			ID:       uuid.NewString(),
			Name:     "Dompet Saya",
			Balance:  0,
			Currency: "IDR",
		})

		// Buat Kategori Pengeluaran[cite: 1]
		expenses := []models.Category{
			{ID: uuid.NewString(), Name: "Makan", Icon: "🍚", Type: "expense", IsActive: true},
			{ID: uuid.NewString(), Name: "Belanja", Icon: "🛒", Type: "expense", IsActive: true},
			{ID: uuid.NewString(), Name: "Transport", Icon: "🚗", Type: "expense", IsActive: true},
			{ID: uuid.NewString(), Name: "Tagihan", Icon: "⚡", Type: "expense", IsActive: true},
			{ID: uuid.NewString(), Name: "Lainnya", Icon: "📝", Type: "expense", IsActive: true},
		}
		for _, cat := range expenses {
			db.Create(&cat)
		}

		// Buat Kategori Pemasukan[cite: 1]
		incomes := []models.Category{
			{ID: uuid.NewString(), Name: "Gaji", Icon: "💰", Type: "income", IsActive: true},
			{ID: uuid.NewString(), Name: "Bonus", Icon: "🎁", Type: "income", IsActive: true},
			{ID: uuid.NewString(), Name: "Jualan", Icon: "🏪", Type: "income", IsActive: true},
			{ID: uuid.NewString(), Name: "Lainnya", Icon: "📝", Type: "income", IsActive: true},
		}
		for _, cat := range incomes {
			db.Create(&cat)
		}
	}
}
