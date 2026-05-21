package usecase

import (
	"errors"
	"strings"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/google/uuid"
)

type petUseCase struct {
	petRepo   port.PetRepository
	ownerRepo port.OwnerRepository // Butuh ini untuk validasi eksistensi owner
}

func NewPetUseCase(p port.PetRepository, o port.OwnerRepository) port.PetUseCase {
	return &petUseCase{p, o}
}

func (uc *petUseCase) RegisterPet(ownerID, name, species, breed string, weight float64, dietNotes string) (*domain.Pet, error) {
	if strings.TrimSpace(name) == "" {
		return nil, errors.New("nama hewan peliharaan wajib diisi")
	}

	// Validasi apakah OwnerID benar-benar ada
	_, err := uc.ownerRepo.FindByID(ownerID)
	if err != nil {
		return nil, errors.New("pemilik tidak ditemukan")
	}

	pet := &domain.Pet{
		ID:        uuid.NewString(),
		OwnerID:   ownerID,
		Name:      name,
		Species:   species,
		Breed:     breed,
		Weight:    weight,
		DietNotes: dietNotes,
	}

	if err := uc.petRepo.Save(pet); err != nil {
		return nil, err
	}
	return pet, nil
}

func (uc *petUseCase) GetPetsByOwner(ownerID string) ([]domain.Pet, error) {
	return uc.petRepo.FindByOwnerID(ownerID)
}
