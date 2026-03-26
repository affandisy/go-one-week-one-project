package handlers

import (
	"github.com/affandisy/padel-booking-system/services"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	service services.UserService
}

func NewUserHandler(service services.UserService) *UserHandler {
	return &UserHandler{service}
}

// UpdateProfile godoc
// @Summary Ubah Data Profil
// @Description Mengubah nama lengkap atau nomor WhatsApp pengguna yang sedang login
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.UpdateProfileRequest true "Data Profil Baru"
// @Success 200 {object} map[string]string "Profil berhasil diperbarui"
// @Failure 400 {object} map[string]string "Format tidak valid"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *fiber.Ctx) error {
	var req services.UpdateProfileRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	// Ambil ID User dari Token JWT
	userID := c.Locals("user_id").(string)

	if err := h.service.UpdateProfile(userID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Profil berhasil diperbarui"})
}

// UpdatePassword godoc
// @Summary Ubah Password
// @Description Mengubah password dengan memvalidasi password lama terlebih dahulu
// @Tags Users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body services.UpdatePasswordRequest true "Kredensial Password"
// @Success 200 {object} map[string]string "Password berhasil diubah"
// @Failure 400 {object} map[string]string "Password lama salah atau format tidak valid"
// @Router /users/password [put]
func (h *UserHandler) UpdatePassword(c *fiber.Ctx) error {
	var req services.UpdatePasswordRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	userID := c.Locals("user_id").(string)

	if err := h.service.UpdatePassword(userID, req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Password berhasil diubah dengan aman"})
}
