package repositories

import (
	"github.com/affandisy/padel-booking-system/models"
	"gorm.io/gorm"
)

type BookingRepository interface {
	BeginTx() *gorm.DB
	FindConflictingBookings(tx *gorm.DB, courtID string, date string, start, end string) ([]models.Booking, error)
	CreateBooking(tx *gorm.DB, booking *models.Booking) error
	GetBookingsByDate(courtID string, date string) ([]models.Booking, error)
	GetPricingRules(courtID string, dayType string) ([]models.PricingRule, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db}
}

// Membuka koneksi transaksi khusus
func (r *bookingRepository) BeginTx() *gorm.DB {
	return r.db.Begin()
}

// Mencari jadwal yang bentrok (di jam yang sama, dan statusnya belum dibatalkan)
func (r *bookingRepository) FindConflictingBookings(tx *gorm.DB, courtID string, date string, start, end string) ([]models.Booking, error) {
	var bookings []models.Booking
	err := tx.Where("court_id = ? AND booking_date = ? AND status IN ? AND ((start_time < ? AND end_time > ?) OR (start_time >= ? AND start_time < ?))",
		courtID, date, []string{"locked", "pending", "paid"}, end, start, start, end).Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) CreateBooking(tx *gorm.DB, booking *models.Booking) error {
	return tx.Create(booking).Error
}

func (r *bookingRepository) GetBookingsByDate(courtID string, date string) ([]models.Booking, error) {
	var bookings []models.Booking
	err := r.db.Where("court_id = ? AND booking_date = ? AND status IN ?", courtID, date, []string{"locked", "pending", "paid"}).Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) GetPricingRules(courtID string, dayType string) ([]models.PricingRule, error) {
	var rules []models.PricingRule
	err := r.db.Where("court_id = ? AND day_type = ?", courtID, dayType).Find(&rules).Error
	return rules, err
}
