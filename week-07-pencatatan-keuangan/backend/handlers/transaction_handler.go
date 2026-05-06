package handlers

import (
	"github.com/affandisy/financial-app/services"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service}
}

func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	var req services.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.RecordTransaction(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Transaksi berhasil disimpan",
	})
}

func (h *TransactionHandler) GetHistory(c *fiber.Ctx) error {
	history, err := h.service.GetDashboardHistory()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat riwayat"})
	}
	return c.JSON(fiber.Map{"data": history})
}

func (h *TransactionHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteTransaction(id); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Transaksi dihapus"})
}

func (h *TransactionHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req services.TransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format invalid"})
	}
	if err := h.service.UpdateTransaction(id, req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Transaksi diperbarui"})
}
