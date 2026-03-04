package services

import (
	"errors"
	"os"
	"time"

	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/models"
	"github.com/affandisy/go-one-week-one-project/week-02-warehouse-system/repositories"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterRequest struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserService interface {
	Register(req RegisterRequest) (*models.User, error)
	Login(req LoginRequest) (string, error)
	UpdateUserAvatar(userID uint, imagePath string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Register(req RegisterRequest) (*models.User, error) {
	if req.Email == "" || req.Password == "" || req.Name == "" {
		return nil, errors.New("Nama, email, dan Password wajid diisi")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("Gagal mengenkripsi password")
	}

	if req.Role == "" {
		req.Role = "operator"
	}

	user := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
		Phone:    req.Phone,
		Role:     req.Role,
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, errors.New("Email mungkin sudah terdaftar")
	}

	return &user, nil
}

func (s *userService) Login(req LoginRequest) (string, error) {
	user, err := s.repo.FindByEmail(req.Email)
	if err != nil {
		return "", errors.New("Email atau password salah")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", errors.New("Email atau password salah")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", errors.New("Gagal membuat token otentikasi")
	}

	return tokenString, nil
}

func (s *userService) UpdateUserAvatar(userID uint, imagePath string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("User tidak ditemukan")
	}

	user.AvatarURL = &imagePath

	return s.repo.Update(user)
}
