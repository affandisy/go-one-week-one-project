package services

import (
	"errors"

	"github.com/affandisy/padel-booking-system/repositories"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	UpdateProfile(userID string, req UpdateProfileRequest) error
	UpdatePassword(userID string, req UpdatePasswordRequest) error
}

type userService struct {
	repo repositories.UserRepository
}

func NewUserService(repo repositories.UserRepository) UserService {
	return &userService{repo}
}

// ================= DTO =================

type UpdateProfileRequest struct {
	FullName string `json:"full_name"`
	Whatsapp string `json:"whatsapp"`
}

type UpdatePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ================= IMPLEMENTASI LOGIKA =================

func (s *userService) UpdateProfile(userID string, req UpdateProfileRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	// Update field yang dikirim (bisa salah satu atau keduanya)
	if req.FullName != "" {
		user.FullName = req.FullName
	}
	if req.Whatsapp != "" {
		// Di dunia nyata, jika mengubah WA, mungkin perlu verifikasi OTP ulang.
		// Untuk MVP, kita langsung izinkan.
		user.Whatsapp = req.Whatsapp
	}

	return s.repo.Update(user)
}

func (s *userService) UpdatePassword(userID string, req UpdatePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	// 1. Validasi Password Lama
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.OldPassword)); err != nil {
		return errors.New("password lama yang Anda masukkan salah")
	}

	// 2. Hash Password Baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal memproses password baru")
	}

	// 3. Simpan Perubahan
	user.PasswordHash = string(hashedPassword)
	return s.repo.Update(user)
}
