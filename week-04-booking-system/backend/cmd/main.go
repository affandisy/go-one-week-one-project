package main

import (
	"log"

	"github.com/affandisy/padel-booking-system/config"
	"github.com/affandisy/padel-booking-system/handlers"
	"github.com/affandisy/padel-booking-system/middlewares"
	"github.com/affandisy/padel-booking-system/repositories"
	"github.com/affandisy/padel-booking-system/services"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	config.ConnectDatabase()
	db := config.DB

	jwtSecret := "testing_key"

	// Inisialisasi Modul Autentikasi
	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, "testing_key")
	authHandler := handlers.NewAuthHandler(authService)

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api/v1")

	// ================= PUBLIC ROUTES =================
	auth := api.Group("/auth")
	auth.Post("/register", authHandler.Register)
	auth.Post("/verify-otp", authHandler.VerifyOTP)
	auth.Post("/login", authHandler.Login)

	// 1. Inisialisasi Modul Court (Lapangan)
	courtRepo := repositories.NewCourtRepository(db)
	courtService := services.NewCourtService(courtRepo)
	courtHandler := handlers.NewCourtHandler(courtService)

	// ================= PROTECTED ROUTES =================
	protected := api.Group("/", middlewares.JWTProtected(jwtSecret))

	courts := protected.Group("/courts")
	courts.Get("/", courtHandler.GetAll)
	courts.Get("/:id", courtHandler.GetByID)

	protected.Get("/me", func(c *fiber.Ctx) error {
		userID := c.Locals("user_id").(string)
		role := c.Locals("role").(string)
		return c.JSON(fiber.Map{
			"message": "Ini adalah data profil Anda",
			"user_id": userID,
			"role":    role,
		})
	})

	adminRoutes := protected.Group("/admin", middlewares.RequireRoles("admin", "owner"))

	adminRoutes.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Selamat datang di ruang rahasia Admin!"})
	})

	adminCourts := adminRoutes.Group("/courts")
	adminCourts.Post("/", courtHandler.Create)
	adminCourts.Put("/:id", courtHandler.Update)
	adminCourts.Delete("/:id", courtHandler.Delete)

	log.Println("Server Padel Booking berjalan di port 3000")
	log.Fatal(app.Listen(":3000"))
}
