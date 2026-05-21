package port

import "github.com/affandisy/petcare-system/internal/core/domain"

// --- DRIVEN PORTS (Database) ---
type OwnerRepository interface {
	Save(owner *domain.Owner) error
	FindAll() ([]domain.Owner, error)
	FindByID(id string) (*domain.Owner, error)
}

type PetRepository interface {
	Save(pet *domain.Pet) error
	FindByOwnerID(ownerID string) ([]domain.Pet, error)
}

// --- DRIVING PORTS (API) ---
type OwnerUseCase interface {
	RegisterOwner(name, phone string) (*domain.Owner, error)
	ListOwners() ([]domain.Owner, error)
}

type PetUseCase interface {
	RegisterPet(ownerID, name, species, breed string, weight float64, dietNotes string) (*domain.Pet, error)
	GetPetsByOwner(ownerID string) ([]domain.Pet, error)
}
