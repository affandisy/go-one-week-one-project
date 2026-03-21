package main

import (
	"log"

	"github.com/affandisy/padel-booking-system/config"
	"github.com/affandisy/padel-booking-system/handlers"
	"github.com/affandisy/padel-booking-system/repositories"
	"github.com/affandisy/padel-booking-system/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.ConnectDatabase()
	db := config.DB

	// Inisialisasi Modul Autentikasi
	userRepo := repositories.NewUserRepository(db)
	// Gunakan kunci rahasia yang kuat untuk JWT (idealnya dari .env)
	authService := services.NewAuthService(userRepo, "PADEL_SECRET_KEY_2026")
	authHandler := handlers.NewAuthHandler(authService)

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api/v1")

	// Routing Autentikasi sesuai PRD
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/verify-otp", authHandler.VerifyOTP)
	auth.Post("/login", authHandler.Login)

	log.Println("Server Padel Booking berjalan di port 3000")
	log.Fatal(app.Listen(":3000"))
}
