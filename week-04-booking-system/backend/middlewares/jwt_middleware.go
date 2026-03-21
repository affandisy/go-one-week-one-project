package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// JWTProtected memvalidasi token JWT dari header Authorization
func JWTProtected(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// 1. Ambil header Authorization
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak, token tidak ditemukan",
			})
		}

		// 2. Pisahkan awalan "Bearer " dari token aslinya
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// 3. Parsing dan validasi token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Pastikan algoritma enkripsinya sesuai (HMAC)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token tidak valid atau sudah kedaluwarsa",
			})
		}

		// 4. Ekstrak isi token (Claims)
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Gagal membaca data dari token",
			})
		}

		// 5. Simpan Data ke Locals (Bisa diakses oleh Handler selanjutnya)
		c.Locals("user_id", claims["user_id"])
		c.Locals("role", claims["role"])

		// 6. Lanjut ke Handler berikutnya
		return c.Next()
	}
}

// RequireRoles membatasi akses berdasarkan role (admin, customer, owner)
func RequireRoles(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Ambil role dari JWT Locals
		userRole := c.Locals("role")
		if userRole == nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Akses ditolak, identitas role tidak ditemukan",
			})
		}

		roleStr := userRole.(string)

		// Cocokkan role user dengan role yang diizinkan
		for _, allowedRole := range allowedRoles {
			if roleStr == allowedRole {
				return c.Next() // Role cocok, izinkan lewat
			}
		}

		// Jika loop selesai dan tidak ada yang cocok, tolak akses
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "Akses ditolak, Anda tidak memiliki izin (Forbidden)",
		})
	}
}
