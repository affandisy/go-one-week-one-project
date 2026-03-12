package models

import (
	"time"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AcademicYear struct {
	ID        uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"` // e.g., "2025/2026"
	StartDate time.Time `gorm:"type:date;not null" json:"start_date"`
	EndDate   time.Time `gorm:"type:date;not null" json:"end_date"`
	IsActive  bool      `gorm:"default:false" json:"is_active"`
}

type StudyProgram struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code string    `gorm:"uniqueIndex;not null" json:"code"` // IPA, IPS
	Name string    `gorm:"not null" json:"name"`
}

type Class struct {
	ID                uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name              string       `gorm:"not null" json:"name"` // e.g., "XII IPA 1"
	StudyProgramID    uuid.UUID    `json:"study_program_id"`
	StudyProgram      StudyProgram `gorm:"foreignKey:StudyProgramID" json:"study_program"`
	AcademicYearID    uuid.UUID    `json:"academic_year_id"`
	AcademicYear      AcademicYear `gorm:"foreignKey:AcademicYearID" json:"academic_year"`
	HomeroomTeacherID *uuid.UUID   `json:"homeroom_teacher_id"` // Wali kelas
	HomeroomTeacher   User         `gorm:"foreignKey:HomeroomTeacherID" json:"homeroom_teacher"`
}

type Subject struct {
	ID             uuid.UUID     `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Code           string        `gorm:"uniqueIndex;not null" json:"code"`
	Name           string        `gorm:"not null" json:"name"`
	StudyProgramID *uuid.UUID    `json:"study_program_id"` // Null jika mata pelajaran umum (Agama, PKn)
	StudyProgram   *StudyProgram `gorm:"foreignKey:StudyProgramID" json:"study_program"`
}

// Hooks untuk UUID
func (m *AcademicYear) BeforeCreate(tx *gorm.DB) (err error) { m.ID = uuid.New(); return }
func (m *StudyProgram) BeforeCreate(tx *gorm.DB) (err error) { m.ID = uuid.New(); return }
func (m *Class) BeforeCreate(tx *gorm.DB) (err error) { m.ID = uuid.New(); return }
func (m *Subject) BeforeCreate(tx *gorm.DB) (err error) { m.ID = uuid.New(); return }