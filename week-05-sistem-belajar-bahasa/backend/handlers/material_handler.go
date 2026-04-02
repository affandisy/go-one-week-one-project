package handlers

import (
	"github.com/affandi/belajar-bahasa/services"
	"github.com/gofiber/fiber/v2"
)

type MaterialHandler struct {
	service services.MaterialService
}

func NewMaterialHandler(service services.MaterialService) *MaterialHandler {
	return &MaterialHandler{service}
}

func (h *MaterialHandler) GetByModule(c *fiber.Ctx) error {
	moduleID := c.Params("moduleId")
	materials, err := h.service.GetMaterialsByModule(moduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil materi"})
	}

	return c.JSON(fiber.Map{"data": materials})
}

func (h *MaterialHandler) Create(c *fiber.Ctx) error {
	var req services.CreateMaterialRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.CreateMaterial(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Materi berhasil ditambahkan"})
}
