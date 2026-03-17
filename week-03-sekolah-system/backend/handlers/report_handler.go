package handlers

import (
	"github.com/affandisy/school-system-sma/service"
	"github.com/gofiber/fiber/v2"
)

type ReportHandler struct {
	service service.ReportService
}

func NewReportHandler(s service.ReportService) *ReportHandler {
	return &ReportHandler{s}
}

func (h *ReportHandler) GenerateReport(c *fiber.Ctx) error {
	var req struct {
		StudentID      string `json:"student_id"`
		Semester       int    `json:"semester"`
		AcademicYearID string `json:"academic_year_id"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tidak valid"})
	}

	err := h.service.TriggerReportGeneration(req.StudentID, req.Semester, req.AcademicYearID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Sesuai PRD, respons instan karena background job
	return c.Status(202).JSON(fiber.Map{"message": "Pembuatan rapor sedang diproses di latar belakang."})
}

func (h *ReportHandler) GetReports(c *fiber.Ctx) error {
	data, _ := h.service.GetReports()
	return c.JSON(fiber.Map{"data": data})
}
