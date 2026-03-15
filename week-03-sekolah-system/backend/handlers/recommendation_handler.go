package handlers

import (
	"github.com/affandisy/school-system-sma/service"
	"github.com/gofiber/fiber/v2"
)

type RecommendationHandler struct {
	service service.RecommendationService
}

func NewRecommendationHandler(service service.RecommendationService) *RecommendationHandler {
	return &RecommendationHandler{service: service}
}

// Struct untuk menangkap body JSON dari Frontend
type CalculateRankingRequest struct {
	ClassID        string `json:"class_id"`
	Semester       int    `json:"semester"`
	AcademicYearID string `json:"academic_year_id"`
}

func (h *RecommendationHandler) CalculateRanking(c *fiber.Ctx) error {
	var req CalculateRankingRequest

	// 1. Parsing Body JSON
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Format request tidak valid",
		})
	}

	// 2. Validasi input dasar
	if req.ClassID == "" || req.AcademicYearID == "" || req.Semester == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Class ID, Academic Year ID, dan Semester wajib diisi",
		})
	}

	// 3. Panggil Otak Aplikasi (Service)
	rankedStudents, err := h.service.CalculateRanking(req.ClassID, req.Semester, req.AcademicYearID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// 4. Kembalikan Response Sukses
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Kalkulasi peringkat siswa berhasil diselesaikan",
		"data":    rankedStudents,
	})
}
