package domain

import "time"

type Service struct {
	ID        string
	Name      string // Contoh: "Premium Grooming Kucing"
	BasePrice float64
}

type Appointment struct {
	ID        string
	PetID     string
	ServiceID string
	StartAt   time.Time
	EndAt     time.Time
	Status    string // "Scheduled", "CheckedIn", "InProgress", "Completed"
}
