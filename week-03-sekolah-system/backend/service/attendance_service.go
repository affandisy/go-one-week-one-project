package service

import (
	"errors"
	"time"

	"github.com/affandisy/school-system-sma/models"
	"github.com/affandisy/school-system-sma/repository"
	"github.com/google/uuid"
)

type AttendanceService interface {
	RecordStudentBatchAttendance(req BatchAttendanceRequest, recordedByID string) error
	GetStudentsForAttendance(classID string) ([]models.Student, error)
}

type attendanceService struct {
	repo repository.AttendanceRepository
}

func NewAttendanceService(repo repository.AttendanceRepository) AttendanceService {
	return &attendanceService{repo}
}

// DTO untuk input massal (satu kelas)
type BatchAttendanceRequest struct {
	ClassID string `json:"class_id"`
	Date    string `json:"date"` // Format: YYYY-MM-DD
	Records []struct {
		StudentID string `json:"student_id"`
		UserID    string `json:"user_id"` // User ID dari siswa tersebut
		Status    string `json:"status"`  // hadir, izin, sakit, alfa
		Notes     string `json:"notes"`
	} `json:"records"`
}

func (s *attendanceService) RecordStudentBatchAttendance(req BatchAttendanceRequest, recordedByID string) error {
	// 1. Validasi Tanggal
	attendanceDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return errors.New("Format tanggal tidak valid (gunakan YYYY-MM-DD)")
	}

	if attendanceDate.After(time.Now()) {
		return errors.New("Tidak bisa melakukan absensi untuk tanggal di masa depan")
	}

	recorderUUID, _ := uuid.Parse(recordedByID)

	// 2. Iterasi dan Validasi Data per Siswa (Golang Cerdas)
	for _, record := range req.Records {
		if record.Status != "hadir" && record.Status != "izin" && record.Status != "sakit" && record.Status != "alfa" {
			return errors.New("Status absensi tidak valid untuk siswa ID: " + record.StudentID)
		}

		userUUID, _ := uuid.Parse(record.UserID)

		attendance := &models.Attendance{
			UserID:         userUUID,
			AttendanceDate: attendanceDate,
			Status:         record.Status,
			Notes:          record.Notes,
			RecordedByID:   &recorderUUID,
		}

		// Jika hadir, catat waktu check-in (opsional)
		if record.Status == "hadir" {
			now := time.Now()
			attendance.CheckInTime = &now
		}

		// Simpan ke DB (Lewat Repo)
		if err := s.repo.UpsertAttendance(attendance); err != nil {
			return errors.New("Gagal menyimpan absensi untuk sebagian siswa")
		}
	}

	return nil
}

func (s *attendanceService) GetStudentsForAttendance(classID string) ([]models.Student, error) {
	return s.repo.GetStudentsByClass(classID)
}
