package handlers

import (
	"github.com/affandi/belajar-bahasa/services"
	"github.com/gofiber/fiber/v2"
)

type ModuleHandler struct {
	service services.ModuleService
}

func NewModuleHandler(service services.ModuleService) *ModuleHandler {
	return &ModuleHandler{service}
}

func (h *ModuleHandler) GetAll(c *fiber.Ctx) error {
	modules, err := h.service.GetAllModules()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data modul"})
	}

	return c.JSON(fiber.Map{"data": modules})
}

func (h *ModuleHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	module, err := h.service.GetModuleByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Modul tidak ditemukan"})
	}

	return c.JSON(fiber.Map{"data": module})
}

// Hanya untuk inisialisasi awal (Idealnya butuh Role Admin, tapi untuk MVP kita buka dulu)
func (h *ModuleHandler) Create(c *fiber.Ctx) error {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		LevelOrder  int    `json:"level_order"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.CreateModule(req.Title, req.Description, req.LevelOrder); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Modul berhasil dibuat"})
}
