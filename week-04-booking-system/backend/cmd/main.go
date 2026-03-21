package main

import (
	"log"

	"github.com/affandisy/padel-booking-system/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// 1. Inisialisasi Database
	config.ConnectDatabase()

	// 2. Setup Fiber
	app := fiber.New()
	app.Use(logger.New())

	// 3. Routing Dasar (Base URL /api/v1)[cite: 2]
	api := app.Group("/api/v1")

	api.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "OK", "message": "Padel Booking API is running"})
	})

	log.Println("Server Padel Booking berjalan di port 3000")
	log.Fatal(app.Listen(":3000"))
}
