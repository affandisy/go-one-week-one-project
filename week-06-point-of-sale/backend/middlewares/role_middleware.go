package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil peran dari Locas (sudah diset oleh jwt_middleware sebelumnya)
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Peran tidak teridentifikasi."})
		}

		roleStr := userRole.(string)

		// Cek apakah peran user ada di dalam daftar peran yang diizinkan
		isAllowed := false
		for _, role := range allowedRoles {
			if roleStr == role {
				isAllowed = true
				break
			}
		}

		if !isAllowed {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Akses ditolak. Anda tidak memiliki izin (Hak Akses) untuk fitur ini.",
			})
		}

		return c.Next()
	}
}
