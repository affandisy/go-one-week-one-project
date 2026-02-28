package main

import (
	"log"

	"github.com/gofiber/fiber/v2/middleware/logger"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/config"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/handlers"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/routes"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/services"
	"github.com/gofiber/fiber/v2"
)

func main() {
	config.ConnectDatabase()
	db := config.DB

	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	userRepo := repositories.NewUserRepository(db)
	userService := services.NewUserService(userRepo)
	userHandler := handlers.NewUserHandler(userService)

	partnerRepo := repositories.NewPartnerRepository(db)
	partnerService := services.NewPartnerService(partnerRepo)
	partnerHandler := handlers.NewPartnerHandler(partnerService)

	txRepo := repositories.NewTransactionRepository(db)
	txService := services.NewTransactionService(txRepo, productRepo, partnerRepo)
	txHandler := handlers.NewTransactionHandler(txService)

	dashboardService := services.NewDashboardService(productRepo, txRepo)
	dashboardHandler := handlers.NewDashboardHandler(dashboardService)

	adjRepo := repositories.NewAdjustmentRepository(db)
	adjService := services.NewAdjustmentService(adjRepo, productRepo)
	adjHandler := handlers.NewAdjustmentHandler(adjService)

	analyticsService := services.NewAnalyticsService(txRepo)
	analyticsHandler := handlers.NewAnalyticsHandler(analyticsService)

	opnameService := services.NewOpnameService(adjRepo, productRepo)
	opnameHandler := handlers.NewOpnameHandler(opnameService)

	reportService := services.NewReportService(txRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	app := fiber.New()
	app.Use(logger.New())

	// app.Get("/", func(c *fiber.Ctx) error {
	// 	return c.JSON(fiber.Map{
	// 		"message": "Selamat datang di Warehouse Management System",
	// 		"status":  "OK",
	// 	})
	// })

	routes.SetupRoutes(app, productHandler, userHandler, txHandler, dashboardHandler, adjHandler, analyticsHandler, partnerHandler, opnameHandler, reportHandler)

	log.Println("Server WMS berjalan di http://localhost:3000")
	err := app.Listen(":3000")
	if err != nil {
		log.Fatal("Server error: ", err)
	}
}
