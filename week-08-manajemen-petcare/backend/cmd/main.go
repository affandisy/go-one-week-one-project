package main

import (
	"log"

	pgrepo "github.com/affandisy/petcare-system/internal/adapter/driven/postgres"
	"github.com/affandisy/petcare-system/internal/adapter/driving/rest"
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
	billingRepo := pgrepo.NewBillingRepository(db)
	billingUseCase := usecase.NewBillingUseCase(billingRepo)
	billingHandler := rest.NewBillingHandler(billingUseCase)
	ownerRepo, petRepo := pgrepo.NewMasterRepository(db)

	// 3. Inisiasi UseCases
	ownerUseCase := usecase.NewOwnerUseCase(ownerRepo)
	petUseCase := usecase.NewPetUseCase(petRepo, ownerRepo) // PetUseCase butuh ownerRepo utk validasi

	// 4. Inisiasi Handlers
	masterHandler := rest.NewMasterHandler(ownerUseCase, petUseCase)

	// 5. Setup Router Fiber
	app := fiber.New()
	api := app.Group("/api/v1")

	api.Post("/invoices", billingHandler.CreateInvoice)
	api.Post("/owners", masterHandler.CreateOwner)
	api.Get("/owners", masterHandler.GetOwners)

	api.Post("/pets", masterHandler.CreatePet)
	api.Get("/pets", masterHandler.GetPets)

	log.Println("Server PetCare berjalan di port 3000...")
	log.Fatal(app.Listen(":3000"))
}
