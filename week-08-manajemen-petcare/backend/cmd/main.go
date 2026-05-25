package main

import (
	"log"

	"github.com/affandisy/petcare-system/internal/adapter/driven/postgres"
	"github.com/affandisy/petcare-system/internal/adapter/driving/rest"
	"github.com/affandisy/petcare-system/internal/adapter/driving/rest/middleware"
	"github.com/affandisy/petcare-system/internal/core/usecase"
	"github.com/gofiber/fiber/v2"
	pgdriver "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// 1. Koneksi PostgreSQL
	dsn := "host=localhost user=postgres password=secret dbname=petcare port=5432 sslmode=disable"
	db, err := gorm.Open(pgdriver.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi database")
	}

	// 2. Setup Adapters & UseCases
	billingRepo := postgres.NewBillingRepository(db)
	billingUseCase := usecase.NewBillingUseCase(billingRepo)
	billingHandler := rest.NewBillingHandler(billingUseCase)
	ownerRepo, petRepo := postgres.NewMasterRepository(db)

	// 3. Inisiasi UseCases
	ownerUseCase := usecase.NewOwnerUseCase(ownerRepo)
	petUseCase := usecase.NewPetUseCase(petRepo, ownerRepo) // PetUseCase butuh ownerRepo utk validasi

	// 4. Inisiasi Handlers
	masterHandler := rest.NewMasterHandler(ownerUseCase, petUseCase)

	// 1. Inisiasi Repositori Gizi
	nutritionRepo := postgres.NewNutritionRepository(db)

	// 2. Inisiasi UseCase Gizi
	nutritionUseCase := usecase.NewNutritionUseCase(nutritionRepo)

	// 3. Inisiasi Handler Gizi
	nutritionHandler := rest.NewNutritionHandler(nutritionUseCase)

	schedulingRepo := postgres.NewSchedulingRepository(db)
	schedulingUseCase := usecase.NewSchedulingUseCase(schedulingRepo)
	schedulingHandler := rest.NewSchedulingHandler(schedulingUseCase)

	userRepo := postgres.NewUserRepository(db)
	authUseCase := usecase.NewAuthUseCase(userRepo)
	authHandler := rest.NewAuthHandler(authUseCase)

	// 5. Setup Router Fiber
	app := fiber.New()

	auth := app.Group("/api/v1/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)

	api := app.Group("/api/v1")

	// Gunakan Middleware Protect untuk semua rute di bawah /api/v1/
	api.Use(middleware.Protect())

	// Contoh Proteksi Role: Hanya Kasir dan Manajer yang bisa membuat tagihan
	api.Post("/invoices", middleware.RequireRole("Cashier", "Manager"), billingHandler.CreateInvoice)

	// Contoh Proteksi Role: Hanya Resepsionis dan Manajer yang bisa mengelola Master Data & Jadwal
	api.Post("/owners", middleware.RequireRole("Receptionist", "Manager"), masterHandler.CreateOwner)
	api.Post("/appointments", middleware.RequireRole("Receptionist", "Manager"), schedulingHandler.CreateAppointment)

	// Contoh Proteksi Role: Groomer dan Manajer bisa mencatat Rekam Gizi
	api.Post("/nutrition", middleware.RequireRole("Groomer", "Manager"), nutritionHandler.CreateLog)

	log.Println("Server PetCare berjalan di port 3000...")
	log.Fatal(app.Listen(":3000"))
}
