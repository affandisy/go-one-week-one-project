package handlers

import (
	"github.com/affandisy/pos-system/services"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	service services.ReportService
}

func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service}
}

func (h *ReportHandler) GetDailySales(c *fiber.Ctx) error {
	// Opsional: Ambil parameter tanggal dari URL (contoh: ?date=2026-04-23)
	dateParam := c.Query("date")

	summary, err := h.service.GetDailySalesSummary(dateParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal tidak valid. Gunakan YYYY-MM-DD"})
	}

	return c.JSON(fiber.Map{
		"message": "Laporan penjualan harian berhasil dimuat",
		"data":    summary,
	})
}

func (h *ReportHandler) GetLowStocks(c *fiber.Ctx) error {
	products, err := h.service.GetLowStockAlerts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat laporan stok"})
	}

	return c.JSON(fiber.Map{
		"message": "Laporan stok menipis berhasil dimuat",
		"data":    products,
	})
}
