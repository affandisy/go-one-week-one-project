package config

import (
	"log"

	"github.com/affandisy/padel-booking-system/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	// Ganti kredensial sesuai dengan database lokal Anda
	dsn := "host=localhost user=postgres password=postgres dbname=padel_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal(" Gagal terhubung ke PostgreSQL: ", err)
	}

	// Wajib untuk mengaktifkan uuid_generate_v4() sesuai PRD[cite: 2]
	database.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	log.Println("✅ Berhasil terhubung ke PostgreSQL")
	DB = database

	// Migrasi Terpusat
	err = DB.AutoMigrate(
		&models.User{},        // Independen
		&models.Court{},       // Independen
		&models.PricingRule{}, // Bergantung pada Court
	)

	if err != nil {
		log.Fatal("Gagal melakukan migrasi: ", err)
	}

	log.Println("Migrasi Database Phase 1 Selesai!")
}
