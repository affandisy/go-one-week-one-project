package config

import (
	"fmt"
	"log"
	"os"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Variabel global agar database bisa dipanggil dari layer Repository nanti
var DB *gorm.DB

func ConnectDatabase() {
	// 1. Baca file .env
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan, menggunakan variabel sistem default")
	}

	// 2. Buat string koneksi (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		os.Getenv("DB_HOST"), os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME"), os.Getenv("DB_PORT"))

	// 3. Buka koneksi dengan GORM
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database: ", err)
	}

	log.Println("Koneksi ke PostgreSQL berhasil")

	// 4. Jalankan AutoMigrate
	// GORM akan otomatis membuat tabel jika belum ada, atau menambah kolom jika ada perubahan struct
	err = database.AutoMigrate(
		&models.User{},
		&models.Product{},
		&models.Transaction{},
		&models.TransactionItem{},
		&models.StockAdjustment{},
		&models.Partner{},
		&models.Location{},
		&models.ProductBatch{},
		&models.Warehouse{},
		&models.WarehouseStock{},
	)

	if err != nil {
		log.Fatal("Gagal melakukan migrasi database: ", err)
	}

	log.Println("Migrasi tabel database selesai")

	// Simpan instance database ke variabel global
	DB = database
}
