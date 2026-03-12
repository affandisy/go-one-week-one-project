package handlers

import (
	"github.com/affandisy/school-system-sma/service"
	"github.com/gofiber/fiber/v2"
)

type AcademicHandler struct {
	service service.AcademicService
}

func NewAcademicHandler(service service.AcademicService) *AcademicHandler {
	return &AcademicHandler{service}
}

func (h *AcademicHandler) CreateAcademicYear(c *fiber.Ctx) error {
	var req service.CreateAcademicYearRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	ay, err := h.service.CreateAcademicYear(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Tahun ajaran berhasil dibuat",
		"data":    ay,
	})
}

func (h *AcademicHandler) GetAcademicYears(c *fiber.Ctx) error {
	ays, err := h.service.GetAllAcademicYears()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data"})
	}
	return c.JSON(fiber.Map{"data": ays})
}
