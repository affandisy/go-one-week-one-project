package services

import (
	"time"

	"github.com/affandisy/financial-app/repositories"
)

type CategorySummary struct {
	CategoryName string  `json:"category_name"`
	Color        string  `json:"color"`
	Amount       float64 `json:"amount"`
	Percentage   float64 `json:"percentage"`
}

type MonthlyReport struct {
	TotalIncome        float64           `json:"total_income"`
	TotalExpense       float64           `json:"total_expense"`
	NetBalance         float64           `json:"net_balance"`
	VsLastMonthPercent float64           `json:"vs_last_month_percent"` // FR-004: Persentase perbandingan
	ExpenseBreakdown   []CategorySummary `json:"expense_breakdown"`     // FR-004: Untuk Pie Chart
}

type ReportService interface {
	GetMonthlyReport(year int, month time.Month) (*MonthlyReport, error)
}

type reportService struct {
	reportRepo repositories.ReportRepository
}

func NewReportService(reportRepo repositories.ReportRepository) ReportService {
	return &reportService{reportRepo}
}

func (s *reportService) GetMonthlyReport(year int, month time.Month) (*MonthlyReport, error) {
	// 1. Tentukan batas tanggal menggunakan Time di Go
	startThisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	startNextMonth := startThisMonth.AddDate(0, 1, 0) // Tanggal 1 bulan depan

	startLastMonth := startThisMonth.AddDate(0, -1, 0)
	// Batas akhir bulan lalu adalah awal bulan ini (startThisMonth)

	// 2. Tarik data dari Repositori Laporan
	thisMonthTrx, err := s.reportRepo.GetMonthlyData(startThisMonth, startNextMonth)
	if err != nil {
		return nil, err
	}
	lastMonthTrx, _ := s.reportRepo.GetMonthlyData(startLastMonth, startThisMonth)

	report := &MonthlyReport{}
	categoryMap := make(map[string]float64)
	colorMap := make(map[string]string)

	// 3. Kalkulasi Bulan Ini
	for _, t := range thisMonthTrx {
		if t.Type == "income" {
			report.TotalIncome += t.Amount
		} else if t.Type == "expense" {
			report.TotalExpense += t.Amount
			categoryMap[t.Category.Name] += t.Amount

			// Ambil warna untuk pie chart, fallback jika tidak ada warna
			color := t.Category.Color
			if color == "" {
				color = "#3B82F6"
			}
			colorMap[t.Category.Name] = color
		}
	}
	report.NetBalance = report.TotalIncome - report.TotalExpense

	// 4. Kalkulasi Bulan Lalu untuk Perbandingan (FR-004)
	var lastMonthExpense float64
	for _, t := range lastMonthTrx {
		if t.Type == "expense" {
			lastMonthExpense += t.Amount
		}
	}

	// Formula Perbandingan: ((Bulan Ini - Bulan Lalu) / Bulan Lalu) * 100
	if lastMonthExpense > 0 {
		report.VsLastMonthPercent = ((report.TotalExpense - lastMonthExpense) / lastMonthExpense) * 100
	} else if report.TotalExpense > 0 {
		// Jika bulan lalu 0 pengeluaran, kenaikannya dianggap 100%
		report.VsLastMonthPercent = 100.0
	}

	// 5. Hitung Persentase per Kategori untuk Tabel & Pie Chart (FR-004)
	for name, amount := range categoryMap {
		percentage := 0.0
		if report.TotalExpense > 0 {
			percentage = (amount / report.TotalExpense) * 100
		}
		report.ExpenseBreakdown = append(report.ExpenseBreakdown, CategorySummary{
			CategoryName: name,
			Color:        colorMap[name],
			Amount:       amount,
			Percentage:   percentage,
		})
	}

	return report, nil
}
