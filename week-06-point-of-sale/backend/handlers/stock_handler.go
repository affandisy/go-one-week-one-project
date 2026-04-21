package handlers

import (
	"github.com/affandisy/pos-system/services"
	"github.com/gofiber/fiber/v2"
)

type StockHandler struct {
	service services.StockService
}

func NewStockHandler(service services.StockService) *StockHandler {
	return &StockHandler{service}
}

func (h *StockHandler) StockIn(c *fiber.Ctx) error {
	var req services.StockInRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.RecordStockIn(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Stok berhasil ditambahkan",
	})
}

func (h *StockHandler) GetHistory(c *fiber.Ctx) error {
	productID := c.Params("productId")

	history, err := h.service.GetStockHistory(productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat riwayat stok"})
	}

	return c.JSON(fiber.Map{"data": history})
}
