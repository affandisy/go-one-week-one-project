package rest

import (
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	useCase port.AnalyticsUseCase
}

func NewAnalyticsHandler(useCase port.AnalyticsUseCase) *AnalyticsHandler {
	return &AnalyticsHandler{useCase}
}

func (h *AnalyticsHandler) GetPetNutritionSummary(c *fiber.Ctx) error {
	petID := c.Params("pet_id")

	summary, err := h.useCase.GeneratePetNutritionReport(petID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Analitik gizi berhasil dimuat",
		"data":    summary,
	})
}
