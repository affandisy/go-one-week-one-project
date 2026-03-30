package config

import (
	"log"

	"github.com/affandi/belajar-bahasa/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=localhost user=postgres password=test dbname=lang_learning_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	// WAJIB: Aktifkan ekstensi UUID di PostgreSQL
	database.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

	// Migrasi Skema
	err = database.AutoMigrate(
		&models.User{},
		&models.Module{},
		&models.Material{},
		&models.UserProgress{},
	)
	if err != nil {
		log.Fatal("Gagal melakukan migrasi database:", err)
	}

	DB = database
	log.Println("✅ Database berhasil terhubung dan dimigrasi!")
}
