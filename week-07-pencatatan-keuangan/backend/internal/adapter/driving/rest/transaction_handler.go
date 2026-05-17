package rest

import (
	"github.com/affandisy/financial-app/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	useCase port.TransactionUseCase
}

func NewTransactionHandler(useCase port.TransactionUseCase) *TransactionHandler {
	return &TransactionHandler{useCase}
}

type RecordTrxReq struct {
	WalletID   string  `json:"wallet_id"`
	Type       string  `json:"type"`
	CategoryID string  `json:"category_id"`
	Amount     float64 `json:"amount"`
	Note       string  `json:"note"`
}

func (h *TransactionHandler) Create(c *fiber.Ctx) error {
	var req RecordTrxReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.useCase.RecordTransaction(req.WalletID, req.Type, req.CategoryID, req.Amount, req.Note); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Transaksi berhasil disimpan"})
}

func (h *TransactionHandler) GetRecent(c *fiber.Ctx) error {
	walletID := c.Query("wallet_id")
	if walletID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "wallet_id wajib diisi"})
	}

	history, err := h.useCase.GetRecentHistory(walletID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat riwayat"})
	}

	return c.JSON(fiber.Map{"data": history})
}
