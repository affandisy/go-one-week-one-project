package services

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/affandisy/padel-booking-system/models"
	"github.com/affandisy/padel-booking-system/repositories"
	"github.com/golang-jwt/jwt/v5"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(email, whatsapp, fullName, password string) error
	VerifyOTP(email, otp string) error
	Login(email, password string) (string, error)
}

type authService struct {
	repo      repositories.UserRepository
	jwtSecret string
}

func NewAuthService(repo repositories.UserRepository, secret string) AuthService {
	return &authService{repo, secret}
}

func (s *authService) Register(email, whatsapp, fullName, password string) error {
	// 1. Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("gagal mengenkripsi password")
	}

	// 2. Buat 6-digit OTP (Simulasi)
	otpCode := fmt.Sprintf("%06d", rand.Intn(1000000))
	otpExpiry := time.Now().Add(5 * time.Minute) // Berlaku 5 menit

	// 3. Simpan ke Database (Status IsVerified = false)
	user := &models.User{
		Email:        email,
		Whatsapp:     whatsapp,
		FullName:     fullName,
		PasswordHash: string(hashedPassword),
		Role:         "customer",
		IsVerified:   false,
		OTP:          otpCode,
		OTPExpiry:    otpExpiry,
	}

	if err := s.repo.Create(user); err != nil {
		return errors.New("email atau whatsapp sudah terdaftar")
	}

	// 4. Simulasi Kirim OTP (Di Phase 3 nanti ini akan diganti dengan API WhatsApp Twilio)
	log.Printf("[MOCK WHATSAPP] Kode OTP untuk %s adalah: %s", whatsapp, otpCode)

	return nil
}

func (s *authService) VerifyOTP(email, otp string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("user tidak ditemukan")
	}

	if user.IsVerified {
		return errors.New("akun sudah diverifikasi sebelumnya")
	}

	if user.OTP != otp {
		return errors.New("kode OTP salah")
	}

	if time.Now().After(user.OTPExpiry) {
		return errors.New("kode OTP sudah kedaluwarsa")
	}

	// Aktivasi akun dan bersihkan OTP
	user.IsVerified = true
	user.OTP = ""

	return s.repo.Update(user)
}

func (s *authService) Login(email, password string) (string, error) {
	// 1. Cari User
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("email atau password salah")
	}

	// 2. Cek Status Verifikasi
	if !user.IsVerified {
		return "", errors.New("akun belum diverifikasi, silakan masukkan OTP")
	}

	// 3. Cek Password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("email atau password salah")
	}

	// 4. Generate JWT Token
	claims := jwt.MapClaims{
		"user_id": user.ID.String(),
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(), // Berlaku 24 Jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", errors.New("gagal membuat token autentikasi")
	}

	return tokenString, nil
}
