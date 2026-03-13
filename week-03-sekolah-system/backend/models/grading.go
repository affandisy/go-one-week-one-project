package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Penugasan Guru ke Kelas & Mata Pelajaran
type TeacherAssignment struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TeacherID      uuid.UUID `gorm:"not null" json:"teacher_id"`
	ClassID        uuid.UUID `gorm:"not null" json:"class_id"`
	SubjectID      uuid.UUID `gorm:"not null" json:"subject_id"`
	AcademicYearID uuid.UUID `gorm:"not null" json:"academic_year_id"`
}

type GradeComponent struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name           string    `gorm:"not null" json:"name"`                     // e.g., "UAS", "Tugas Harian"
	Weight         float64   `gorm:"type:decimal(3,2);not null" json:"weight"` // e.g., 0.30
	AcademicYearID uuid.UUID `gorm:"not null" json:"academic_year_id"`
}

type Grade struct {
	ID          uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StudentID   uuid.UUID `gorm:"not null;uniqueIndex:idx_student_subject_comp" json:"student_id"`
	SubjectID   uuid.UUID `gorm:"not null;uniqueIndex:idx_student_subject_comp" json:"subject_id"`
	ComponentID uuid.UUID `gorm:"not null;uniqueIndex:idx_student_subject_comp" json:"component_id"`
	Score       float64   `gorm:"type:decimal(5,2);not null" json:"score"`
	Notes       string    `json:"notes"`
	CreatedByID uuid.UUID `gorm:"not null" json:"created_by_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ... tambahkan fungsi BeforeCreate untuk UUID seperti sebelumnya ...
func (m *Grade) BeforeCreate(tx *gorm.DB) (err error)             { m.ID = uuid.New(); return }
func (m *GradeComponent) BeforeCreate(tx *gorm.DB) (err error)    { m.ID = uuid.New(); return }
func (m *TeacherAssignment) BeforeCreate(tx *gorm.DB) (err error) { m.ID = uuid.New(); return }
