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
	attendanceDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return errors.New("Format tanggal tidak valid")
	}
	if attendanceDate.After(time.Now()) {
		return errors.New("Tidak bisa mengabsen untuk masa depan")
	}

	recorderUUID, _ := uuid.Parse(recordedByID)

	// Siapkan penampung (Slice)
	var attendances []models.Attendance

	for _, record := range req.Records {
		if record.Status != "hadir" && record.Status != "izin" && record.Status != "sakit" && record.Status != "alfa" {
			return errors.New("Status absensi tidak valid")
		}

		userUUID, _ := uuid.Parse(record.UserID)
		att := models.Attendance{
			UserID:         userUUID,
			AttendanceDate: attendanceDate,
			Status:         record.Status,
			Notes:          record.Notes,
			RecordedByID:   &recorderUUID,
		}

		if record.Status == "hadir" {
			now := time.Now()
			att.CheckInTime = &now
		}

		// Masukkan ke penampung, BUKAN langsung di-save
		attendances = append(attendances, att)
	}

	// Eksekusi 1 kueri massal yang sangat cepat
	if err := s.repo.UpsertBatch(attendances); err != nil {
		return errors.New("Gagal menyimpan absensi massal")
	}

	return nil
}

func (s *attendanceService) GetStudentsForAttendance(classID string) ([]models.Student, error) {
	return s.repo.GetStudentsByClass(classID)
}
