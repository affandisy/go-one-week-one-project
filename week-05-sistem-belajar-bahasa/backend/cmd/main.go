package main

import (
	"log"

	"github.com/affandi/belajar-bahasa/config"
	"github.com/affandi/belajar-bahasa/handlers"
	"github.com/affandi/belajar-bahasa/repositories"
	"github.com/affandi/belajar-bahasa/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Inisialisasi Database
	config.ConnectDatabase()
	db := config.DB

	// Kunci rahasia untuk JWT (Idealnya ditaruh di .env)
	jwtSecret := "LANGUAGE_LEARNING_SECRET_KEY_2026"

	// 2. Inisialisasi Modul Auth (Dependency Injection)
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	// 3. Setup Fiber App
	app := fiber.New()

	// Middleware standar
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*", // Izinkan frontend SvelteKit nanti
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))

	// 4. Setup Routing sesuai PRD
	api := app.Group("/api/v1")

	// Endpoint Publik (Auth)
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/login", authHandler.Login)

	// 5. Jalankan Server
	log.Println("🚀 Server Language Learning berjalan di port 3000")
	log.Fatal(app.Listen(":3000"))
}
