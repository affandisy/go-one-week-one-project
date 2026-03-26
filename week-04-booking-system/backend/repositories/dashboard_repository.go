package repositories

import (
	"github.com/affandisy/padel-booking-system/models"
	"gorm.io/gorm"
)

type DashboardRepository interface {
	GetRevenue(period string) (float64, error)
	GetActiveBookingsCount() (int64, error)
	GetTodayOccupancy() (int64, error)
}

type dashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) DashboardRepository {
	return &dashboardRepository{db}
}

// Menghitung total pendapatan (hanya yang sudah PAID)
func (r *dashboardRepository) GetRevenue(period string) (float64, error) {
	var total float64
	query := r.db.Model(&models.Booking{}).Where("status = ?", "paid")

	switch period {
	case "daily":
		// Pendapatan Hari Ini
		query = query.Where("booking_date = CURRENT_DATE")
	case "monthly":
		// Pendapatan Bulan Ini
		query = query.Where("EXTRACT(MONTH FROM booking_date) = EXTRACT(MONTH FROM CURRENT_DATE) AND EXTRACT(YEAR FROM booking_date) = EXTRACT(YEAR FROM CURRENT_DATE)")
	case "yearly":
		// Pendapatan Tahun Ini
		query = query.Where("EXTRACT(YEAR FROM booking_date) = EXTRACT(YEAR FROM CURRENT_DATE)")
	}

	// Gunakan COALESCE agar mengembalikan 0 jika tidak ada data (mencegah error NULL)
	err := query.Select("COALESCE(SUM(total_price), 0)").Scan(&total).Error
	return total, err
}

// Menghitung booking yang sedang berjalan atau akan datang (Hari ini ke depan)
func (r *dashboardRepository) GetActiveBookingsCount() (int64, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).
		Where("status IN ? AND booking_date >= CURRENT_DATE", []string{"locked", "pending", "paid"}).
		Count(&count).Error
	return count, err
}

// Menghitung berapa banyak slot lapangan yang terpakai HARI INI
func (r *dashboardRepository) GetTodayOccupancy() (int64, error) {
	var count int64
	err := r.db.Model(&models.Booking{}).
		Where("status IN ? AND booking_date = CURRENT_DATE", []string{"locked", "pending", "paid"}).
		Count(&count).Error
	return count, err
}
