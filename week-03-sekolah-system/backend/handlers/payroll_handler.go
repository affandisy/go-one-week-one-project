package handlers

import (
	"github.com/affandisy/school-system-sma/service"
	"github.com/gofiber/fiber/v2"
)

type PayrollHandler struct {
	service service.PayrollService
}

func NewPayrollHandler(service service.PayrollService) *PayrollHandler {
	return &PayrollHandler{service: service}
}

func (h *PayrollHandler) CalculatePayslip(c *fiber.Ctx) error {
	var req struct {
		TeacherID string `json:"teacher_id"`
		Month     int    `json:"month"`
		Year      int    `json:"year"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	payslip, err := h.service.CalculatePayslip(req.TeacherID, req.Month, req.Year)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Slip gaji berhasil dihitung",
		"data":    payslip,
	})
}
