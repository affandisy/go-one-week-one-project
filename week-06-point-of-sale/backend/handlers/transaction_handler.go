package handlers

import (
	"github.com/affandisy/pos-system/services"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	service services.TransactionService
}

func NewTransactionHandler(service services.TransactionService) *TransactionHandler {
	return &TransactionHandler{service}
}

func (h *TransactionHandler) Checkout(c *fiber.Ctx) error {
	// Ambil ID Kasir dari Token JWT
	cashierID := c.Locals("user_id").(string)

	var req services.CheckoutRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	transaction, err := h.service.Checkout(cashierID, req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Transaksi berhasil",
		"data": fiber.Map{
			"receipt_number": transaction.ReceiptNumber,
			"final_amount":   transaction.FinalAmount,
			"change_amount":  transaction.ChangeAmount,
			"created_at":     transaction.CreatedAt,
		},
	})
}
