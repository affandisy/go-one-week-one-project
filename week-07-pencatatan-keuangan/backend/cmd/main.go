// cmd/main.go
package main

import (
	"log"

	"github.com/affandisy/financial-app/config"
	"github.com/affandisy/financial-app/handlers"
	"github.com/affandisy/financial-app/repositories"
	"github.com/affandisy/financial-app/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// Inisialisasi Database (SQLite Local - FR-005)[cite: 1]
	config.ConnectDatabase()

	// Inisiasi Modul Master
	masterRepo := repositories.NewMasterRepository(config.DB)
	masterService := services.NewMasterService(masterRepo)
	masterHandler := handlers.NewMasterHandler(masterService)

	// Inisiasi Modul Transaksi (BARU)
	transactionRepo := repositories.NewTransactionRepository(config.DB)
	transactionService := services.NewTransactionService(transactionRepo, masterRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	app := fiber.New()

	// CORS agar bisa diakses oleh React Frontend nantinya
	app.Use(cors.New())

	api := app.Group("/api/v1")

	// Rute Modul Master
	api.Get("/wallet", masterHandler.GetWallet)
	api.Get("/categories", masterHandler.GetCategories)

	// Rute Modul Transaksi (BARU)
	api.Post("/transactions", transactionHandler.Create)
	api.Get("/transactions/recent", transactionHandler.GetHistory)

	log.Println("Peladen berjalan di port 3000...")
	log.Fatal(app.Listen(":3000"))
}
