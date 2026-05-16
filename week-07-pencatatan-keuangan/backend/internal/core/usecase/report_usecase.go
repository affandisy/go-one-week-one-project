package usecase

import (
	"time"

	"github.com/affandisy/financial-app/internal/core/domain"
	"github.com/affandisy/financial-app/internal/core/port"
)

type reportUseCase struct {
	repo port.ReportRepository
}

func NewReportUseCase(repo port.ReportRepository) port.ReportUseCase {
	return &reportUseCase{repo}
}

func (uc *reportUseCase) GetMonthlySummary(walletID string, year int, month time.Month) (*domain.MonthlyReport, error) {
	// 1. Tentukan batas waktu
	startThisMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	startNextMonth := startThisMonth.AddDate(0, 1, 0)
	startLastMonth := startThisMonth.AddDate(0, -1, 0)

	// 2. Tarik data dari Port (Repository)
	thisMonthTrx, err := uc.repo.GetTransactionsForReport(walletID, startThisMonth, startNextMonth)
	if err != nil {
		return nil, err
	}
	lastMonthTrx, _ := uc.repo.GetTransactionsForReport(walletID, startLastMonth, startThisMonth)

	report := &domain.MonthlyReport{}
	categoryMap := make(map[string]float64)
	colorMap := make(map[string]string)

	// 3. Kalkulasi Bulan Ini
	for _, t := range thisMonthTrx {
		if t.Type == "income" {
			report.TotalIncome += t.Amount
		} else if t.Type == "expense" {
			report.TotalExpense += t.Amount
			categoryMap[t.Category.Name] += t.Amount

			color := t.Category.Color
			if color == "" {
				color = "#3B82F6"
			}
			colorMap[t.Category.Name] = color
		}
	}
	report.NetBalance = report.TotalIncome - report.TotalExpense

	// 4. Kalkulasi Bulan Lalu (FR-004)
	var lastMonthExpense float64
	for _, t := range lastMonthTrx {
		if t.Type == "expense" {
			lastMonthExpense += t.Amount
		}
	}

	if lastMonthExpense > 0 {
		report.VsLastMonthPercent = ((report.TotalExpense - lastMonthExpense) / lastMonthExpense) * 100
	} else if report.TotalExpense > 0 {
		report.VsLastMonthPercent = 100.0
	}

	// 5. Kalkulasi Persentase Kategori untuk Pie Chart
	for name, amount := range categoryMap {
		percentage := 0.0
		if report.TotalExpense > 0 {
			percentage = (amount / report.TotalExpense) * 100
		}
		report.ExpenseBreakdown = append(report.ExpenseBreakdown, domain.CategorySummary{
			CategoryName: name,
			Color:        colorMap[name],
			Amount:       amount,
			Percentage:   percentage,
		})
	}

	return report, nil
}
