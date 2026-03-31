package handlers

import (
	"github.com/affandi/belajar-bahasa/services"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

// Struct untuk membaca body JSON
type AuthRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req AuthRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.Register(req.Username, req.Password); err != nil {
		// Menggunakan status 409 Conflict jika username sudah ada
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registrasi berhasil, silakan login",
	})
}

func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req AuthRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	token, user, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Login berhasil",
		"data": fiber.Map{
			"accessToken": token,
			"user": fiber.Map{
				"id":       user.ID,
				"username": user.Username,
			},
		},
	})
}
