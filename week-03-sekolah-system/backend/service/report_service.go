package service

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/affandisy/school-system-sma/models"
	"github.com/affandisy/school-system-sma/repository"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type ReportService interface {
	TriggerReportGeneration(studentID string, semester int, academicYearID string) error
	GetReports() ([]models.ReportCard, error)
}

type reportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) ReportService {
	return &reportService{repo}
}

func (s *reportService) GetReports() ([]models.ReportCard, error) {
	return s.repo.GetReportCards()
}

// Memicu proses asinkron (Goroutine)
func (s *reportService) TriggerReportGeneration(studentID string, semester int, academicYearID string) error {
	if studentID == "" || academicYearID == "" || semester == 0 {
		return errors.New("Data tidak lengkap untuk generate rapor")
	}

	// JALANKAN DI LATAR BELAKANG! (Tidak membuat API menunggu)
	go s.generatePDFWorker(studentID, semester, academicYearID)

	return nil
}

// Pekerja Pembuat PDF
func (s *reportService) generatePDFWorker(studentID string, semester int, academicYearID string) {
	log.Printf("📄 [WORKER] Memulai generate PDF Rapor untuk Siswa: %s", studentID)

	// 1. Ambil Data Siswa & Nilai
	student, err := s.repo.GetStudentByID(studentID)
	if err != nil {
		log.Println("❌ [WORKER] Siswa tidak ditemukan:", err)
		return
	}

	grades, _ := s.repo.GetStudentGrades(studentID)

	// 2. Hitung Nilai Akhir (Golang Cerdas)
	var totalScore float64 = 0
	for _, g := range grades {
		totalScore += g.Score
	}
	finalScore := float64(0)
	if len(grades) > 0 {
		finalScore = totalScore / float64(len(grades))
	}

	// 3. Buat Folder Penyimpanan jika belum ada
	storageDir := "./public/reports"
	os.MkdirAll(storageDir, os.ModePerm)

	// 4. Proses Menggambar PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Header
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(190, 10, "RAPOR AKADEMIK SISWA", "0", 1, "C", false, 0, "")
	pdf.Ln(10)

	// Identitas Siswa
	pdf.SetFont("Arial", "", 12)
	pdf.CellFormat(40, 8, "Nama Siswa", "0", 0, "L", false, 0, "")
	pdf.CellFormat(150, 8, ": "+student.User.FullName, "0", 1, "L", false, 0, "")
	pdf.CellFormat(40, 8, "NIS", "0", 0, "L", false, 0, "")
	pdf.CellFormat(150, 8, ": "+student.NIS, "0", 1, "L", false, 0, "")
	pdf.CellFormat(40, 8, "Semester", "0", 0, "L", false, 0, "")
	pdf.CellFormat(150, 8, fmt.Sprintf(": %d", semester), "0", 1, "L", false, 0, "")
	pdf.Ln(10)

	// Tabel Nilai
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(100, 10, "Mata Pelajaran", "1", 0, "C", false, 0, "")
	pdf.CellFormat(40, 10, "Nilai", "1", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	for _, g := range grades {
		mapelName := fmt.Sprintf("Mapel %s", g.SubjectID.String())
		pdf.CellFormat(100, 10, mapelName, "1", 0, "L", false, 0, "")
		pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", g.Score), "1", 1, "C", false, 0, "")
	}

	pdf.Ln(10)
	pdf.SetFont("Arial", "B", 12)
	pdf.CellFormat(100, 10, "NILAI RATA-RATA AKHIR", "1", 0, "R", false, 0, "")
	pdf.CellFormat(40, 10, fmt.Sprintf("%.2f", finalScore), "1", 1, "C", false, 0, "")

	// 5. Simpan File PDF
	fileName := fmt.Sprintf("rapor_%s_sem%d.pdf", student.NIS, semester)
	filePath := fmt.Sprintf("%s/%s", storageDir, fileName)
	pdfURL := fmt.Sprintf("http://localhost:3000/public/reports/%s", fileName)

	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		log.Println("❌ [WORKER] Gagal membuat file PDF:", err)
		return
	}

	// 6. Simpan URL PDF ke Database
	studUUID, _ := uuid.Parse(studentID)
	ayUUID, _ := uuid.Parse(academicYearID)

	reportCard := &models.ReportCard{
		StudentID:      studUUID,
		Semester:       semester,
		AcademicYearID: ayUUID,
		FinalScore:     finalScore,
		PdfURL:         pdfURL,
	}

	s.repo.UpsertReportCard(reportCard)
	log.Printf("✅ [WORKER] PDF Rapor %s berhasil dibuat dan URL disimpan ke DB!", fileName)
}
