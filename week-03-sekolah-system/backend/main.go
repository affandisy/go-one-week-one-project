package main

import (
	"log"

	"github.com/affandisy/school-system-sma/handlers"
	"github.com/affandisy/school-system-sma/middleware"
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

	authHandler := handlers.NewAuthHandler(authService)

	// Inisialisasi Academic Module
	academicRepo := repository.NewAcademicRepository(db)
	academicService := service.NewAcademicService(academicRepo)
	academicHandler := handlers.NewAcademicHandler(academicService)

	// Public Routes
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)

	jwtMid := middleware.JWTProtected("SUPER_SECRET_KEY")

	// Routing (idealnya dilindungi JWT Middleware, tapi untuk tes kita buka dulu)
	master := api.Group("/academic-years", jwtMid)
	// Hanya Admin/TU yang boleh membuat Tahun Ajaran
	master.Post("/", middleware.RequireRoles("ADMIN", "TU"), academicHandler.CreateAcademicYear)
	master.Get("/", academicHandler.GetAcademicYears)

	log.Println("Server berjalan di port 3000")
	app.Listen(":3000")
}
