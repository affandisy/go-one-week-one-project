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

	// Inisialisasi Modul Absensi
	attRepo := repository.NewAttendanceRepository(db)
	attService := service.NewAttendanceService(attRepo)
	attHandler := handlers.NewAttendanceHandler(attService)

	masterRepo := repository.NewMasterRepository(db)
	masterService := service.NewMasterService(masterRepo)
	masterHandler := handlers.NewMasterHandler(masterService)

	reportRepo := repository.NewReportRepository(db)
	reportService := service.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

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
	announcements.Get("/", annHandler.GetAnnouncements)

	// Laporan (Reports)
	reports := protected.Group("/reports")
	reports.Post("/academic", middlewares.RequireRoles("ADMIN", "KEPSEK"), annHandler.RequestAcademicReport)

	attendance := protected.Group("/attendance")
	// Sesuai PRD, Guru yang mencatat absensi murid
	attendance.Post("/student/batch", middlewares.RequireRoles("GURU", "ADMIN"), attHandler.RecordBatchAttendance)
	attendance.Get("/class/:classId/students", middlewares.RequireRoles("GURU", "ADMIN"), attHandler.GetClassStudents)

	// CLASSES
	classes := protected.Group("/classes")
	classes.Get("/", masterHandler.GetClasses) // Semua bisa melihat daftar kelas
	classes.Post("/", middlewares.RequireRoles("ADMIN", "TU"), masterHandler.CreateClass)

	// SUBJECTS
	subjects := protected.Group("/subjects")
	subjects.Get("/", masterHandler.GetSubjects) // Semua bisa melihat daftar mapel
	subjects.Post("/", middlewares.RequireRoles("ADMIN", "TU"), masterHandler.CreateSubject)

	// STUDENTS
	students := protected.Group("/students")
	students.Get("/", middlewares.RequireRoles("ADMIN", "TU", "GURU", "KEPSEK"), masterHandler.GetStudents)
	students.Post("/", middlewares.RequireRoles("ADMIN", "TU"), masterHandler.CreateStudent)

	reportsCard := protected.Group("/report-cards")
	reportsCard.Get("/", reportHandler.GetReports)
	reportsCard.Post("/generate", middlewares.RequireRoles("KEPSEK", "ADMIN", "GURU"), reportHandler.GenerateReport)

	log.Println("Server School System berjalan di port 3000")
	app.Listen(":3000")
}
