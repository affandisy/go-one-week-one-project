package handlers

import (
	"github.com/affandisy/pos-system/services"
	"github.com/gofiber/fiber/v2"
)

type SupplierHandler struct {
	service services.SupplierService
}

func NewSupplierHandler(service services.SupplierService) *SupplierHandler {
	return &SupplierHandler{service}
}

func (h *SupplierHandler) GetAll(c *fiber.Ctx) error {
	suppliers, err := h.service.GetAllSuppliers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data pemasok"})
	}
	return c.JSON(fiber.Map{"data": suppliers})
}

func (h *SupplierHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	supplier, err := h.service.GetSupplierByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(fiber.Map{"data": supplier})
}

func (h *SupplierHandler) Create(c *fiber.Ctx) error {
	var req services.SupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.CreateSupplier(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Pemasok berhasil ditambahkan"})
}

func (h *SupplierHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req services.SupplierRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.UpdateSupplier(id, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data pemasok berhasil diperbarui"})
}

func (h *SupplierHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteSupplier(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Pemasok berhasil dihapus"})
}
