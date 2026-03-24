package services

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/affandisy/padel-booking-system/models"
	"github.com/affandisy/padel-booking-system/repositories"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type BookingService interface {
	GetAvailability(courtID, dateStr string) ([]map[string]interface{}, error)
	CreateBooking(userID, courtID, dateStr, start, end string) (*models.Booking, error)
	GetUserBookings(userID string) ([]models.Booking, error)
	GetAllBookings(page int, limit int, dateFilter string) ([]models.Booking, int64, error)
	CancelBooking(bookingID string, userID string, role string) error
	ProcessPayment(bookingID string, userID string) error
	GenerateReceiptPDF(bookingID string, userID string) (string, error)
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

func (s *bookingService) GetAllBookings(page int, limit int, dateFilter string) ([]models.Booking, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	} // Batas maksimal 100 per request
	offset := (page - 1) * limit

	return s.repo.GetAllBookings(limit, offset, dateFilter)
}

func (s *bookingService) CancelBooking(bookingID string, userID string, role string) error {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking tidak ditemukan")
	}

	// Validasi Hak Akses
	if role == "customer" && booking.UserID.String() != userID {
		return errors.New("anda tidak memiliki izin membatalkan booking ini")
	}

	// IMPROVISASI 2: Kebijakan Waktu Pembatalan (Hanya untuk Customer)
	if role == "customer" {
		// Gabungkan Tanggal dan Jam Main menjadi tipe time.Time
		playTimeStr := fmt.Sprintf("%s %s", booking.BookingDate.Format("2006-01-02"), booking.StartTime)
		playTime, _ := time.Parse("2006-01-02 15:04", playTimeStr)

		// Hitung selisih waktu sekarang dengan waktu bermain
		hoursUntilPlay := time.Until(playTime).Hours()

		if hoursUntilPlay < 24 {
			return errors.New("pembatalan ditolak. Anda hanya bisa membatalkan maksimal 24 jam sebelum jadwal bermain")
		}
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

func (s *bookingService) ProcessPayment(bookingID string, userID string) error {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil {
		return errors.New("booking tidak ditemukan")
	}

	// 1. Validasi Kepemilikan
	if booking.UserID.String() != userID {
		return errors.New("anda tidak memiliki akses untuk membayar booking ini")
	}

	// 2. Validasi Status & Waktu Kedaluwarsa
	if booking.Status == "paid" {
		return errors.New("booking ini sudah dibayar")
	}
	if booking.Status != "locked" {
		return errors.New("status booking tidak valid untuk pembayaran")
	}
	if time.Now().After(*booking.LockExpiry) {
		// Ubah status ke expired di database
		_ = s.repo.UpdateStatus(bookingID, "expired")
		return errors.New("waktu pembayaran telah habis (expired)")
	}

	// 3. Eksekusi Pembayaran (Simulasi sukses)
	// Di dunia nyata, di sinilah Anda memanggil API Midtrans/Xendit
	return s.repo.UpdateStatus(bookingID, "paid")
}

// Fitur 4: Cetak Bukti Booking (PDF)
func (s *bookingService) GenerateReceiptPDF(bookingID string, userID string) (string, error) {
	booking, err := s.repo.FindByID(bookingID)
	if err != nil || booking.UserID.String() != userID {
		return "", errors.New("data booking tidak ditemukan atau tidak valid")
	}

	if booking.Status != "paid" {
		return "", errors.New("bukti booking hanya bisa dicetak setelah pembayaran lunas")
	}

	// 1. Siapkan direktori penyimpanan
	storageDir := "./public/receipts"
	os.MkdirAll(storageDir, os.ModePerm)

	// 2. Menggambar Invoice PDF
	pdf := gofpdf.New("P", "mm", "A5", "") // Menggunakan ukuran A5 agar pas di layar HP
	pdf.AddPage()

	// Header Logo / Judul
	pdf.SetFont("Arial", "B", 18)
	pdf.CellFormat(130, 10, "PADEL BOOKING E-RECEIPT", "0", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "", 10)
	pdf.CellFormat(130, 6, "Bukti Pemesanan Lapangan Sah", "0", 1, "C", false, 0, "")
	pdf.Ln(8)

	// Info Booking
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 8, "Kode Booking", "0", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(90, 8, ": "+booking.BookingCode, "0", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 8, "Nama Pemesan", "0", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(90, 8, ": "+booking.User.FullName, "0", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 8, "Lapangan", "0", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(90, 8, ": "+booking.Court.Name, "0", 1, "L", false, 0, "")

	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(40, 8, "Tanggal Main", "0", 0, "L", false, 0, "")
	pdf.SetFont("Arial", "", 12)
	dateStr := booking.BookingDate.Format("02 Jan 2006")
	pdf.CellFormat(90, 8, fmt.Sprintf(": %s (%s - %s)", dateStr, booking.StartTime, booking.EndTime), "0", 1, "L", false, 0, "")

	pdf.Ln(10)

	// Kotak Total Harga & Status
	pdf.SetFillColor(230, 255, 230) // Warna hijau muda
	pdf.SetFont("Arial", "B", 14)
	pdf.CellFormat(130, 15, fmt.Sprintf("TOTAL DIBAYAR: Rp %.2f (LUNAS)", booking.TotalPrice), "1", 1, "C", true, 0, "")

	pdf.Ln(10)
	pdf.SetFont("Arial", "I", 9)
	pdf.CellFormat(130, 6, "*Tunjukkan e-receipt ini kepada petugas lapangan sebelum bermain.", "0", 1, "C", false, 0, "")

	// 3. Simpan dan Hasilkan URL
	fileName := fmt.Sprintf("receipt_%s.pdf", booking.BookingCode)
	filePath := fmt.Sprintf("%s/%s", storageDir, fileName)

	if err := pdf.OutputFileAndClose(filePath); err != nil {
		return "", errors.New("gagal membuat file PDF")
	}

	// URL yang bisa diakses dari Frontend
	pdfURL := fmt.Sprintf("http://localhost:3000/public/receipts/%s", fileName)
	return pdfURL, nil
}
