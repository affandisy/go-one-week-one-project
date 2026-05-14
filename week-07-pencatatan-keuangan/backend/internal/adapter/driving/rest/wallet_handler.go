package rest

import (
	"github.com/affandisy/financial-app/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	useCase port.WalletUseCase
}

func NewWalletHandler(useCase port.WalletUseCase) *WalletHandler {
	return &WalletHandler{useCase}
}

// Request DTO
type CreateWalletReq struct {
	Name           string  `json:"name"`
	InitialBalance float64 `json:"initial_balance"`
}

func (h *WalletHandler) Create(c *fiber.Ctx) error {
	var req CreateWalletReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid payload"})
	}

	if err := h.useCase.CreateWallet(req.Name, req.InitialBalance); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Dompet berhasil dibuat"})
}

func (h *WalletHandler) GetAll(c *fiber.Ctx) error {
	wallets, err := h.useCase.ListWallets()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "internal error"})
	}

	// Response DTO mapping bisa dilakukan di sini jika diperlukan
	return c.JSON(fiber.Map{"data": wallets})
}
