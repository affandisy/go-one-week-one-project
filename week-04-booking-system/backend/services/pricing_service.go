package services

import (
	"errors"

	"github.com/affandisy/padel-booking-system/models"
	"github.com/affandisy/padel-booking-system/repositories"
	"github.com/google/uuid"
)

type PricingService interface {
	CreateRule(req CreatePricingRuleRequest) error
	GetRulesByCourt(courtID string) ([]models.PricingRule, error)
	UpdateRule(id string, req UpdatePricingRuleRequest) error
	DeleteRule(id string) error
}

type pricingService struct {
	repo repositories.PricingRepository
}

func NewPricingService(repo repositories.PricingRepository) PricingService {
	return &pricingService{repo}
}

// ================= DTO =================

type CreatePricingRuleRequest struct {
	CourtID    string  `json:"court_id"`
	DayType    string  `json:"day_type"`   // weekday, weekend
	StartTime  string  `json:"start_time"` // Format HH:MM, e.g., "08:00"
	EndTime    string  `json:"end_time"`   // Format HH:MM, e.g., "10:00"
	BasePrice  float64 `json:"base_price"`
	Multiplier float64 `json:"multiplier"`
}

type UpdatePricingRuleRequest struct {
	DayType    string  `json:"day_type"`
	StartTime  string  `json:"start_time"`
	EndTime    string  `json:"end_time"`
	BasePrice  float64 `json:"base_price"`
	Multiplier float64 `json:"multiplier"`
}

// ================= IMPLEMENTASI LOGIKA =================

func (s *pricingService) CreateRule(req CreatePricingRuleRequest) error {
	if req.BasePrice <= 0 {
		return errors.New("harga dasar harus lebih besar dari 0")
	}
	if req.DayType != "weekday" && req.DayType != "weekend" {
		return errors.New("tipe hari harus 'weekday' atau 'weekend'")
	}
	if req.Multiplier < 1.0 {
		req.Multiplier = 1.0 // Default ke 1.0 jika tidak diisi atau salah
	}

	courtUUID, err := uuid.Parse(req.CourtID)
	if err != nil {
		return errors.New("ID lapangan tidak valid")
	}

	rule := &models.PricingRule{
		CourtID:    courtUUID,
		DayType:    req.DayType,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		BasePrice:  req.BasePrice,
		Multiplier: req.Multiplier,
	}

	return s.repo.Create(rule)
}

func (s *pricingService) GetRulesByCourt(courtID string) ([]models.PricingRule, error) {
	return s.repo.FindByCourtID(courtID)
}

func (s *pricingService) UpdateRule(id string, req UpdatePricingRuleRequest) error {
	rule, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("aturan harga tidak ditemukan")
	}

	// Update kolom yang dikirim
	if req.DayType != "" {
		rule.DayType = req.DayType
	}
	if req.StartTime != "" {
		rule.StartTime = req.StartTime
	}
	if req.EndTime != "" {
		rule.EndTime = req.EndTime
	}
	if req.BasePrice > 0 {
		rule.BasePrice = req.BasePrice
	}
	if req.Multiplier >= 1.0 {
		rule.Multiplier = req.Multiplier
	}

	return s.repo.Update(rule)
}

func (s *pricingService) DeleteRule(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("aturan harga tidak ditemukan")
	}
	return s.repo.Delete(id)
}
