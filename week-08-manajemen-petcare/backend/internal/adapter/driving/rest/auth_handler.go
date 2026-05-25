package rest

import (
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	useCase port.AuthUseCase
}

func NewAuthHandler(useCase port.AuthUseCase) *AuthHandler {
	return &AuthHandler{useCase}
}

type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req AuthRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	token, err := h.useCase.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": fiber.Map{"token": token}, "message": "Login berhasil"})
}

// Opsional untuk MVP: Endpoint registrasi awal staf
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		AuthRequest
		Role string `json:"role"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	err := h.useCase.RegisterUser(req.Username, req.Password, req.Role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Pengguna berhasil didaftarkan"})
}
