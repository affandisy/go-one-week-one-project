package handlers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *fiber.Ctx) error {
	var req services.RegisterRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request tidak valid",
		})
	}

	user, err := h.service.Register(req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registrasi berhasil",
		"data":    user,
	})
}

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req services.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request tidak valid",
		})
	}

	token, err := h.service.Login(req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   token,
	})
}

func (h *UserHandler) UploadAvatarImage(c *fiber.Ctx) error {
	userId, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "ID User tidak valid",
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

	filename := fmt.Sprintf("user_avatar_%d_%d%s", userId, time.Now().Unix(), ext)
	savePath := fmt.Sprintf("./uploads/users/%s", filename)

	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal menyimpan gambar",
		})
	}

	dbPath := fmt.Sprintf("/uploads/users/%s", filename)
	if err := h.service.UpdateUserAvatar(uint(userId), dbPath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":   "Gambar produk berhasil diunggah",
		"image_url": dbPath,
	})
}
