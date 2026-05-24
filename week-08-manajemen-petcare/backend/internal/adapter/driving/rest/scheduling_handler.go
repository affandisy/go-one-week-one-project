package rest

import (
	"time"

	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type SchedulingHandler struct {
	useCase port.SchedulingUseCase
}

func NewSchedulingHandler(useCase port.SchedulingUseCase) *SchedulingHandler {
	return &SchedulingHandler{useCase}
}

type BookAppointmentReq struct {
	PetID     string `json:"pet_id"`
	ServiceID string `json:"service_id"`
	StartAt   string `json:"start_at"` // Format ISO8601
	EndAt     string `json:"end_at"`   // Format ISO8601
}

func (h *SchedulingHandler) CreateAppointment(c *fiber.Ctx) error {
	var req BookAppointmentReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Payload tidak valid"})
	}

	start, errStart := time.Parse(time.RFC3339, req.StartAt)
	end, errEnd := time.Parse(time.RFC3339, req.EndAt)
	if errStart != nil || errEnd != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format tanggal tidak valid, gunakan RFC3339"})
	}

	app, err := h.useCase.BookAppointment(req.PetID, req.ServiceID, start, end)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"data": app})
}
