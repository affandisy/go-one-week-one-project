package handlers

import (
	"github.com/affandisy/school-system-sma/service"
	"github.com/gofiber/fiber/v2"
)

type MasterHandler struct {
	service service.MasterService
}

func NewMasterHandler(service service.MasterService) *MasterHandler {
	return &MasterHandler{service}
}

// ==== CLASS ====
func (h *MasterHandler) CreateClass(c *fiber.Ctx) error {
	var req service.CreateClassRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tidak valid"})
	}
	if err := h.service.CreateClass(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Kelas berhasil ditambahkan"})
}
func (h *MasterHandler) GetClasses(c *fiber.Ctx) error {
	data, _ := h.service.GetClasses()
	return c.JSON(fiber.Map{"data": data})
}

// ==== SUBJECT ====
func (h *MasterHandler) CreateSubject(c *fiber.Ctx) error {
	var req service.CreateSubjectRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tidak valid"})
	}
	if err := h.service.CreateSubject(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Mata pelajaran berhasil ditambahkan"})
}
func (h *MasterHandler) GetSubjects(c *fiber.Ctx) error {
	data, _ := h.service.GetSubjects()
	return c.JSON(fiber.Map{"data": data})
}

// ==== STUDENT ====
func (h *MasterHandler) CreateStudent(c *fiber.Ctx) error {
	var req service.CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Format tidak valid"})
	}
	if err := h.service.CreateStudent(req); err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"message": "Siswa berhasil didaftarkan"})
}
func (h *MasterHandler) GetStudents(c *fiber.Ctx) error {
	data, _ := h.service.GetStudents()
	return c.JSON(fiber.Map{"data": data})
}
