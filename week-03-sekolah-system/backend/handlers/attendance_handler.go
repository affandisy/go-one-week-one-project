package handlers

import (
	"github.com/affandisy/school-system-sma/service"
	"github.com/gofiber/fiber/v2"
)

type AttendanceHandler struct {
	service service.AttendanceService
}

func NewAttendanceHandler(service service.AttendanceService) *AttendanceHandler {
	return &AttendanceHandler{service}
}

func (h *AttendanceHandler) RecordBatchAttendance(c *fiber.Ctx) error {
	var req service.BatchAttendanceRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	// Ambil ID Guru/Piket yang sedang login dari token JWT
	recordedByID := c.Locals("user_id").(string)

	if err := h.service.RecordStudentBatchAttendance(req, recordedByID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Absensi berhasil disimpan"})
}

func (h *AttendanceHandler) GetClassStudents(c *fiber.Ctx) error {
	classID := c.Params("classId")
	students, err := h.service.GetStudentsForAttendance(classID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memuat daftar siswa"})
	}
	return c.JSON(fiber.Map{"data": students})
}
