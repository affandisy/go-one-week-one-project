package main

import (
	"log"
	"os"

	"github.com/affandisy/pos-system/config"
	"github.com/affandisy/pos-system/handlers"
	"github.com/affandisy/pos-system/middlewares"
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

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	stockRepo := repositories.NewStockRepository(db)
	stockService := services.NewStockService(stockRepo)
	stockHandler := handlers.NewStockHandler(stockService)

	// --- INISIASI MODUL TRANSAKSI ---
	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo, productRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	// --- INISIASI MODUL LAPORAN ---
	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	// Inisiasi Modul User (Baru)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

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

	// Middleware Proteksi
	jwtAuth := middlewares.Protected(jwtSecret)
	onlyOwnerAdmin := middlewares.RequireRole("owner", "admin")
	// Kasir diizinkan mengakses produk dan transaksi
	anyRole := middlewares.RequireRole("owner", "admin", "cashier")

	api := app.Group("/api/v1")

	// Rute Publik
	api.Post("/auth/login", authHandler.Login)

	// --- RUTE TERPROTEKSI ---
	protected := api.Group("/", jwtAuth)

	// Produk: Kasir cuma bisa lihat, Pemilik/Admin bisa tambah/ubah/hapus
	products := protected.Group("/products")
	products.Get("/", anyRole, productHandler.GetAll)
	products.Get("/:id", anyRole, productHandler.GetByID)
	products.Post("/", onlyOwnerAdmin, productHandler.Create)
	products.Put("/:id", onlyOwnerAdmin, productHandler.Update)
	products.Delete("/:id", onlyOwnerAdmin, productHandler.Delete)

	// Manajemen Stok: Hanya Pemilik & Admin (PRD 5.3)
	stocks := protected.Group("/stocks", onlyOwnerAdmin)
	stocks.Post("/in", stockHandler.StockIn)
	stocks.Get("/history/:productId", stockHandler.GetHistory)

	// Kasir / Checkout: Bisa diakses semua (Kasir, Pemilik)
	protected.Post("/checkout", anyRole, transactionHandler.Checkout)

	// Laporan: HANYA PEMILIK
	reports := protected.Group("/reports", middlewares.RequireRole("owner"))
	reports.Get("/sales/daily", reportHandler.GetDailySales)
	reports.Get("/stocks/low", reportHandler.GetLowStocks)

	// Rute Manajemen Pengguna (HANYA OWNER & ADMIN)
	users := protected.Group("/users", onlyOwnerAdmin)
	users.Get("/", userHandler.GetAll)
	users.Post("/", userHandler.Create)
	users.Put("/:id", userHandler.Update)
	users.Delete("/:id", userHandler.Delete)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("🚀 Server POS berjalan di port %s", port)
	log.Fatal(app.Listen(":" + port))
}
