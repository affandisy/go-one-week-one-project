package services

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/config"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
)

type ReportService interface {
	GenerateMonthlyTransactionCSV(month, year int) ([]byte, error)
	TriggerMonthlyTransactionExport(month, year int, email string) error
	GetMovementAnalytics(month, year int) ([]repositories.ProductMovementResult, error)
}

type reportService struct {
	txRepo repositories.TransactionRepository
}

func NewReportService(txRepo repositories.TransactionRepository) ReportService {
	return &reportService{txRepo: txRepo}
}

func (s *reportService) TriggerMonthlyTransactionExport(month, year int, email string) error {

	go s.processExportBackground(month, year, email)

	return nil
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

func (s *reportService) processExportBackground(month, year int, email string) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	os.MkdirAll("./exports", os.ModePerm)
	fileName := fmt.Sprintf("laporan_%d_%d_%d.csv", month, year, time.Now().Unix())
	filePath := filepath.Join("./exports", fileName)

	file, err := os.Create(filePath)
	if err != nil {
		config.Log.Error().Err(err).Msg("Gagal membuat file export CSV")
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	header := []string{"ID Transaksi", "Tanggal", "Tipe", "Status", "Nama Partner", "Jumlah Macam Barang", "Catatan"}
	writer.Write(header)

	limit := 10000
	offset := 0

	config.Log.Info().Msgf("Memulai proses export CSV latar belakang untuk periode %d-%d", month, year)

	for {
		transactions, err := s.txRepo.FindTransactionsPaginatedByDate(startDate, endDate, limit, offset)
		if err != nil {
			config.Log.Error().Err(err).Msg("Gagal mengambil data untuk export")
			return
		}

		if len(transactions) == 0 {
			break
		}

		for _, tx := range transactions {
			partnerName := "-"
			if tx.Partner.Name != "" {
				partnerName = tx.Partner.Name
			}

			row := []string{
				fmt.Sprintf("%d", tx.ID),
				tx.TransactionDate.Format("2006-01-02 15:04"),
				tx.Type,
				tx.Status,
				partnerName,
				fmt.Sprintf("%d items", len(tx.Items)),
				tx.Notes,
			}
			writer.Write(row)
		}

		writer.Flush()
		offset += limit
	}

	config.Log.Info().Msgf("Export CSV sukses! Data tersimpan di %s. Siap dikirim ke email %s", filePath, email)
}

func (s *reportService) GetMovementAnalytics(month, year int) ([]repositories.ProductMovementResult, error) {
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, 0).Add(-time.Nanosecond)

	results, err := s.txRepo.AnalyzeProductMovement(startDate, endDate)
	if err != nil {
		return nil, err
	}

	return results, err
}
