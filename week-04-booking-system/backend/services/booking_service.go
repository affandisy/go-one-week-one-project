package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/affandisy/padel-booking-system/models"
	"github.com/affandisy/padel-booking-system/repositories"
	"github.com/google/uuid"
)

type BookingService interface {
	GetAvailability(courtID, dateStr string) ([]map[string]interface{}, error)
	CreateBooking(userID, courtID, dateStr, start, end string) (*models.Booking, error)
	GetUserBookings(userID string) ([]models.Booking, error)
	GetAllBookings() ([]models.Booking, error)
	CancelBooking(bookingID string, userID string, role string) error
	RunAutoExpireJob()
}

type bookingService struct {
	repo repositories.BookingRepository
}

func NewBookingService(repo repositories.BookingRepository) BookingService {
	return &bookingService{repo}
}

// Helper: Menghitung Harga Dinamis
func (s *bookingService) calculatePrice(courtID, dateStr, start string) (float64, error) {
	date, _ := time.Parse("2006-01-02", dateStr)

	// Tentukan apakah ini hari biasa atau akhir pekan
	dayType := "weekday"
	if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		dayType = "weekend"
	}

	rules, _ := s.repo.GetPricingRules(courtID, dayType)

	// Nilai default jika admin belum mengatur harga (misal: Rp 150.000)
	basePrice := 150000.0
	multiplier := 1.0

	// Cari kecocokan aturan harga dengan jam mulai (start time)
	for _, rule := range rules {
		if start >= rule.StartTime && start < rule.EndTime {
			basePrice = rule.BasePrice
			multiplier = rule.Multiplier
			break
		}
	}

	return basePrice * multiplier, nil
}

// Fitur 1: Melihat Jadwal Kosong
func (s *bookingService) GetAvailability(courtID, dateStr string) ([]map[string]interface{}, error) {
	bookings, _ := s.repo.GetBookingsByDate(courtID, dateStr)

	var slots []map[string]interface{}

	// Asumsi lapangan buka jam 06:00 sampai 22:00 (slot per 1 jam)
	for hour := 6; hour < 22; hour++ {
		start := fmt.Sprintf("%02d:00", hour)
		end := fmt.Sprintf("%02d:00", hour+1)

		status := "available"
		for _, b := range bookings {
			if b.StartTime == start {
				status = b.Status // slot akan ditandai locked/pending/paid
				break
			}
		}

		price, _ := s.calculatePrice(courtID, dateStr, start)

		slots = append(slots, map[string]interface{}{
			"start":  start,
			"end":    end,
			"price":  price,
			"status": status,
		})
	}

	return slots, nil
}

// Fitur 2: Mengunci Jadwal (Create Booking)
func (s *bookingService) CreateBooking(userID, courtID, dateStr, start, end string) (*models.Booking, error) {
	// Memulai Transaksi Database!
	tx := s.repo.BeginTx()
	// WAJIB: Jika ada kegagalan di tengah jalan, batalkan semuanya
	defer tx.Rollback()

	// 1. Cek Ketersediaan (Pencegahan Double Booking)
	existing, err := s.repo.FindConflictingBookings(tx, courtID, dateStr, start, end)
	if err != nil || len(existing) > 0 {
		return nil, errors.New("maaf, slot waktu tersebut sudah tidak tersedia (baru saja diambil orang lain)")
	}

	// 2. Kalkulasi Harga Final
	price, _ := s.calculatePrice(courtID, dateStr, start)

	// 3. Konversi Data & Generate Kode Unik
	uID, _ := uuid.Parse(userID)
	cID, _ := uuid.Parse(courtID)
	bookingDate, _ := time.Parse("2006-01-02", dateStr)

	lockExpiry := time.Now().Add(10 * time.Minute) // Kunci selama 10 menit untuk pembayaran
	bookingCode := fmt.Sprintf("PDL-%s-%04d", time.Now().Format("060102"), time.Now().Nanosecond()%10000)

	booking := &models.Booking{
		BookingCode: bookingCode,
		UserID:      uID,
		CourtID:     cID,
		BookingDate: bookingDate,
		StartTime:   start,
		EndTime:     end,
		Status:      "locked", // Slot ini resmi dikunci!
		LockExpiry:  &lockExpiry,
		TotalPrice:  price,
	}

	// 4. Simpan ke Database melalui Transaksi
	if err := s.repo.CreateBooking(tx, booking); err != nil {
		return nil, errors.New("sistem gagal memproses booking Anda")
	}

	// 5. Berhasil! Simpan perubahan secara permanen
	tx.Commit()

	return booking, nil
}

func (s *bookingService) GetUserBookings(userID string) ([]models.Booking, error) {
	return s.repo.GetBookingsByUser(userID)
}

func (s *bookingService) GetAllBookings() ([]models.Booking, error) {
	return s.repo.GetAllBookings()
}

func (s *bookingService) CancelBooking(bookingID string, userID string, role string) error {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking tidak ditemukan")
	}

	// Validasi Otorisasi: Hanya Admin/Owner ATAU pemilik booking yang boleh membatalkan
	if role == "customer" && booking.UserID.String() != userID {
		return errors.New("anda tidak memiliki izin membatalkan booking ini")
	}

	// Hanya booking yang belum dibayar atau masih pending yang bisa dibatalkan dengan mudah
	if booking.Status == "paid" {
		return errors.New("booking yang sudah dibayar harus melalui proses refund ke Admin")
	}

	return s.repo.UpdateStatus(bookingID, "cancelled")
}

func (s *bookingService) RunAutoExpireJob() {
	ticker := time.NewTicker(1 * time.Minute) // Berjalan setiap 1 Menit
	go func() {
		for {
			<-ticker.C // Tunggu detak ticker

			// 1. Cari booking yang waktunya habis
			expiredBookings, err := s.repo.GetExpiredLockedBookings()
			if err != nil || len(expiredBookings) == 0 {
				continue // Aman, tidak ada yang expired
			}

			// 2. Ubah statusnya menjadi expired agar slotnya kembali tersedia
			for _, b := range expiredBookings {
				_ = s.repo.UpdateStatus(b.ID.String(), "expired")
				fmt.Printf("[AUTO-EXPIRE JOB] Booking %s kedaluwarsa dan slot dirilis.\n", b.BookingCode)
			}
		}
	}()
}
