package routes

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/handlers"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, productHandler *handlers.ProductHandler, userHandler *handlers.UserHandler, txHandler *handlers.TransactionHandler, dashboardHandler *handlers.DashboardHandler) {
	api := app.Group("/api/v1")

	auth := api.Group("/auth")
	auth.Post("/register", userHandler.Register)
	auth.Post("/login", userHandler.Login)

	protected := api.Group("/", middlewares.Protected())

	dashboard := protected.Group("/dashboard")
	dashboard.Get("/", dashboardHandler.GetSummary)

	products := protected.Group("/products")
	products.Post("/", productHandler.CreateProduct)
	products.Get("/", productHandler.GetAllProducts)

	transactions := protected.Group("/transactions")
	transactions.Post("/", txHandler.Create)
}
