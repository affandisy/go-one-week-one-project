package main

import (
	"log"
	"os"

	"github.com/affandisy/hardware-erp/internal/adapter/driven/postgres"
	"github.com/affandisy/hardware-erp/internal/core/engine"
	"github.com/affandisy/hardware-erp/internal/core/service"
	"github.com/affandisy/hardware-erp/internal/core/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"

	gorm_postgres "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DTO untuk Request HTTP
type CreateComponentReq struct {
	SKU          string                 `json:"sku"`
	Name         string                 `json:"name"`
	Category     string                 `json:"category"`
	Manufacturer string                 `json:"manufacturer"`
	Price        float64                `json:"price"`
	Specs        map[string]interface{} `json:"specs"`
}

type ValidateAssemblyReq struct {
	ComponentSKUs []string `json:"component_skus"`
}

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		// Fallback DSN jika environment variable tidak diset
		dsn = "host=localhost user=postgres password=rahasia_super_kuat_123 dbname=petcare port=5432 sslmode=disable"
	}

	db, err := gorm.Open(gorm_postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Kritikal: Gagal terhubung ke PostgreSQL: %v", err)
	}
	log.Println("Database PostgreSQL berhasil terhubung.")

	// =========================================================================
	// 2. INISIALISASI DOMAIN SERVICES & DRIVEN ADAPTERS
	// =========================================================================
	// Muat Schema Registry untuk validasi JSONB di hulu
	schemaValidator, err := service.NewSchemaValidator()
	if err != nil {
		log.Fatalf("Kritikal: Gagal memuat JSON Schema Registry: %v", err)
	}
	log.Println("JSON Schema Registry untuk CPU, GPU, dan Mobile sukses dimuat ke memori.")

	// Inisialisasi Driven Repository GORM
	compRepo := postgres.NewComponentRepository(db)

	socketRule := engine.NewSocketRule()
	powerRule := engine.NewPowerRule()

	// Suntikkan semua aturan spesifik ke dalam Core Assembly Engine
	hardwareRuleEngine := engine.NewAssemblyEngine(socketRule, powerRule)
	log.Println("Hardware Rule Engine aktif dengan SocketRule dan PowerBudgetRule.")

	compUseCase := usecase.NewComponentUseCase(compRepo, schemaValidator)
	assemblyUseCase := usecase.NewAssemblyUseCase(compRepo, hardwareRuleEngine)

	app := fiber.New(fiber.Config{
		AppName: "Hardware Assembly & Build ERP v1.2",
	})

	// Middleware Logger untuk observabilitas REST API
	app.Use(logger.New())

	api := app.Group("/api")

	// --- ENDPOINT FR-001: Katalog Komponen Dinamis ---
	api.Post("/components", func(c *fiber.Ctx) error {
		var req CreateComponentReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON Request tidak valid"})
		}

		comp, err := compUseCase.RegisterComponent(req.SKU, req.Name, req.Category, req.Manufacturer, req.Price, req.Specs)
		if err != nil {
			// Jika error berasal dari JSON Schema Validation, status otomatis Bad Request
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusCreated).JSON(fiber.Map{
			"message": "Komponen baru berhasil didaftarkan ke katalog",
			"data":    comp,
		})
	})

	// --- ENDPOINT FR-002 & FR-003: Rule Engine & Assembly Simulator UI ---
	api.Post("/assemblies/validate", func(c *fiber.Ctx) error {
		var req ValidateAssemblyReq
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON Request tidak valid"})
		}

		buildReport, err := assemblyUseCase.SimulateBuild(req.ComponentSKUs)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
		}

		// Jika rakitan tidak kompatibel, kita kembalikan status 200 OK namun is_compatible = false
		// Sesuai dengan spesifikasi AC-FR-003 agar UI dapat merespons visual peringatan
		return c.JSON(fiber.Map{
			"message": "Simulasi perakitan selesai dievaluasi",
			"data":    buildReport,
		})
	})

	api.Post("/assemblies/checkout", func(c *fiber.Ctx) error {
		var req ValidateAssemblyReq // Menggunakan format request yang sama
		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
		}

		completedBuild, err := assemblyUseCase.CheckoutAssembly(req.ComponentSKUs)
		if err != nil {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Checkout berhasil. Stok komponen telah diperbarui.",
			"data":    completedBuild,
		})
	})

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Peladen Hardware ERP berjalan aktif pada port %s...", port)
	log.Fatal(app.Listen(":" + port))
}
