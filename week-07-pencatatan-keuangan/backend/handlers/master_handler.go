// handlers/master_handler.go
package handlers

import (
	"github.com/affandisy/financial-app/services"
	"github.com/gofiber/fiber/v2"
)

type MasterHandler struct {
	service services.MasterService
}

func NewMasterHandler(service services.MasterService) *MasterHandler {
	return &MasterHandler{service}
}

func (h *MasterHandler) GetWallet(c *fiber.Ctx) error {
	wallet, err := h.service.GetWalletInfo()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data dompet"})
	}
	return c.JSON(fiber.Map{"data": wallet})
}

func (h *MasterHandler) GetCategories(c *fiber.Ctx) error {
	txType := c.Query("type") // Menangkap ?type=income atau ?type=expense
	if txType == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Parameter type (income/expense) wajib diisi"})
	}

	categories, err := h.service.GetActiveCategories(txType)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil daftar kategori"})
	}
	return c.JSON(fiber.Map{"data": categories})
}
