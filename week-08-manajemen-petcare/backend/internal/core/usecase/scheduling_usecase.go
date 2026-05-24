package usecase

import (
	"errors"
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/google/uuid"
)

type schedulingUseCase struct {
	repo port.SchedulingRepository
}

func NewSchedulingUseCase(repo port.SchedulingRepository) port.SchedulingUseCase {
	return &schedulingUseCase{repo}
}

func (uc *schedulingUseCase) BookAppointment(petID, serviceID string, start, end time.Time) (*domain.Appointment, error) {
	if start.After(end) || start.Equal(end) {
		return nil, errors.New("waktu selesai harus lebih besar dari waktu mulai")
	}

	// Aturan Bisnis: Cek bentrok jadwal untuk hewan yang sama (FR-007)
	isOverlap, err := uc.repo.CheckOverlap(petID, start, end)
	if err != nil {
		return nil, errors.New("gagal memvalidasi ketersediaan jadwal")
	}
	if isOverlap {
		return nil, errors.New("jadwal bentrok: hewan ini sudah memiliki jadwal perawatan di waktu tersebut")
	}

	app := &domain.Appointment{
		ID:        uuid.NewString(),
		PetID:     petID,
		ServiceID: serviceID,
		StartAt:   start,
		EndAt:     end,
		Status:    "Scheduled",
	}

	if err := uc.repo.SaveAppointment(app); err != nil {
		return nil, err
	}

	return app, nil
}
