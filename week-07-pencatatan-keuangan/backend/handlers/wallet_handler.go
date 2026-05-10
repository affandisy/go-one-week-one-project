package handlers

import (
	"github.com/affandisy/financial-app/services"
	"github.com/gofiber/fiber/v2"
)

type WalletHandler struct {
	service services.WalletService
}

func NewWalletHandler(service services.WalletService) *WalletHandler {
	return &WalletHandler{service}
}

func (h *WalletHandler) GetAll(c *fiber.Ctx) error {
	wallets, err := h.service.GetAllWallets()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil daftar dompet"})
	}
	return c.JSON(fiber.Map{"data": wallets})
}

func (h *WalletHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	wallet, err := h.service.GetWallet(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": wallet})
}

func (h *WalletHandler) Create(c *fiber.Ctx) error {
	var req services.WalletRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	if err := h.service.CreateWallet(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Dompet berhasil dibuat"})
}

func (h *WalletHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req services.WalletRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format data tidak valid"})
	}

	if err := h.service.UpdateWallet(id, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Dompet berhasil diperbarui"})
}

func (h *WalletHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteWallet(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"message": "Dompet berhasil dihapus"})
}
