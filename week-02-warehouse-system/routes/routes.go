package routes

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/handlers"
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App, productHandler *handlers.ProductHandler) {
	api := app.Group("/api/v1")

	products := api.Group("/products")
	products.Post("/", productHandler.CreateProduct)
	products.Get("/", productHandler.GetAllProducts)
}
