package handlers

import (
	"github.com/affandisy/padel-booking-system/services"
	"github.com/gofiber/fiber/v2"
)

type CourtHandler struct {
	service services.CourtService
}

func NewCourtHandler(service services.CourtService) *CourtHandler {
	return &CourtHandler{service}
}

func (h *CourtHandler) Create(c *fiber.Ctx) error {
	var req services.CreateCourtRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.CreateCourt(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Lapangan berhasil ditambahkan"})
}

func (h *CourtHandler) GetAll(c *fiber.Ctx) error {
	courts, err := h.service.GetAllCourts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data lapangan"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": courts})
}

func (h *CourtHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	court, err := h.service.GetCourtByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": court})
}

func (h *CourtHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req services.UpdateCourtRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.UpdateCourt(id, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Data lapangan berhasil diperbarui"})
}

func (h *CourtHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteCourt(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Lapangan berhasil dihapus"})
}
