package repository

import (
	"github.com/affandisy/school-system-sma/models"
	"gorm.io/gorm"
)

type GradeRepository interface {
	// Cek penugasan
	CheckTeacherAssignment(teacherID, classID, subjectID string) (bool, error)
	// Simpan Nilai (Upsert: Update kalau ada, Insert kalau belum)
	UpsertGrade(grade *models.Grade) error
}

type gradeRepository struct {
	db *gorm.DB
}

func NewGradeRepository(db *gorm.DB) GradeRepository {
	return &gradeRepository{db}
}

func (r *gradeRepository) CheckTeacherAssignment(teacherID, classID, subjectID string) (bool, error) {
	var count int64
	err := r.db.Model(&models.TeacherAssignment{}).
		Where("teacher_id = ? AND class_id = ? AND subject_id = ?", teacherID, classID, subjectID).
		Count(&count).Error
	return count > 0, err
}

func (r *gradeRepository) UpsertGrade(grade *models.Grade) error {
	// GORM Save akan melakukan Insert atau Update berdasarkan Primary Key / Unique Constraints
	// Untuk Upsert yang aman, kita cari dulu datanya
	var existing models.Grade
	err := r.db.Where("student_id = ? AND subject_id = ? AND component_id = ?",
		grade.StudentID, grade.SubjectID, grade.ComponentID).First(&existing).Error

	if err == nil {
		// Update jika sudah ada
		grade.ID = existing.ID
		return r.db.Save(grade).Error
	}

	// Create jika belum ada
	return r.db.Create(grade).Error
}
