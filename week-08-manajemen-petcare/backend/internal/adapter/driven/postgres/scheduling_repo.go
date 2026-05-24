package postgres

import (
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"gorm.io/gorm"
)

// DTO GORM
type AppointmentModel struct {
	ID        string    `gorm:"type:uuid;primaryKey"`
	PetID     string    `gorm:"type:uuid;index"`
	ServiceID string    `gorm:"type:uuid;index"`
	StartAt   time.Time `gorm:"index"`
	EndAt     time.Time
	Status    string
}

func (AppointmentModel) TableName() string { return "appointments" }

type schedulingRepository struct {
	db *gorm.DB
}

func NewSchedulingRepository(db *gorm.DB) port.SchedulingRepository {
	db.AutoMigrate(&AppointmentModel{}) // Pastikan Anda juga auto-migrate ServiceModel nanti
	return &schedulingRepository{db}
}

func (r *schedulingRepository) CheckOverlap(petID string, start, end time.Time) (bool, error) {
	var count int64

	// Cek apakah ada jadwal dengan PetID yang sama di mana waktunya bersinggungan
	err := r.db.Model(&AppointmentModel{}).
		Where("pet_id = ? AND start_at < ? AND end_at > ? AND status != 'Cancelled'", petID, end, start).
		Count(&count).Error

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *schedulingRepository) SaveAppointment(app *domain.Appointment) error {
	model := AppointmentModel{
		ID:        app.ID,
		PetID:     app.PetID,
		ServiceID: app.ServiceID,
		StartAt:   app.StartAt,
		EndAt:     app.EndAt,
		Status:    app.Status,
	}
	return r.db.Create(&model).Error
}
