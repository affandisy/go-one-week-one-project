package handlers

import (
	"fmt"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/services"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	service services.ReportService
}

func NewReportHandler(service services.ReportService) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) DonwloadMonthlyCSV(c *fiber.Ctx) error {
	month := c.QueryInt("month", int(time.Now().Month()))
	year := c.QueryInt("year", time.Now().Year())

	csvData, err := h.service.GenerateMonthlyTransactionCSV(month, year)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal membuat laporan CSV",
		})
	}

	filename := fmt.Sprintf("Laporan_Transaksi_Gudang_%02d_%d.csv", month, year)

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return c.Send(csvData)
}

func (h *ReportHandler) TriggerExportCSV(c *fiber.Ctx) error {
	month := c.QueryInt("month", int(time.Now().Month()))
	year := c.QueryInt("year", time.Now().Year())
	email := "manager@perusahaan.com"

	err := h.service.TriggerMonthlyTransactionExport(month, year, email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal memulai pembuatan laporan",
		})
	}

	return c.Status(fiber.StatusAccepted).JSON(fiber.Map{
		"success": true,
		"message": "Permintaan export laporan sedang diproses. Tautan unduhan akan dikirimkan ke email Anda saat selesai.",
	})
}

func (h *ReportHandler) GetMovementAnalytics(c *fiber.Ctx) error {
	month := c.QueryInt("month", int(time.Now().Month()))
	year := c.QueryInt("year", time.Now().Year())

	results, err := h.service.GetMovementAnalytics(month, year)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"error":   "Gagal menarik data analitik pergerakan barang",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": fmt.Sprintf("Laporan pergerakan barang periode %02d-%d", month, year),
		"data":    results,
	})
}
