package rest

import (
	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type BillingHandler struct {
	useCase port.BillingUseCase
}

func NewBillingHandler(useCase port.BillingUseCase) *BillingHandler {
	return &BillingHandler{useCase}
}

type CreateInvoiceReq struct {
	OwnerID string `json:"owner_id"`
	Items   []struct {
		PetID     string  `json:"pet_id"`
		ServiceID string  `json:"service_id"`
		Price     float64 `json:"price"`
	} `json:"items"`
}

func (h *BillingHandler) CreateInvoice(c *fiber.Ctx) error {
	var req CreateInvoiceReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	// Mapper Request -> Domain
	var inputItems []domain.InvoiceItem
	for _, it := range req.Items {
		inputItems = append(inputItems, domain.InvoiceItem{
			PetID:     it.PetID,
			ServiceID: it.ServiceID,
			Price:     it.Price,
		})
	}

	invoice, err := h.useCase.GenerateInvoice(req.OwnerID, inputItems)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"data": invoice})
}
