package service

import (
	"errors"
	"time"

	"github.com/affandisy/school-system-sma/models"
	"github.com/affandisy/school-system-sma/repository"
)

type AcademicService interface {
	CreateAcademicYear(req CreateAcademicYearRequest) (*models.AcademicYear, error)
	GetAllAcademicYears() ([]models.AcademicYear, error)
}

type academicService struct {
	repo repository.AcademicRepository
}

func NewAcademicService(repo repository.AcademicRepository) AcademicService {
	return &academicService{repo}
}

type CreateAcademicYearRequest struct {
	Name      string `json:"name"`
	StartDate string `json:"start_date"` // Format: YYYY-MM-DD
	EndDate   string `json:"end_date"`   // Format: YYYY-MM-DD
	IsActive  bool   `json:"is_active"`
}

func (s *academicService) CreateAcademicYear(req CreateAcademicYearRequest) (*models.AcademicYear, error) {
	// Parsing tanggal (Logika validasi di Golang)
	start, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return nil, errors.New("Format StartDate tidak valid")
	}

	end, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		return nil, errors.New("Format EndDate tidak valid")
	}

	if end.Before(start) {
		return nil, errors.New("Tanggal akhir tidak boleh sebelum tanggal mulai")
	}

	ay := &models.AcademicYear{
		Name:      req.Name,
		StartDate: start,
		EndDate:   end,
		IsActive:  req.IsActive,
	}

	// Jika User ingin mengaktifkan tahun ajaran baru ini
	if req.IsActive {
		// Gunakan Database Transaction agar aman
		tx := s.repo.BeginTransaction()

		// 1. Matikan semua tahun ajaran yang ada (Golang yang mengontrol alur ini)
		if err := s.repo.UpdateAcademicYearsStatus(tx, false); err != nil {
			tx.Rollback()
			return nil, errors.New("Gagal menonaktifkan tahun ajaran sebelumnya")
		}

		// 2. Simpan tahun ajaran baru
		if err := tx.Create(ay).Error; err != nil {
			tx.Rollback()
			return nil, errors.New("Gagal menyimpan tahun ajaran baru")
		}

		tx.Commit()
	} else {
		// Jika tidak aktif, simpan biasa
		if err := s.repo.CreateAcademicYear(ay); err != nil {
			return nil, errors.New("Gagal menyimpan tahun ajaran")
		}
	}

	return ay, nil
}

func (s *academicService) GetAllAcademicYears() ([]models.AcademicYear, error) {
	return s.repo.FindAllAcademicYears()
}
