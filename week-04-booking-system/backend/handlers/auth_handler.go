package handlers

import (
	"github.com/affandisy/padel-booking-system/services"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{service}
}

func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Whatsapp string `json:"whatsapp"`
		FullName string `json:"full_name"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.Register(req.Email, req.Whatsapp, req.FullName, req.Password); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(201).JSON(fiber.Map{"message": "Registrasi berhasil, silakan periksa WhatsApp Anda untuk OTP"})
}

func (h *AuthHandler) VerifyOTP(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		OTP   string `json:"otp"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.service.VerifyOTP(req.Email, req.OTP); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{"message": "Akun berhasil diaktifkan, silakan login"})
}

// Login godoc
// @Summary Login Pelanggan / Admin
// @Description Mengautentikasi user dan mengembalikan token JWT
// @Tags Authentication
// @Accept json
// @Produce json
// @Param request body object true "Kredensial Login" SchemaExample({\"email\":\"budi@email.com\",\"password\":\"rahasia123\"})
// @Success 200 {object} map[string]interface{} "Login berhasil"
// @Failure 400 {object} map[string]string "Format tidak valid"
// @Failure 401 {object} map[string]string "Kredensial salah"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	token, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(401).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   token,
	})
}
