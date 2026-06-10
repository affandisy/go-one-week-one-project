package usecase

import (
	"errors"
	"time"

	"github.com/affandisy/hardware-erp/internal/core/domain"
	"github.com/affandisy/hardware-erp/internal/core/port"
	"github.com/google/uuid"
)

type assemblyUseCase struct {
	compRepo port.ComponentRepository
	engine   port.AssemblyEngine // Injeksi Rule Engine
}

func NewAssemblyUseCase(repo port.ComponentRepository, engine port.AssemblyEngine) *assemblyUseCase {
	return &assemblyUseCase{repo, engine}
}

// Simulasi Validasi (Digunakan Frontend untuk tampilan real-time UI)
func (uc *assemblyUseCase) SimulateBuild(componentSKUs []string) (*domain.AssemblyBuild, error) {
	if len(componentSKUs) == 0 {
		return nil, errors.New("tidak ada komponen yang dipilih")
	}

	var selectedComponents []domain.Component

	// Ambil data detail komponen dari database berdasarkan SKU
	for _, sku := range componentSKUs {
		comp, err := uc.compRepo.FindBySKU(sku)
		if err != nil || comp == nil {
			return nil, errors.New("komponen dengan SKU " + sku + " tidak ditemukan di katalog")
		}
		selectedComponents = append(selectedComponents, *comp)
	}

	// Serahkan validasi berat kepada Engine
	isCompatible, notes, power, price := uc.engine.ValidateBuild(selectedComponents)

	build := &domain.AssemblyBuild{
		ID:              uuid.NewString(),
		BuildName:       "Draft Rakitan - " + time.Now().Format("02 Jan 2006"),
		Components:      selectedComponents,
		TotalPowerW:     power,
		TotalPrice:      price,
		IsCompatible:    isCompatible,
		ValidationNotes: notes,
	}

	return build, nil
}
