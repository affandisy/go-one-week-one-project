package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type ReportService interface {
	GenerateMonthlyTransactionCSV(month, year int) ([]byte, error)
}

type reportService struct {
	txRepo repositories.TransactionRepository
}

func NewReportService(txRepo repositories.TransactionRepository) ReportService {
	return &reportService{txRepo: txRepo}
}

func (s *reportService) GenerateMonthlyTransactionCSV(month, year int) ([]byte, error) {

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	transactions, err := s.txRepo.FindTransactionsByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	var buffer bytes.Buffer
	writer := csv.NewWriter(&buffer)

	header := []string{
		"ID Transaksi",
		"Tanggal",
		"Tipe",
		"No Referensi",
		"Status",
		"Nama Partner (Supplier/Customer)",
		"Jumlah Macam Barang",
		"Catatan",
	}

	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for _, tx := range transactions {
		totalItems := fmt.Sprintf("%d item(s)", len(tx.Items))

		partnerName := "-"
		if tx.Partner.Name != "" {
			partnerName = tx.Partner.Name
		}

		row := []string{
			fmt.Sprintf("%d", tx.ID),
			tx.TransactionDate.Format("2006-01-02 15:04"),
			tx.Type,
			tx.ReferenceNo,
			tx.Status,
			partnerName,
			totalItems,
			tx.Notes,
		}

		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
