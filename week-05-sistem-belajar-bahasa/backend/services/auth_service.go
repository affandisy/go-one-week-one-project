package services

import (
	"errors"
	"time"

	"github.com/affandi/belajar-bahasa/models"
	"github.com/affandi/belajar-bahasa/repositories"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, password string) error
	Login(username, password string) (string, *models.User, error)
}

type authService struct {
	repo      repositories.UserRepository
	jwtSecret string
}

func NewAuthService(repo repositories.UserRepository, secret string) AuthService {
	return &authService{repo, secret}
}

func (s *authService) Register(username, password string) error {
	// 1. Validasi input dasar
	if username == "" || len(password) < 6 {
		return errors.New("username tidak boleh kosong dan password minimal 6 karakter")
	}

	// 2. Hash Password menggunakan bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal memproses password")
	}

	// 3. Simpan ke Database
	user := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return errors.New("username sudah digunakan, silakan pilih yang lain")
	}

	return nil
}

func (s *authService) Login(username, password string) (string, *models.User, error) {
	// 1. Cari user berdasarkan username
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", nil, errors.New("username atau password salah")
	}

	// 2. Bandingkan password yang diinput dengan hash di database
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", nil, errors.New("username atau password salah")
	}

	// 3. Generate JWT Token
	claims := jwt.MapClaims{
		"user_id":  user.ID.String(),
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(), // Token berlaku 7 hari
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, errors.New("gagal membuat token akses")
	}

	return tokenString, user, nil
}
