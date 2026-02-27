package handlers

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/services"
	"github.com/gofiber/fiber/v2"
)

type PartnerHandler struct {
	service services.PartnerService
}

func NewPartnerHandler(service services.PartnerService) *PartnerHandler {
	return &PartnerHandler{service: service}
}

func (h *PartnerHandler) Create(c *fiber.Ctx) error {
	var req services.CreatePartnerRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format tidak valid",
		})
	}

	partner, err := h.service.CreatePartner(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Partner berhasil ditambahkan",
		"data":    partner,
	})
}

func (h *PartnerHandler) GetAll(c *fiber.Ctx) error {
	partners, err := h.service.GetAllPartners()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": partners,
	})
}
