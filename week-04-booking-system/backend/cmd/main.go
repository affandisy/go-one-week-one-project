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

	app.Static("/public", "./public")

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

	pricingRepo := repositories.NewPricingRepository(db)
	pricingService := services.NewPricingService(pricingRepo)
	pricingHandler := handlers.NewPricingHandler(pricingService)

	bookingRepo := repositories.NewBookingRepository(db)
	bookingService := services.NewBookingService(bookingRepo)
	bookingHandler := handlers.NewBookingHandler(bookingService)

	bookingService.RunAutoExpireJob()
	log.Println("⚙️ Background Job Auto-Expire Booking telah diaktifkan")

	// ================= PROTECTED ROUTES =================
	protected := api.Group("/", middlewares.JWTProtected(jwtSecret))

	courts := protected.Group("/courts")
	courts.Get("/", courtHandler.GetAll)
	courts.Get("/:id", courtHandler.GetByID)

	api.Get("/availability", bookingHandler.GetAvailability)

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

	adminPricing := adminRoutes.Group("/pricing-rules")
	adminPricing.Get("/", pricingHandler.GetByCourt)
	adminPricing.Post("/", pricingHandler.Create)
	adminPricing.Put("/:id", pricingHandler.Update)
	adminPricing.Delete("/:id", pricingHandler.Delete)

	adminBookings := adminRoutes.Group("/bookings")
	adminBookings.Get("/", bookingHandler.GetAllBookings) // Lihat semua booking
	adminBookings.Put("/:id/cancel", bookingHandler.CancelBooking)

	customerBookings := protected.Group("/bookings")
	customerBookings.Post("/", bookingHandler.CreateBooking)          // Buat Booking
	customerBookings.Get("/me", bookingHandler.GetMyBookings)         // Riwayat Pribadi
	customerBookings.Put("/:id/cancel", bookingHandler.CancelBooking) // Batal Booking

	customerBookings.Post("/:id/pay", bookingHandler.PayBooking)
	customerBookings.Get("/:id/receipt", bookingHandler.DownloadReceipt)

	log.Println("Server Padel Booking berjalan di port 3000")
	log.Fatal(app.Listen(":3000"))
}
