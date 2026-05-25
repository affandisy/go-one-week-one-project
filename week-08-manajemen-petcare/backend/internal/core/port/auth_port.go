package port

import "github.com/affandisy/petcare-system/internal/core/domain"

// Driven Port (Database)
type UserRepository interface {
	FindByUsername(username string) (*domain.User, error)
	Save(user *domain.User) error
}

// Driving Port (API)
type AuthUseCase interface {
	RegisterUser(username, password, role string) error
	Login(username, password string) (string, error) // Mengembalikan JWT Token
}
