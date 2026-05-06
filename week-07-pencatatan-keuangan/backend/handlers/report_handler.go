package handlers

import (
	"strconv"
	"time"

	"github.com/affandisy/financial-app/services"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	service services.ReportService
}

func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service}
}

func (h *ReportHandler) GetMonthlySummary(c *fiber.Ctx) error {
	// Tangkap query ?year=2026&month=5
	yearStr := c.Query("year")
	monthStr := c.Query("month")

	// Fallback ke bulan dan tahun ini jika parameter kosong
	now := time.Now()
	year := now.Year()
	month := now.Month()

	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = y
		}
	}
	if monthStr != "" {
		if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
			month = time.Month(m)
		}
	}

	report, err := h.service.GetMonthlyReport(year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghasilkan laporan"})
	}

	return c.JSON(fiber.Map{"data": report})
}
