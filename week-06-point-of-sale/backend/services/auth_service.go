package services

import (
	"errors"
	"time"

	"github.com/affandisy/pos-system/models"
	"github.com/affandisy/pos-system/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Login(username, pin string) (string, *models.User, error)
}

type authService struct {
	repo      repositories.UserRepository
	jwtSecret string
}

func NewAuthService(repo repositories.UserRepository, secret string) AuthService {
	return &authService{repo, secret}
}

func (s *authService) Login(username, pin string) (string, *models.User, error) {
	// 1. Cari pengguna
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", nil, errors.New("username atau PIN salah")
	}

	// 2. Verifikasi PIN menggunakan Bcrypt
	if err := bcrypt.CompareHashAndPassword([]byte(user.PINHash), []byte(pin)); err != nil {
		return "", nil, errors.New("username atau PIN salah")
	}

	// 3. Buat JWT Token
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"role":     user.Role,                             // Sangat penting untuk Otorisasi nanti
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam (Shift harian)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, errors.New("gagal membuat sesi")
	}

	return tokenString, user, nil
}
