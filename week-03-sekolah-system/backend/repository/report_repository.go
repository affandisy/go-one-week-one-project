package repository

import (
	"github.com/affandisy/school-system-sma/models"
	"gorm.io/gorm"
)

type ReportRepository interface {
	GetStudentGrades(studentID string) ([]models.Grade, error)
	GetStudentByID(studentID string) (*models.Student, error)
	UpsertReportCard(report *models.ReportCard) error
	GetReportCards() ([]models.ReportCard, error)
}

type reportRepository struct{ db *gorm.DB }

func NewReportRepository(db *gorm.DB) ReportRepository { return &reportRepository{db} }

func (r *reportRepository) GetStudentGrades(studentID string) ([]models.Grade, error) {
	var grades []models.Grade
	// Preload Subject untuk PDF
	err := r.db.Preload("Subject").Where("student_id = ?", studentID).Find(&grades).Error
	return grades, err
}

func (r *reportRepository) GetStudentByID(studentID string) (*models.Student, error) {
	var student models.Student
	err := r.db.Preload("User").Preload("Class").Where("id = ?", studentID).First(&student).Error
	return &student, err
}

func (r *reportRepository) UpsertReportCard(report *models.ReportCard) error {
	var existing models.ReportCard
	err := r.db.Where("student_id = ? AND semester = ? AND academic_year_id = ?",
		report.StudentID, report.Semester, report.AcademicYearID).First(&existing).Error

	if err == nil {
		report.ID = existing.ID
		return r.db.Save(report).Error
	}
	return r.db.Create(report).Error
}

func (r *reportRepository) GetReportCards() ([]models.ReportCard, error) {
	var reports []models.ReportCard
	err := r.db.Preload("Student").Preload("Student.User").Order("created_at DESC").Find(&reports).Error
	return reports, err
}
