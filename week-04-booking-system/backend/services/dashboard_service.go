package services

import "github.com/affandisy/padel-booking-system/repositories"

type DashboardStatsResponse struct {
	Revenue struct {
		Daily   float64 `json:"daily"`
		Monthly float64 `json:"monthly"`
		Yearly  float64 `json:"yearly"`
	} `json:"revenue"`
	ActiveBookings int64 `json:"active_bookings"`
	TodayOccupancy int64 `json:"today_occupancy_slots"`
}

type DashboardService interface {
	GetDashboardStats() (*DashboardStatsResponse, error)
}

type dashboardService struct {
	repo repositories.DashboardRepository
}

func NewDashboardService(repo repositories.DashboardRepository) DashboardService {
	return &dashboardService{repo}
}

func (s *dashboardService) GetDashboardStats() (*DashboardStatsResponse, error) {
	var stats DashboardStatsResponse

	// 1. Ambil Data Pendapatan (Revenue)
	daily, _ := s.repo.GetRevenue("daily")
	monthly, _ := s.repo.GetRevenue("monthly")
	yearly, _ := s.repo.GetRevenue("yearly")

	stats.Revenue.Daily = daily
	stats.Revenue.Monthly = monthly
	stats.Revenue.Yearly = yearly

	// 2. Ambil Data Booking Aktif
	active, _ := s.repo.GetActiveBookingsCount()
	stats.ActiveBookings = active

	// 3. Ambil Data Okupansi Lapangan Hari Ini
	occupancy, _ := s.repo.GetTodayOccupancy()
	stats.TodayOccupancy = occupancy

	return &stats, nil
}
