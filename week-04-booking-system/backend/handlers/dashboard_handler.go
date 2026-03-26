package handlers

import (
	"github.com/affandisy/padel-booking-system/services"
	"github.com/gofiber/fiber/v2"
)

type DashboardHandler struct {
	service services.DashboardService
}

func NewDashboardHandler(service services.DashboardService) *DashboardHandler {
	return &DashboardHandler{service}
}

// GetStats godoc
// @Summary Dapatkan Statistik Dashboard Admin
// @Description Mengembalikan ringkasan pendapatan, booking aktif, dan okupansi lapangan
// @Tags Dashboard
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} map[string]interface{} "Statistik berhasil diambil"
// @Failure 500 {object} map[string]string "Terjadi kesalahan server"
// @Router /admin/dashboard-stats [get]
func (h *DashboardHandler) GetStats(c *fiber.Ctx) error {
	stats, err := h.service.GetDashboardStats()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data statistik dashboard",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data statistik dashboard berhasil diambil",
		"data":    stats,
	})
}
