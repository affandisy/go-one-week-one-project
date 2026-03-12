package main

import (
	"log"

	"github.com/affandisy/school-system-sma/handlers"
	"github.com/affandisy/school-system-sma/models"
	"github.com/affandisy/school-system-sma/repository"
	"github.com/affandisy/school-system-sma/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Koneksi DB
	dsn := "host=localhost user=school_user password=school_password dbname=school_db port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal koneksi ke database")
	}

	// Auto Migrate
	db.AutoMigrate(
		&models.User{},
		&models.AcademicYear{},
		&models.Class{},
		&models.StudyProgram{},
		&models.Subject{},
	)

	// Dependency Injection
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "SUPER_SECRET_KEY")

	app := fiber.New()
	app.Use(cors.New()) // Penting agar Frontend Svelte bisa komunikasi

	api := app.Group("/api/v1")
	auth := api.Group("/auth")

	// Endpoint Login
	auth.Post("/login", func(c *fiber.Ctx) error {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tidak valid"})
		}
		token, err := authService.Login(req.Email, req.Password)
		if err != nil {
			return c.Status(401).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"data": fiber.Map{"token": token}})
	})

	// Endpoint Register Dummy (Untuk testing awal)
	auth.Post("/register", func(c *fiber.Ctx) error {
		var req service.RegisterRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(400).JSON(fiber.Map{"error": "Format tidak valid"})
		}
		if err := authService.Register(req); err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Gagal register"})
		}
		return c.JSON(fiber.Map{"message": "Berhasil daftar!"})
	})

	// Inisialisasi Academic Module
	academicRepo := repository.NewAcademicRepository(db)
	academicService := service.NewAcademicService(academicRepo)
	academicHandler := handlers.NewAcademicHandler(academicService)

	// Routing (idealnya dilindungi JWT Middleware, tapi untuk tes kita buka dulu)
	masterData := api.Group("/academic-years")
	masterData.Post("/", academicHandler.CreateAcademicYear)
	masterData.Get("/", academicHandler.GetAcademicYears)

	log.Println("Server berjalan di port 3000")
	app.Listen(":3000")
}
