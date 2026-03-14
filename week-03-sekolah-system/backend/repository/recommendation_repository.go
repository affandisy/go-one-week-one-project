package repository

import (
	"github.com/affandisy/school-system-sma/models"
	"gorm.io/gorm"
)

type RecommendationRepository interface {
	GetCriteria() ([]models.RecommendationCriteria, error)
	GetStudentsByClass(classID string) ([]models.Student, error)
	GetAverageGrade(studentID string) (float64, error)
	GetAttendanceRate(userID string) (float64, error)
	SaveRecommendations(recs []models.StudentRecommendation) error
	ClearPreviousRecommendations(classID string, semester int, academicYearID string) error
}

type recommendationRepository struct {
	db *gorm.DB
}

func NewRecommendationRepository(db *gorm.DB) RecommendationRepository {
	return &recommendationRepository{db}
}

func (r *recommendationRepository) GetStudentsByClass(classID string) ([]models.Student, error) {
	var students []models.Student
	err := r.db.
		Where("class_id = ?", classID).
		Find(&students).Error
	return students, err
}

func (r *recommendationRepository) ClearPreviousRecommendations(classID string, semester int, academicYearID string) error {
	studentSubQuery := r.db.
		Model(&models.Student{}).
		Select("id").
		Where("class_id = ?", classID)

	return r.db.
		Where("student_id IN (?)", studentSubQuery).
		Where("semester = ?", semester).
		Where("academic_year_id = ?", academicYearID).
		Delete(&models.StudentRecommendation{}).Error
}

// Contoh cara mengambil nilai rata-rata murni:
func (r *recommendationRepository) GetAverageGrade(studentID string) (float64, error) {
	var avg float64
	err := r.db.Model(&models.Grade{}).
		Where("student_id = ?", studentID).
		Select("COALESCE(AVG(score), 0)").Scan(&avg).Error
	return avg, err
}

// Persentase kehadiran (Hadir / Total Hari)
func (r *recommendationRepository) GetAttendanceRate(userID string) (float64, error) {
	var total, hadir int64
	r.db.Model(&models.Attendance{}).Where("user_id = ?", userID).Count(&total)
	if total == 0 {
		return 0, nil // Belum ada data absensi
	}
	r.db.Model(&models.Attendance{}).Where("user_id = ? AND status = ?", userID, "hadir").Count(&hadir)

	rate := float64(hadir) / float64(total) * 100.0
	return rate, nil
}

func (r *recommendationRepository) GetCriteria() ([]models.RecommendationCriteria, error) {
	var criteria []models.RecommendationCriteria
	err := r.db.Find(&criteria).Error
	return criteria, err
}

func (r *recommendationRepository) SaveRecommendations(recs []models.StudentRecommendation) error {
	return r.db.CreateInBatches(recs, 100).Error // Simpan masal
}
