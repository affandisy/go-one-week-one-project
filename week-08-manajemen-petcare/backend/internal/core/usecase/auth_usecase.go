package usecase

import (
	"errors"
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// Di lingkungan produksi, gunakan os.Getenv("JWT_SECRET")
var jwtSecretKey = []byte("super-rahasia-petcare-123")

type authUseCase struct {
	repo port.UserRepository
}

func NewAuthUseCase(repo port.UserRepository) port.AuthUseCase {
	return &authUseCase{repo}
}

func (uc *authUseCase) RegisterUser(username, password, role string) error {
	// Validasi peran (Role)
	validRoles := map[string]bool{"Receptionist": true, "Cashier": true, "Groomer": true, "Manager": true}
	if !validRoles[role] {
		return errors.New("peran pengguna tidak valid")
	}

	// Hash password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi kata sandi")
	}

	user := &domain.User{
		ID:           uuid.NewString(),
		Username:     username,
		PasswordHash: string(hashedPassword),
		Role:         role,
		CreatedAt:    time.Now(),
	}

	return uc.repo.Save(user)
}

func (uc *authUseCase) Login(username, password string) (string, error) {
	// 1. Cari pengguna di database
	user, err := uc.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("kredensial tidak valid")
	}

	// 2. Verifikasi kata sandi
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("kredensial tidak valid")
	}

	// 3. Buat JWT Token
	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtSecretKey)
	if err != nil {
		return "", errors.New("gagal menghasilkan token")
	}

	return tokenString, nil
}
