package handlers

import (
	"github.com/affandisy/school-system-sma/service"
	"github.com/gofiber/fiber/v2"
)

type AnnouncementHandler struct {
	service service.AnnouncementService
}

func NewAnnouncementHandler(service service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{service: service}
}

func (h *AnnouncementHandler) CreateAnnouncement(c *fiber.Ctx) error {
	var req struct {
		Title       string   `json:"title"`
		Content     string   `json:"content"`
		TargetRoles []string `json:"target_roles"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tidak valid"})
	}

	// Ambil ID pembuat dari JWT Token yang sudah dipasang oleh Middleware
	creatorID := c.Locals("user_id").(string)

	err := h.service.CreateAnnouncement(req.Title, req.Content, req.TargetRoles, creatorID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Pengumuman berhasil disiarkan"})
}

// Handler untuk memicu laporan asinkron
func (h *AnnouncementHandler) RequestAcademicReport(c *fiber.Ctx) error {
	academicYearID := c.Query("academic_year_id")
	requestedByID := c.Locals("user_id").(string)

	err := h.service.TriggerAcademicReportGeneration(academicYearID, requestedByID)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	// Response Instan!
	return c.Status(202).JSON(fiber.Map{
		"message": "Permintaan laporan sedang diproses di latar belakang. Anda akan diberitahu saat laporan siap.",
	})
}
