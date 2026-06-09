package usecase

import (
	"errors"
	"strings"

	"github.com/affandisy/hardware-erp/internal/core/domain"
	"github.com/affandisy/hardware-erp/internal/core/port"
	"github.com/affandisy/hardware-erp/internal/core/service"
	"github.com/google/uuid"
)

type componentUseCase struct {
	repo      port.ComponentRepository
	validator *service.SchemaValidator // Injeksi validator
}

// NewComponentUseCase menerima repository dan validator
func NewComponentUseCase(repo port.ComponentRepository, validator *service.SchemaValidator) port.ComponentUseCase {
	return &componentUseCase{repo, validator}
}

func (uc *componentUseCase) RegisterComponent(sku, name, category, manufacturer string, price float64, specs map[string]interface{}) (*domain.Component, error) {
	category = strings.TrimSpace(category)

	// 1. Validasi Skema JSON secara ketat sebelum menyentuh entitas
	if err := uc.validator.Validate(category, specs); err != nil {
		return nil, errors.New("gagal validasi spesifikasi: " + err.Error())
	}

	// 2. Validasi Duplikasi SKU
	existing, _ := uc.repo.FindBySKU(sku)
	if existing != nil {
		return nil, errors.New("komponen dengan SKU tersebut sudah ada")
	}

	// 3. Buat Entitas
	comp := &domain.Component{
		ID:           uuid.NewString(),
		SKU:          sku,
		Name:         name,
		Category:     category,
		Manufacturer: manufacturer,
		BasePrice:    price,
		StockOnHand:  0,
		Specs:        specs, // Data dipastikan aman dan berstruktur benar
	}

	// 4. Simpan ke Database
	if err := uc.repo.Save(comp); err != nil {
		return nil, errors.New("gagal menyimpan komponen ke database")
	}

	return comp, nil
}
