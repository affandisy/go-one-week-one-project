package handlers

import (
	"strconv"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/services"
	"github.com/gofiber/fiber/v2"
)

type AnalyticsHandler struct {
	service services.AnalyticsService
}

func NewAnalyticsHandler(service services.AnalyticsService) *AnalyticsHandler {
	return &AnalyticsHandler{service: service}
}

func (h *AnalyticsHandler) GetBestSellers(c *fiber.Ctx) error {
	days, err := strconv.Atoi(c.Query("days", "30"))
	if err != nil || days < 1 {
		days = 30
	}

	limit, err := strconv.Atoi(c.Query("limit", "5"))
	if err != nil || limit < 1 {
		limit = 5
	}

	results, err := h.service.GetTopBestSellers(days, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal menghitung data analitik",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil memuat Top Best Seller",
		"data": fiber.Map{
			"periode_hari": days,
			"limit":        limit,
			"items":        results,
		},
	})
}
