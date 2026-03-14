package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RecommendationCriteria struct {
	ID     uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name   string    `gorm:"not null" json:"name"`                     // "Nilai Akademik", "Kehadiran", dll
	Weight float64   `gorm:"type:decimal(3,2);not null" json:"weight"` // Total semua kriteria harus 1.0 (100%)
}

type StudentRecommendation struct {
	ID             uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	StudentID      uuid.UUID    `gorm:"not null" json:"student_id"`
	Student        Student      `gorm:"foreignKey:StudentID" json:"student"`
	Semester       int          `gorm:"not null" json:"semester"`
	AcademicYearID uuid.UUID    `gorm:"not null" json:"academic_year_id"`
	AcademicYear   AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academic_year"`
	Score          float64      `gorm:"type:decimal(5,4);not null" json:"score"` // Nilai akhir perhitungan
	Rank           int          `json:"rank"`
	IsSelected     bool         `gorm:"default:false" json:"is_selected"`
	CreatedAt      time.Time    `json:"created_at"`
}

func (m *RecommendationCriteria) BeforeCreate(tx *gorm.DB) (err error) { m.ID = uuid.New(); return }
func (m *StudentRecommendation) BeforeCreate(tx *gorm.DB) (err error)  { m.ID = uuid.New(); return }
