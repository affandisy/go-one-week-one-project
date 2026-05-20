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

	// 3. Setup Router Fiber
	app := fiber.New()
	api := app.Group("/api/v1")

	api.Post("/invoices", billingHandler.CreateInvoice)

	log.Println("Server PetCare berjalan di port 3000...")
	log.Fatal(app.Listen(":3000"))
}
