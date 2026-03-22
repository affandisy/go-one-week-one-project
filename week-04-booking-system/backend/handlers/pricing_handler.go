package handlers

import (
	"github.com/affandisy/padel-booking-system/services"
	"github.com/gofiber/fiber/v2"
)

type PricingHandler struct {
	service services.PricingService
}

func NewPricingHandler(service services.PricingService) *PricingHandler {
	return &PricingHandler{service}
}

func (h *PricingHandler) Create(c *fiber.Ctx) error {
	var req services.CreatePricingRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.CreateRule(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Aturan harga berhasil ditambahkan"})
}

func (h *PricingHandler) GetByCourt(c *fiber.Ctx) error {
	// Ambil courtID dari query params (misal: /pricing-rules?court_id=123)
	courtID := c.Query("court_id")
	if courtID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Parameter court_id wajib disertakan"})
	}

	rules, err := h.service.GetRulesByCourt(courtID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data aturan harga"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": rules})
}

func (h *PricingHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req services.UpdatePricingRuleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.UpdateRule(id, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Aturan harga berhasil diperbarui"})
}

func (h *PricingHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteRule(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Aturan harga berhasil dihapus"})
}
