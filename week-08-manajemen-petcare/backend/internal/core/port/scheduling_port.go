package port

import (
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
)

// Driven Port (Database)
type SchedulingRepository interface {
	CheckOverlap(petID string, start, end time.Time) (bool, error)
	SaveAppointment(app *domain.Appointment) error
}

// Driving Port (API)
type SchedulingUseCase interface {
	BookAppointment(petID, serviceID string, start, end time.Time) (*domain.Appointment, error)
}
