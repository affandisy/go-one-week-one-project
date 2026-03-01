package handlers

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/services"
	"github.com/gofiber/fiber/v2"
)

type LocationHandler struct {
	service services.LocationService
}

func NewLocationHandler(service services.LocationService) *LocationHandler {
	return &LocationHandler{service: service}
}

func (h *LocationHandler) Create(c *fiber.Ctx) error {
	var req services.CreateLocationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request tidak valid",
		})
	}

	location, err := h.service.CreateLocation(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Lokasi/Rak berhasil ditambahkan",
		"data":    location,
	})
}

func (h *LocationHandler) GetAll(c *fiber.Ctx) error {
	locations, err := h.service.GetAllLocations()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil memuat daftar lokasi",
		"data":    locations,
	})
}
