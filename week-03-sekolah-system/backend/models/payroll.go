package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SalaryComponent struct {
	ID   uuid.UUID `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	Name string    `gorm:"not null" json:"name"` // Tunjangan Jabatan, Potongan Telat, dll
	Type string    `gorm:"not null" json:"type"` // "earning" atau "deduction"
}

// Konfigurasi gaji per guru
type TeacherSalaryConfig struct {
	ID         uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TeacherID  uuid.UUID              `gorm:"uniqueIndex;not null" json:"teacher_id"`
	Teacher    User                   `gorm:"foreignKey:TeacherID" json:"teacher"`
	BaseSalary float64                `gorm:"type:decimal(12,2);not null" json:"base_salary"`
	Config     map[string]interface{} `gorm:"type:jsonb;serializer:json;not null;default:'{}'" json:"config"`
	// Contoh Config: {"Tunjangan Wali Kelas": 500000, "Potongan Koperasi": 100000}
}

// Slip Gaji Bulanan
type SalaryPayslip struct {
	ID              uuid.UUID              `gorm:"type:uuid;default:gen_random_uuid();primaryKey" json:"id"`
	TeacherID       uuid.UUID              `gorm:"not null" json:"teacher_id"`
	Teacher         User                   `gorm:"foreignKey:TeacherID" json:"teacher"`
	Month           int                    `gorm:"not null" json:"month"`
	Year            int                    `gorm:"not null" json:"year"`
	TotalEarnings   float64                `gorm:"type:decimal(12,2);not null" json:"total_earnings"`
	TotalDeductions float64                `gorm:"type:decimal(12,2);not null" json:"total_deductions"`
	NetSalary       float64                `gorm:"type:decimal(12,2);not null" json:"net_salary"`
	Details         map[string]interface{} `gorm:"type:jsonb;serializer:json" json:"details"`
	ApprovedByID    *uuid.UUID             `json:"approved_by_id"`
	ApprovedAt      *time.Time             `json:"approved_at"`
	CreatedAt       time.Time              `json:"created_at"`
}

// Hooks untuk UUID
func (m *SalaryComponent) BeforeCreate(tx *gorm.DB) (err error)     { m.ID = uuid.New(); return }
func (m *TeacherSalaryConfig) BeforeCreate(tx *gorm.DB) (err error) { m.ID = uuid.New(); return }
func (m *SalaryPayslip) BeforeCreate(tx *gorm.DB) (err error)       { m.ID = uuid.New(); return }
