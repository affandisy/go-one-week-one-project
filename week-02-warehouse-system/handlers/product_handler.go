package handlers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/services"
	"github.com/gofiber/fiber/v2"
)

type ProductHandler struct {
	service services.ProductService
}

func NewProductHandler(service services.ProductService) *ProductHandler {
	return &ProductHandler{service: service}
}

func (h *ProductHandler) CreateProduct(c *fiber.Ctx) error {
	var req services.CreateProductRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Format request tidak valid",
			"error":   err.Error(),
		})
	}

	product, err := h.service.CreateProduct(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Gagal membuat produk",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"success": true,
		"message": "Produk berhasil ditambahkan",
		"data":    product,
	})
}

func (h *ProductHandler) GetAllProducts(c *fiber.Ctx) error {
	products, err := h.service.GetAllProducts()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data produk",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil mengambil data produk",
		"data":    products,
	})
}

func (h *ProductHandler) GetProductsPaginated(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)

	result, err := h.service.GetAllProductsPaginated(page, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Gagal mengambil data produk",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Berhasil memuat halaman produk",
		"data":    result,
	})
}

func (h *ProductHandler) UploadImage(c *fiber.Ctx) error {
	productID, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID Produk tidak valid",
		})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File gambar tidak ditemukan",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".png" && ext != ".jpg" && ext != ".jpeg" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format file harus PNG atau JPG",
		})
	}

	if file.Size > 2*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Ukuran gambar maksimal 2MB",
		})
	}

	filename := fmt.Sprintf("product_%d_%d%s", productID, time.Now().Unix(), ext)
	savePath := fmt.Sprintf("./uploads/products/%s", filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan gambar",
		})
	}

	dbPath := fmt.Sprintf("/uploads/products/%s", filename)
	if err := h.service.UpdateProductImage(uint(productID), dbPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Gambar produk berhasil diunggah",
		"image_url": dbPath,
	})
}
