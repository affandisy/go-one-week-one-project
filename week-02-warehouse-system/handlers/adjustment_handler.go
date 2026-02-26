package handlers

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/services"
	"github.com/gofiber/fiber/v2"
)

type AdjustmentHandler struct {
	service services.AdjustmentService
}

func NewAdjustmentHandler(service services.AdjustmentService) *AdjustmentHandler {
	return &AdjustmentHandler{
		service: service,
	}
}

func (h *AdjustmentHandler) Create(c *fiber.Ctx) error {
	var req services.CreateAdjustmentRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	userID := c.Locals("user_id").(uint)

	result, err := h.service.ProcessAdjustment(req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Penyesuaian stok berhasil diproses",
		"data":    result,
	})
}
