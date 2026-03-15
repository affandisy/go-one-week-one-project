package main

import (
	"log"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/middlewares"
	"github.com/affandisy/school-system-sma/config"
	"github.com/affandisy/school-system-sma/handlers"
	"github.com/affandisy/school-system-sma/middleware"
	"github.com/affandisy/school-system-sma/repository"
	"github.com/affandisy/school-system-sma/service"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func main() {
	// 1. Inisialisasi Database yang sudah dirapikan
	config.ConnectDatabase()
	db := config.DB

	// 2. Dependency Injection
	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, "SUPER_SECRET_KEY")
	authHandler := handlers.NewAuthHandler(authService)

	academicRepo := repository.NewAcademicRepository(db)
	academicService := service.NewAcademicService(academicRepo)
	academicHandler := handlers.NewAcademicHandler(academicService)

	payrollRepo := repository.NewPayrollRepository(db)
	payrollService := service.NewPayrollService(payrollRepo)
	payrollHandler := handlers.NewPayrollHandler(payrollService)

	recomRepo := repository.NewRecommendationRepository(db)
	recomService := service.NewRecommendationService(recomRepo)
	recomHandler := handlers.NewRecommendationHandler(recomService)

	annRepo := repository.NewAnnouncementRepository(db)
	annService := service.NewAnnouncementService(annRepo)
	annHandler := handlers.NewAnnouncementHandler(annService)

	// ... (di bagian rute)

	// 3. Setup Fiber & Middleware
	app := fiber.New()
	app.Use(cors.New())

	api := app.Group("/api/v1")

	// ================= ROUTES =================

	// Public
	auth := api.Group("/auth")
	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register) // Pastikan Anda sudah memindahkan ini ke handler

	// Protected (Membutuhkan JWT)
	protected := api.Group("/", middleware.JWTProtected("SUPER_SECRET_KEY"))

	// Master Data Akademik
	master := protected.Group("/academic-years")
	master.Post("/", middleware.RequireRoles("ADMIN", "TU"), academicHandler.CreateAcademicYear)
	master.Get("/", academicHandler.GetAcademicYears)

	// Payroll / Penggajian
	payroll := protected.Group("/salary-payslips")
	payroll.Post("/calculate", middleware.RequireRoles("KEUANGAN", "ADMIN"), payrollHandler.CalculatePayslip)

	recom := protected.Group("/student-recommendations")
	// Hanya Guru BK dan Kepala Sekolah yang bisa mengaktifkan mesin kalkulasi ini
	recom.Post("/calculate", middlewares.RequireRoles("BK", "KEPSEK", "ADMIN"), recomHandler.CalculateRanking)

	// Pengumuman
	announcements := protected.Group("/announcements")
	announcements.Post("/", middlewares.RequireRoles("ADMIN", "TU", "KEPSEK"), annHandler.CreateAnnouncement)

	// Laporan (Reports)
	reports := protected.Group("/reports")
	reports.Post("/academic", middlewares.RequireRoles("ADMIN", "KEPSEK"), annHandler.RequestAcademicReport)

	log.Println("🚀 Server School System berjalan di port 3000")
	app.Listen(":3000")
}
