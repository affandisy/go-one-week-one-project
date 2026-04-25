package services

import (
	"errors"

	"github.com/affandisy/pos-system/models"
	"github.com/affandisy/pos-system/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserRequest struct {
	Username string `json:"username"`
	PIN      string `json:"pin"`
	Role     string `json:"role"` // "cashier", "owner", "admin"
}

type UserService interface {
	GetAllUsers() ([]models.User, error)
	CreateUser(req UserRequest) error
	UpdateUser(id string, req UserRequest) error
	DeleteUser(id string) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

func (s *userService) GetAllUsers() ([]models.User, error) {
	return s.repo.FindAll()
}

func (s *userService) CreateUser(req UserRequest) error {
	if req.Username == "" || req.PIN == "" || req.Role == "" {
		return errors.New("username, PIN, dan role wajib diisi")
	}

	if req.Role != "cashier" && req.Role != "owner" && req.Role != "admin" {
		return errors.New("role tidak valid (pilih: cashier, owner, admin)")
	}

	// Cek apakah username sudah ada
	existing, _ := s.repo.FindByUsername(req.Username)
	if existing != nil {
		return errors.New("username sudah digunakan")
	}

	// Enkripsi PIN Kasir
	hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengamankan PIN")
	}

	user := &models.User{
		Username: req.Username,
		PINHash:  string(hashedPIN),
		Role:     req.Role,
	}

	return s.repo.Create(user)
}

func (s *userService) UpdateUser(id string, req UserRequest) error {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("pengguna tidak ditemukan")
	}

	// Jika username diubah, pastikan tidak bentrok dengan milik orang lain
	if req.Username != "" && req.Username != user.Username {
		existing, _ := s.repo.FindByUsername(req.Username)
		if existing != nil {
			return errors.New("username sudah digunakan pengguna lain")
		}
		user.Username = req.Username
	}

	if req.Role != "" {
		if req.Role != "cashier" && req.Role != "owner" && req.Role != "admin" {
			return errors.New("role tidak valid")
		}
		user.Role = req.Role
	}

	// Jika kolom PIN diisi, berarti admin ingin mereset PIN pengguna tersebut
	if req.PIN != "" {
		hashedPIN, err := bcrypt.GenerateFromPassword([]byte(req.PIN), bcrypt.DefaultCost)
		if err != nil {
			return errors.New("gagal mengamankan PIN baru")
		}
		user.PINHash = string(hashedPIN)
	}

	return s.repo.Update(user)
}

func (s *userService) DeleteUser(id string) error {
	// Fitur Keamanan: Cek apakah ini user terakhir (jangan sampai owner menghapus dirinya sendiri dan tidak ada yang bisa login)
	count, _ := s.repo.Count()
	if count <= 1 {
		return errors.New("tidak dapat menghapus pengguna terakhir di sistem")
	}

	return s.repo.Delete(id)
}
