package rest

import (
	"strconv"
	"time"

	"github.com/affandisy/financial-app/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	useCase port.ReportUseCase
}

func NewReportHandler(useCase port.ReportUseCase) *ReportHandler {
	return &ReportHandler{useCase}
}

func (h *ReportHandler) GetMonthlySummary(c *fiber.Ctx) error {
	walletID := c.Query("wallet_id")
	if walletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id wajib diisi"})
	}

	yearStr := c.Query("year")
	monthStr := c.Query("month")

	now := time.Now()
	year := now.Year()
	month := now.Month()

	if y, err := strconv.Atoi(yearStr); err == nil {
		year = y
	}
	if m, err := strconv.Atoi(monthStr); err == nil && m >= 1 && m <= 12 {
		month = time.Month(m)
	}

	// Tembak Port (Interface)
	report, err := h.useCase.GetMonthlySummary(walletID, year, month)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghasilkan laporan"})
	}

	return c.JSON(fiber.Map{"data": report})
}
