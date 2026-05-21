package postgres

import (
	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"gorm.io/gorm"
)

// --- GORM MODELS ---
type OwnerModel struct {
	ID    string `gorm:"type:uuid;primaryKey"`
	Name  string `gorm:"not null"`
	Phone string
	Pets  []PetModel `gorm:"foreignKey:OwnerID"` // Relasi One-to-Many
}

func (OwnerModel) TableName() string { return "owners" }

type PetModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	OwnerID   string `gorm:"type:uuid;index"`
	Name      string `gorm:"not null"`
	Species   string
	Breed     string
	Weight    float64
	DietNotes string
}

func (PetModel) TableName() string { return "pets" }

// --- REPOSITORY IMPLEMENTATION ---
type ownerRepository struct {
	db *gorm.DB
}

type petRepository struct {
	db *gorm.DB
}

// Menggabungkan implementasi OwnerRepository dan PetRepository
func NewMasterRepository(db *gorm.DB) (port.OwnerRepository, port.PetRepository) {
	db.AutoMigrate(&OwnerModel{}, &PetModel{}) // Auto migrasi skema
	ownerRepo := &ownerRepository{db}
	petRepo := &petRepository{db}
	return ownerRepo, petRepo
}

// -- Owner Repo --
func (r *ownerRepository) Save(owner *domain.Owner) error {
	model := OwnerModel{
		ID:    owner.ID,
		Name:  owner.Name,
		Phone: owner.Phone,
	}
	return r.db.Save(&model).Error
}

func (r *ownerRepository) FindAll() ([]domain.Owner, error) {
	var models []OwnerModel
	// Gunakan Preload agar data Pets langsung ikut tertarik saat memanggil Owner
	if err := r.db.Preload("Pets").Find(&models).Error; err != nil {
		return nil, err
	}

	var owners []domain.Owner
	for _, m := range models {
		var pets []domain.Pet
		for _, p := range m.Pets {
			pets = append(pets, domain.Pet{ID: p.ID, Name: p.Name, Species: p.Species})
		}
		owners = append(owners, domain.Owner{ID: m.ID, Name: m.Name, Phone: m.Phone, Pets: pets})
	}
	return owners, nil
}

func (r *ownerRepository) FindByID(id string) (*domain.Owner, error) {
	var model OwnerModel
	if err := r.db.First(&model, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &domain.Owner{ID: model.ID, Name: model.Name, Phone: model.Phone}, nil
}

// -- Pet Repo --
func (r *petRepository) Save(pet *domain.Pet) error {
	model := PetModel{
		ID:        pet.ID,
		OwnerID:   pet.OwnerID,
		Name:      pet.Name,
		Species:   pet.Species,
		Breed:     pet.Breed,
		Weight:    pet.Weight,
		DietNotes: pet.DietNotes,
	}
	return r.db.Save(&model).Error
}

func (r *petRepository) FindByOwnerID(ownerID string) ([]domain.Pet, error) {
	var models []PetModel
	if err := r.db.Where("owner_id = ?", ownerID).Find(&models).Error; err != nil {
		return nil, err
	}

	var pets []domain.Pet
	for _, m := range models {
		pets = append(pets, domain.Pet{
			ID:        m.ID,
			OwnerID:   m.OwnerID,
			Name:      m.Name,
			Species:   m.Species,
			Breed:     m.Breed,
			Weight:    m.Weight,
			DietNotes: m.DietNotes,
		})
	}
	return pets, nil
}
