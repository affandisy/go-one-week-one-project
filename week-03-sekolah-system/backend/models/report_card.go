package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ReportCard struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StudentID      uuid.UUID    `gorm:"not null" json:"student_id"`
	Student        Student      `gorm:"foreignKey:StudentID" json:"student"`
	Semester       int          `gorm:"not null" json:"semester"`
	AcademicYearID uuid.UUID    `gorm:"not null" json:"academic_year_id"`
	AcademicYear   AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academic_year"`
	FinalScore     float64      `gorm:"type:decimal(5,2)" json:"final_score"`
	PdfURL         string       `gorm:"type:text" json:"pdf_url"`
	ApprovedByID   *uuid.UUID   `json:"approved_by_id"`
	ApprovedAt     *time.Time   `json:"approved_at"`
	CreatedAt      time.Time    `json:"created_at"`
}

func (m *ReportCard) BeforeCreate(tx *gorm.DB) (err error) { m.ID = uuid.New(); return }
