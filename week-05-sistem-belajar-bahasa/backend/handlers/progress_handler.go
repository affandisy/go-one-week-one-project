package handlers

import (
	"github.com/affandi/belajar-bahasa/services"
	"github.com/gofiber/fiber/v2"
)

type ProgressHandler struct {
	service services.ProgressService
}

func NewProgressHandler(service services.ProgressService) *ProgressHandler {
	return &ProgressHandler{service}
}

func (h *ProgressHandler) GetProgress(c *fiber.Ctx) error {
	// Ambil ID User dari Token JWT (disimpan di Locals oleh Middleware)
	userID := c.Locals("user_id").(string)

	dashboard, err := h.service.GetDashboardProgress(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat dashboard progres"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Dashboard progres berhasil dimuat",
		"data":    dashboard,
	})
}

func (h *ProgressHandler) ResetProgress(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	if err := h.service.ResetProgress(userID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mereset progres"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Progres belajar Anda telah berhasil direset dari awal",
	})
}
