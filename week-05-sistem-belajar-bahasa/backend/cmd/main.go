package main

import (
	"log"

	"github.com/affandi/belajar-bahasa/config"
	"github.com/affandi/belajar-bahasa/handlers"
	"github.com/affandi/belajar-bahasa/middlewares"
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

	jwtSecret := "test"

	userRepo := repositories.NewUserRepository(db)
	authService := services.NewAuthService(userRepo, jwtSecret)
	authHandler := handlers.NewAuthHandler(authService)

	moduleRepo := repositories.NewModuleRepository(db)
	moduleService := services.NewModuleService(moduleRepo)
	moduleHandler := handlers.NewModuleHandler(moduleService)

	materialRepo := repositories.NewMaterialRepository(db)
	materialService := services.NewMaterialService(materialRepo)
	materialHandler := handlers.NewMaterialHandler(materialService)
	uploadHandler := handlers.NewUploadHandler()

	progressRepo := repositories.NewProgressRepository(db)
	quizService := services.NewQuizService(materialRepo, progressRepo, moduleRepo)
	quizHandler := handlers.NewQuizHandler(quizService)

	progressService := services.NewProgressService(progressRepo, moduleRepo)
	progressHandler := handlers.NewProgressHandler(progressService)

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

	api.Post("/modules/seed", moduleHandler.Create)

	protected := api.Group("/", middlewares.JWTProtected(jwtSecret))

	modules := protected.Group("/modules")
	modules.Get("/", moduleHandler.GetAll)
	modules.Get("/:id", moduleHandler.GetByID)

	modules.Get("/:moduleId/materials", materialHandler.GetByModule) // Lihat isi kartu modul
	modules.Post("/:moduleId/materials", materialHandler.Create)     // Tambah kartu ke modul

	protected.Post("/upload", uploadHandler.UploadFile)

	modules.Get("/:moduleId/quiz", quizHandler.StartQuiz)          // Mendapatkan soal
	modules.Post("/:moduleId/quiz/submit", quizHandler.SubmitQuiz) // Mengirim jawaban

	protected.Get("/progress", progressHandler.GetProgress)
	protected.Post("/progress/reset", progressHandler.ResetProgress)

	// 5. Jalankan Server
	log.Println("🚀 Server Language Learning berjalan di port 3000")
	log.Fatal(app.Listen(":3000"))
}
