package handlers

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

func (h *UploadHandler) UploadFile(c *fiber.Ctx) error {
	// Ambil file dari form-data dengan key "file"
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Gagal membaca file upload"})
	}

	// Validasi ekstensi file (hanya izinkan gambar dan audio)
	ext := filepath.Ext(file.Filename)
	allowedExts := map[string]bool{
		".png": true, ".jpg": true, ".jpeg": true, ".svg": true,
		".mp3": true, ".wav": true,
	}
	if !allowedExts[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format file tidak diizinkan"})
	}

	// Buat folder public/uploads jika belum ada
	uploadDir := "./public/uploads"
	os.MkdirAll(uploadDir, os.ModePerm)

	// Buat nama file unik (timestamp + ekstensi)
	fileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	// Simpan file ke server
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan file"})
	}

	// Kembalikan URL yang bisa diakses oleh Frontend
	fileURL := fmt.Sprintf("http://localhost:3000/public/uploads/%s", fileName)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":  "File berhasil diunggah",
		"file_url": fileURL,
	})
}
