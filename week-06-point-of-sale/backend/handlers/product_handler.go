package handlers

import (
	"github.com/affandisy/pos-system/services"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service}
}

func (h *ProductHandler) GetAll(c *fiber.Ctx) error {
	// Menangkap query parameter ?search=
	search := c.Query("search")

	products, err := h.service.GetAllProducts(search)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data produk"})
	}

	return c.JSON(fiber.Map{"data": products})
}

func (h *ProductHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	product, err := h.service.GetProductByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Produk tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"data": product})
}

func (h *ProductHandler) Create(c *fiber.Ctx) error {
	var req services.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.CreateProduct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Produk berhasil ditambahkan"})
}

func (h *ProductHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req services.ProductRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.UpdateProduct(id, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Produk berhasil diperbarui"})
}

func (h *ProductHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")
	if err := h.service.DeleteProduct(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Produk berhasil dihapus"})
}
