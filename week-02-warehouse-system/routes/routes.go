package routes

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/handlers"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/middlewares"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, productHandler *handlers.ProductHandler, userHandler *handlers.UserHandler, txHandler *handlers.TransactionHandler, dashboardHandler *handlers.DashboardHandler, adjHandler *handlers.AdjustmentHandler, analyticsHandler *handlers.AnalyticsHandler, partnerHandler *handlers.PartnerHandler, opnameHandler *handlers.OpnameHandler, reportHandler *handlers.ReportHandler, locationHandler *handlers.LocationHandler) {
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
	products.Get("/page", productHandler.GetProductsPaginated)

	transactions := protected.Group("/transactions")
	transactions.Post("/", txHandler.Create)

	transactions.Put("/:id/approve", middlewares.RequireRoles("admin", "manager"), txHandler.Approve)

	adjustments := protected.Group("/adjustments")
	adjustments.Post("/", adjHandler.Create)

	analytics := protected.Group("/analytics", middlewares.RequireRoles("admin", "manager"))
	analytics.Get("/best-sellers", analyticsHandler.GetBestSellers)

	partners := protected.Group("/partners")
	partners.Post("/", partnerHandler.Create)
	partners.Get("/", partnerHandler.GetAll)

	opname := protected.Group("/opname")
	opname.Post("/", middlewares.RequireRoles("admin", "manager"), opnameHandler.ProcessOpname)

	reports := protected.Group("/reports", middlewares.RequireRoles("admin", "manager"))
	reports.Get("/transactions/csv", reportHandler.DonwloadMonthlyCSV)
	reports.Get("/transactions/batch-csv", reportHandler.TriggerExportCSV)

	locations := protected.Group("/locations")
	locations.Post("/", middlewares.RequireRoles("admin", "manager"), locationHandler.Create)
	locations.Get("/", locationHandler.GetAll)

}
