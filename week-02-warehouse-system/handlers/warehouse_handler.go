package handlers

import (
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
	"github.com/gofiber/fiber/v2"
)

type WarehouseHandler struct {
	repo repositories.WarehouseRepository
}

func NewWarehouseHandler(repo repositories.WarehouseRepository) *WarehouseHandler {
	return &WarehouseHandler{repo: repo}
}

func (h *WarehouseHandler) Create(c *fiber.Ctx) error {
	var wh models.Warehouse

	if err := c.BodyParser(&wh); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format tidak valid",
		})
	}

	if wh.Code == "" || wh.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Kode dan Nama Gudang wajib diisi",
		})
	}

	if err := h.repo.Create(&wh); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan gudang",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Gudang berhasil ditambahkan",
		"data":    wh,
	})
}

func (h *WarehouseHandler) GetAll(c *fiber.Ctx) error {
	warehouses, err := h.repo.FindAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat data gudang"})
	}
	return c.JSON(fiber.Map{"data": warehouses})
}
