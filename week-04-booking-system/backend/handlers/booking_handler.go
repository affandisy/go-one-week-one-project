package handlers

import (
	"github.com/affandisy/padel-booking-system/services"
	"github.com/gofiber/fiber/v2"
)

type BookingHandler struct {
	service services.BookingService
}

func NewBookingHandler(service services.BookingService) *BookingHandler {
	return &BookingHandler{service}
}

// GET /availability?court_id=xxx&date=2026-03-25
func (h *BookingHandler) GetAvailability(c *fiber.Ctx) error {
	courtID := c.Query("court_id")
	date := c.Query("date")

	if courtID == "" || date == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "parameter court_id dan date wajib diisi"})
	}

	slots, err := h.service.GetAvailability(courtID, date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"data": slots})
}

// POST /bookings
func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	var req struct {
		CourtID     string `json:"court_id"`
		BookingDate string `json:"booking_date"`
		StartTime   string `json:"start_time"`
		EndTime     string `json:"end_time"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	// Mengambil ID User dari token JWT yang sudah diverifikasi oleh Middleware
	userID := c.Locals("user_id").(string)

	booking, err := h.service.CreateBooking(userID, req.CourtID, req.BookingDate, req.StartTime, req.EndTime)
	if err != nil {
		// Menggunakan status 409 Conflict jika slot bentrok
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Slot berhasil dikunci. Selesaikan pembayaran dalam 10 menit.",
		"data": map[string]interface{}{
			"id":           booking.ID,
			"booking_code": booking.BookingCode,
			"status":       booking.Status,
			"lock_expiry":  booking.LockExpiry,
			"total_price":  booking.TotalPrice,
		},
	})
}
