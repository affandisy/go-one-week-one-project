package domain

import "time"

type User struct {
	ID           string
	Username     string
	PasswordHash string
	Role         string // "Receptionist", "Cashier", "Groomer", "Manager"
	CreatedAt    time.Time
}
