package config

import (
	"fmt"
	"log"
	"os"

	"github.com/affandi/belajar-bahasa/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Jakarta",
		host, user, password, dbName, port)

	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	database.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`)

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
	log.Println("✅ Database berhasil terhubung menggunakan Environment Variables!")
}
