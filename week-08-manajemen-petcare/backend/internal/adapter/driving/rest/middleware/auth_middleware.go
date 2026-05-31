package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey []byte

// Fungsi untuk menyetel kunci dari main.go
func SetJWTSecret(secret string) {
	jwtSecretKey = []byte(secret)
}

// Protect mendeteksi dan memvalidasi JWT token
func Protect() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token otorisasi tidak ditemukan"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return jwtSecretKey, nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token tidak valid atau sudah kedaluwarsa"})
		}

		claims := token.Claims.(jwt.MapClaims)

		// Simpan data pengguna ke context Fiber agar bisa dipakai di Handler
		c.Locals("userID", claims["sub"])
		c.Locals("userRole", claims["role"])

		return c.Next()
	}
}

// RequireRole membatasi akses berdasarkan peran
func RequireRole(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("userRole").(string)

		for _, role := range allowedRoles {
			if role == userRole {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Anda tidak memiliki akses ke sumber daya ini"})
	}
}
