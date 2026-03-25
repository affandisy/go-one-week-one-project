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

// CreateBooking godoc
// @Summary Mengunci Slot Lapangan (Create Booking)
// @Description Mem-booking slot lapangan yang kosong dan menguncinya selama 10 menit
// @Tags Bookings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body object true "Data Booking" SchemaExample({\"court_id\":\"uuid-lapangan\",\"booking_date\":\"2026-03-25\",\"start_time\":\"19:00\",\"end_time\":\"20:00\"})
// @Success 201 {object} map[string]interface{} "Slot berhasil dikunci"
// @Failure 400 {object} map[string]string "Format tidak valid"
// @Failure 401 {object} map[string]string "Unauthorized (Token tidak ada/salah)"
// @Failure 409 {object} map[string]string "Conflict (Slot sudah diambil orang lain)"
// @Router /bookings [post]
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

// GET /bookings/me
func (h *BookingHandler) GetMyBookings(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(string)

	bookings, err := h.service.GetUserBookings(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "gagal mengambil riwayat booking"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": bookings})
}

// GET /admin/bookings
func (h *BookingHandler) GetAllBookings(c *fiber.Ctx) error {
	// Ambil query parameter
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	dateFilter := c.Query("date", "") // Format: YYYY-MM-DD

	bookings, total, err := h.service.GetAllBookings(page, limit, dateFilter)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "gagal mengambil data"})
	}

	return c.Status(200).JSON(fiber.Map{
		"data": bookings,
		"meta": fiber.Map{
			"current_page": page,
			"per_page":     limit,
			"total_data":   total,
		},
	})
}

// PUT /bookings/:id/cancel
func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	bookingID := c.Params("id")
	userID := c.Locals("user_id").(string)
	role := c.Locals("role").(string)

	if err := h.service.CancelBooking(bookingID, userID, role); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Booking berhasil dibatalkan"})
}

// POST /bookings/:id/pay (Simulasi Pembayaran)
func (h *BookingHandler) PayBooking(c *fiber.Ctx) error {
	bookingID := c.Params("id")
	userID := c.Locals("user_id").(string)

	if err := h.service.ProcessPayment(bookingID, userID); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Pembayaran berhasil diverifikasi. Status menjadi PAID."})
}

// GET /bookings/:id/receipt (Generate & Dapatkan URL PDF)
func (h *BookingHandler) DownloadReceipt(c *fiber.Ctx) error {
	bookingID := c.Params("id")
	userID := c.Locals("user_id").(string)

	pdfURL, err := h.service.GenerateReceiptPDF(bookingID, userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Bukti pembayaran berhasil dibuat",
		"pdf_url": pdfURL,
	})
}
