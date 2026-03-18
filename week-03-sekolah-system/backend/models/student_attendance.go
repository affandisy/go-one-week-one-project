package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Student struct {
	ID               uuid.UUID    `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID           uuid.UUID    `gorm:"uniqueIndex;not null" json:"user_id"`
	User             User         `gorm:"foreignKey:UserID" json:"user"`
	NIS              string       `gorm:"uniqueIndex;not null" json:"nis"`
	ClassID          *uuid.UUID   `json:"class_id"`
	Class            Class        `gorm:"foreignKey:ClassID" json:"class"`
	EnrollmentYearID uuid.UUID    `json:"enrollment_year_id"`
	EnrollmentYear   AcademicYear `gorm:"foreignKey:EnrollmentYearID" json:"enrollment_year"`
}

type Attendance struct {
	ID             uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	UserID         uuid.UUID `gorm:"not null;uniqueIndex:idx_user_date" json:"user_id"`
	AttendanceDate time.Time `gorm:"type:date;not null;uniqueIndex:idx_user_date" json:"attendance_date"`

	User         User       `gorm:"foreignKey:UserID" json:"user"`
	Status       string     `gorm:"not null" json:"status"`
	CheckInTime  *time.Time `json:"check_in_time"`
	Notes        string     `json:"notes"`
	RecordedByID *uuid.UUID `json:"recorded_by_id"`
	RecordedBy   User       `gorm:"foreignKey:RecordedByID" json:"recorded_by"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (m *Student) BeforeCreate(tx *gorm.DB) (err error)    { m.ID = uuid.New(); return }
func (m *Attendance) BeforeCreate(tx *gorm.DB) (err error) { m.ID = uuid.New(); return }
