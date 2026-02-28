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
