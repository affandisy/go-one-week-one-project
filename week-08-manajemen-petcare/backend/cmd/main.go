package main

import (
	"fmt"
	"log"

	"github.com/affandisy/petcare-system/internal/adapter/driven/postgres"
	"github.com/affandisy/petcare-system/internal/adapter/driving/rest"
	"github.com/affandisy/petcare-system/internal/adapter/driving/rest/middleware"
	"github.com/affandisy/petcare-system/internal/config"
	"github.com/affandisy/petcare-system/internal/core/usecase"
	"github.com/gofiber/fiber/v2"
	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Muat Konfigurasi dari .env
	cfg := config.Load()

	// 2. Koneksi PostgreSQL menggunakan variabel lingkungan
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(gorm_postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}

	// 3. Inisiasi Repositori
	userRepo := postgres.NewUserRepository(db)
	billingRepo := postgres.NewBillingRepository(db)

	authUseCase := usecase.NewAuthUseCase(userRepo, cfg.JWTSecret)
	billingUseCase := usecase.NewBillingUseCase(billingRepo)

	authHandler := rest.NewAuthHandler(authUseCase)
	billingHandler := rest.NewBillingHandler(billingUseCase)

	paymentRepo := postgres.NewPaymentRepository(db)
	auditRepo := postgres.NewAuditRepository(db) // Repo baru

	// Inisiasi UseCase
	auditUseCase := usecase.NewAuditUseCase(auditRepo) // UseCase baru

	paymentUseCase := usecase.NewPaymentUseCase(paymentRepo, auditUseCase)
	paymentHandler := rest.NewPaymentHandler(paymentUseCase)

	// Suntikkan rahasia ke Middleware juga
	middleware.SetJWTSecret(cfg.JWTSecret)

	// 6. Router Fiber
	app := fiber.New()
	api := app.Group("/api/v1")

	// Rute Publik
	api.Post("/auth/login", authHandler.Login)

	// Rute Terproteksi
	api.Use(middleware.Protect())
	api.Post("/invoices", middleware.RequireRole("Cashier", "Manager"), billingHandler.CreateInvoice)
	api.Post("/payments", middleware.RequireRole("Cashier", "Manager"), paymentHandler.Process)

	// Jalankan Peladen menggunakan AppPort dari Config
	log.Printf("Server PetCare berjalan di port %s...", cfg.AppPort)
	log.Fatal(app.Listen(":" + cfg.AppPort))
}
