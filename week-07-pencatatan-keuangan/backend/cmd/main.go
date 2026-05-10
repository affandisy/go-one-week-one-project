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

	// 2. Modul Dompet (BARU)
	walletRepo := repositories.NewWalletRepository(config.DB)
	walletService := services.NewWalletService(walletRepo)
	walletHandler := handlers.NewWalletHandler(walletService)

	// Inisiasi Modul Transaksi (BARU)
	transactionRepo := repositories.NewTransactionRepository(config.DB)
	transactionService := services.NewTransactionService(transactionRepo, walletRepo, masterRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	reportRepo := repositories.NewReportRepository(config.DB)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	app := fiber.New()

	// CORS agar bisa diakses oleh React Frontend nantinya
	app.Use(cors.New())

	api := app.Group("/api/v1")

	// Rute Modul Master
	api.Get("/wallet", masterHandler.GetWallet)
	api.Get("/categories", masterHandler.GetCategories)

	// Rute Transaksi
	api.Post("/transactions", transactionHandler.Create)
	api.Get("/transactions/recent", transactionHandler.GetHistory)
	api.Put("/transactions/:id", transactionHandler.Update)    // FR-003
	api.Delete("/transactions/:id", transactionHandler.Delete) // FR-003

	wallets := api.Group("/wallets")
	wallets.Get("/", walletHandler.GetAll)
	wallets.Get("/:id", walletHandler.GetByID)
	wallets.Post("/", walletHandler.Create)
	wallets.Put("/:id", walletHandler.Update)
	wallets.Delete("/:id", walletHandler.Delete)

	// Rute Laporan
	api.Get("/reports/monthly", reportHandler.GetMonthlySummary) // FR-004

	log.Println("Server berjalan di port 3000...")
	log.Fatal(app.Listen(":3000"))
}
