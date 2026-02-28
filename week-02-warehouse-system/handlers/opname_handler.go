package handlers

import (
	"fmt"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/services"
	"github.com/gofiber/fiber/v2"
)

type OpnameHandler struct {
	service services.OpnameService
}

func NewOpnameHandler(service services.OpnameService) *OpnameHandler {
	return &OpnameHandler{service: service}
}

func (h *OpnameHandler) ProcessOpname(c *fiber.Ctx) error {
	var req services.ProcessOpnameReq

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request tidak valid",
		})
	}

	userID := c.Locals("user_id").(uint)

	discrepancyCount, err := h.service.ProcessMassOpname(req, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	msg := fmt.Sprintf("Opname selesai. Ditemukan selisih pada %d barang dan stok berhasil disesuaikan.", discrepancyCount)
	if discrepancyCount == 0 {
		msg = "Opname selesai. Stok fisik 100% cocok dengan sistem, tidak ada penyesuaian yang dilakukan."
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data": fiber.Map{
			"item_processed": len(req.Items),
			"item_adjusted":  discrepancyCount,
		},
	})
}
