package repository

import (
	"github.com/affandisy/school-system-sma/models"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type AttendanceRepository interface {
	UpsertAttendance(attendance *models.Attendance) error
	GetStudentsByClass(classID string) ([]models.Student, error)
	UpsertBatch(attendances []models.Attendance) error
}

type attendanceRepository struct {
	db *gorm.DB
}

func NewAttendanceRepository(db *gorm.DB) AttendanceRepository {
	return &attendanceRepository{db}
}

func (r *attendanceRepository) UpsertAttendance(attendance *models.Attendance) error {
	var existing models.Attendance

	// Cari apakah user ini sudah punya record absensi di tanggal yang sama
	err := r.db.Where("user_id = ? AND DATE(attendance_date) = DATE(?)",
		attendance.UserID, attendance.AttendanceDate).First(&existing).Error

	if err == nil {
		// Jika sudah ada, update status dan catatan
		attendance.ID = existing.ID
		return r.db.Save(attendance).Error
	}

	// Jika belum ada, buat baru
	return r.db.Create(attendance).Error
}

func (r *attendanceRepository) GetStudentsByClass(classID string) ([]models.Student, error) {
	var students []models.Student
	// Preload User untuk mendapatkan nama siswa
	err := r.db.Preload("User").Where("class_id = ?", classID).Find(&students).Error
	return students, err
}

func (r *attendanceRepository) UpsertBatch(attendances []models.Attendance) error {
	if len(attendances) == 0 {
		return nil
	}

	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}, {Name: "attendance_date"}},
		DoUpdates: clause.AssignmentColumns([]string{"status", "notes", "check_in_time", "recorded_by_id"}),
	}).Create(&attendances).Error
}
