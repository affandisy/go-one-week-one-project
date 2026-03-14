package service

import (
	"errors"

	"github.com/affandisy/school-system-sma/models"
	"github.com/affandisy/school-system-sma/repository"
	"github.com/google/uuid"
)

type PayrollService interface {
	CalculatePayslip(teacherID string, month, year int) (*models.SalaryPayslip, error)
}

type payrollService struct {
	repo repository.PayrollRepository
}

func NewPayrollService(repo repository.PayrollRepository) PayrollService {
	return &payrollService{repo}
}

func (s *payrollService) CalculatePayslip(teacherID string, month, year int) (*models.SalaryPayslip, error) {
	// 1. Cek apakah slip gaji bulan ini sudah pernah dihitung?
	if s.repo.CheckPayslipExists(teacherID, month, year) {
		return nil, errors.New("Slip gaji untuk bulan dan tahun ini sudah di-generate sebelumnya")
	}

	// 2. Ambil data konfigurasi guru dari Repository
	config, err := s.repo.GetTeacherConfig(teacherID)
	if err != nil {
		return nil, errors.New("Konfigurasi gaji guru tidak ditemukan. Harap atur gaji pokok terlebih dahulu.")
	}

	// 3. LOGIKA AKUNTANSI (Golang Cerdas)
	var totalEarnings float64 = config.BaseSalary
	var totalDeductions float64 = 0

	// Kita akan merekam jejak perhitungan untuk dimasukkan ke JSONB "Details"
	details := make(map[string]interface{})
	details["Gaji Pokok"] = config.BaseSalary

	// Iterasi isi JSON Config yang dinamis dari database
	for itemName, itemValue := range config.Config {
		// Konversi nilai dari JSON interface{} ke float64
		valFloat, ok := itemValue.(float64)
		if !ok {
			continue // Lewati jika format tidak valid
		}

		// Anggap saja aturan bisnisnya: Jika namanya mengandung kata "Potongan", maka itu Deduction
		// (Di sistem yang lebih kompleks, Anda bisa me-lookup ke tabel SalaryComponent)
		if isDeduction(itemName) {
			totalDeductions += valFloat
			details[itemName] = -valFloat // Catat sebagai minus di slip
		} else {
			totalEarnings += valFloat
			details[itemName] = valFloat // Catat sebagai plus di slip
		}
	}

	netSalary := totalEarnings - totalDeductions

	// 4. Siapkan Objek Slip Gaji
	tUUID, _ := uuid.Parse(teacherID)
	payslip := &models.SalaryPayslip{
		TeacherID:       tUUID,
		Month:           month,
		Year:            year,
		TotalEarnings:   totalEarnings,
		TotalDeductions: totalDeductions,
		NetSalary:       netSalary,
		Details:         details,
	}

	// 5. Simpan ke database
	if err := s.repo.SavePayslip(payslip); err != nil {
		return nil, errors.New("Gagal menyimpan data slip gaji")
	}

	return payslip, nil
}

// Helper cerdas Golang untuk menentukan tipe komponen dari namanya
func isDeduction(name string) bool {
	// Pengecekan sederhana. Di dunia nyata, Anda bisa menggunakan regex atau mapping ke ID Component
	if len(name) >= 8 && name[:8] == "Potongan" {
		return true
	}
	return false
}
