package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config menampung semua variabel lingkungan
type Config struct {
	AppPort   string
	DBHost    string
	DBUser    string
	DBPass    string
	DBName    string
	DBPort    string
	JWTSecret string
}

// Load membaca file .env dan mengisinya ke struct Config
func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: File .env tidak ditemukan. Menggunakan variabel lingkungan sistem.")
	}

	return &Config{
		AppPort:   getEnv("APP_PORT", "3000"),
		DBHost:    getEnv("DB_HOST", "localhost"),
		DBUser:    getEnv("DB_USER", "postgres"),
		DBPass:    getEnv("DB_PASSWORD", "secret"),
		DBName:    getEnv("DB_NAME", "petcare"),
		DBPort:    getEnv("DB_PORT", "5432"),
		JWTSecret: getEnv("JWT_SECRET", "default_secret_key_ubah_segera"),
	}
}

// getEnv membaca variabel lingkungan atau menggunakan nilai fallback jika kosong
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
