package main

import (
	"log"
	"os"

	"github.com/affandisy/pos-system/config"
	"github.com/affandisy/pos-system/handlers"
	"github.com/affandisy/pos-system/models"
	"github.com/affandisy/pos-system/repositories"
	"github.com/affandisy/pos-system/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Catatan: .env tidak ditemukan, menggunakan variabel OS.")
	}

	config.ConnectDatabase()
	db := config.DB
	jwtSecret := os.Getenv("JWT_SECRET")

	// --- INISIASI DEPENDENSI ---
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	// --- AUTO SEEDER UNTUK AKUN PEMILIK ---
	count, _ := userRepo.Count()
	if count == 0 {
		log.Println("Tabel pengguna kosong. Membuat akun Pemilik (Owner) bawaan...")
		hashedPIN, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)

		owner := &models.User{
			Username: "admin",
			PINHash:  string(hashedPIN),
			Role:     "owner",
		}
		userRepo.Create(owner)
		log.Println("✅ Akun Pemilik berhasil dibuat! Username: admin | PIN: 123456")
	}

	// --- SETUP FIBER ---
	app := fiber.New()
	app.Use(logger.New())
	app.Use(cors.New())

	api := app.Group("/api/v1")

	// Rute Publik
	api.Post("/auth/login", authHandler.Login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server POS berjalan di port %s", port)
	log.Fatal(app.Listen(":" + port))
}
