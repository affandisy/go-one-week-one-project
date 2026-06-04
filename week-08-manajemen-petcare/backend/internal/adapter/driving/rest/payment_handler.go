package rest

import (
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type PaymentHandler struct {
	useCase port.PaymentUseCase
}

func NewPaymentHandler(useCase port.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase}
}

type ProcessPaymentReq struct {
	InvoiceID string  `json:"invoice_id"`
	Method    string  `json:"method"`
	Amount    float64 `json:"amount"`
	Reference string  `json:"reference"`
}

func (h *PaymentHandler) Process(c *fiber.Ctx) error {
	var req ProcessPaymentReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	// Ambil userID dari context Fiber (hasil ekstraksi JWT oleh Middleware)
	userID := c.Locals("userID").(string)

	// Teruskan userID ke UseCase
	payment, err := h.useCase.ProcessPayment(userID, req.InvoiceID, req.Method, req.Amount, req.Reference)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pembayaran berhasil diproses",
		"data":    payment,
	})
}
