package handlers

import (
	"github.com/affandisy/pos-system/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service}
}

func (h *UserHandler) GetAll(c *fiber.Ctx) error {
	users, err := h.service.GetAllUsers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat data pengguna"})
	}
	return c.JSON(fiber.Map{"data": users})
}

func (h *UserHandler) Create(c *fiber.Ctx) error {
	var req services.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.CreateUser(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Akun pengguna berhasil dibuat"})
}

func (h *UserHandler) Update(c *fiber.Ctx) error {
	id := c.Params("id")
	var req services.UserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.UpdateUser(id, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Data pengguna berhasil diperbarui"})
}

func (h *UserHandler) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	// Keamanan tambahan: Cegah pengguna menghapus akunnya sendiri yang sedang dipakai login
	currentUserID := c.Locals("user_id").(string)
	if id == currentUserID {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Anda tidak dapat menghapus akun Anda sendiri yang sedang aktif"})
	}

	if err := h.service.DeleteUser(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "Akun pengguna berhasil dihapus"})
}
