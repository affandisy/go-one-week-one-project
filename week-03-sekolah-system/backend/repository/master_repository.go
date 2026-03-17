package repository

import (
	"github.com/affandisy/school-system-sma/models"
	"gorm.io/gorm"
)

type MasterRepository interface {
	// Kelas
	CreateClass(c *models.Class) error
	GetClasses() ([]models.Class, error)

	// Mata Pelajaran
	CreateSubject(s *models.Subject) error
	GetSubjects() ([]models.Subject, error)

	// Siswa
	CreateStudent(s *models.Student) error
	GetStudents() ([]models.Student, error)
}

type masterRepository struct {
	db *gorm.DB
}

func NewMasterRepository(db *gorm.DB) MasterRepository {
	return &masterRepository{db}
}

func (r *masterRepository) CreateClass(c *models.Class) error { return r.db.Create(c).Error }
func (r *masterRepository) GetClasses() ([]models.Class, error) {
	var classes []models.Class
	err := r.db.Preload("StudyProgram").Preload("AcademicYear").Preload("HomeroomTeacher").Find(&classes).Error
	return classes, err
}

func (r *masterRepository) CreateSubject(s *models.Subject) error { return r.db.Create(s).Error }
func (r *masterRepository) GetSubjects() ([]models.Subject, error) {
	var subjects []models.Subject
	err := r.db.Preload("StudyProgram").Find(&subjects).Error
	return subjects, err
}

func (r *masterRepository) CreateStudent(s *models.Student) error { return r.db.Create(s).Error }
func (r *masterRepository) GetStudents() ([]models.Student, error) {
	var students []models.Student
	// Preload User untuk mendapatkan nama dan email siswa
	err := r.db.Preload("User").Preload("Class").Find(&students).Error
	return students, err
}
