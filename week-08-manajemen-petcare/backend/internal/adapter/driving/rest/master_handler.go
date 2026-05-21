package rest

import (
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type MasterHandler struct {
	ownerUC port.OwnerUseCase
	petUC   port.PetUseCase
}

func NewMasterHandler(o port.OwnerUseCase, p port.PetUseCase) *MasterHandler {
	return &MasterHandler{o, p}
}

// --- OWNER HANDLERS ---
type CreateOwnerReq struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

func (h *MasterHandler) CreateOwner(c *fiber.Ctx) error {
	var req CreateOwnerReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Payload invalid"})
	}

	owner, err := h.ownerUC.RegisterOwner(req.Name, req.Phone)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"data": owner})
}

func (h *MasterHandler) GetOwners(c *fiber.Ctx) error {
	owners, err := h.ownerUC.ListOwners()
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data pemilik"})
	}
	return c.JSON(fiber.Map{"data": owners})
}

// --- PET HANDLERS ---
type CreatePetReq struct {
	OwnerID   string  `json:"owner_id"`
	Name      string  `json:"name"`
	Species   string  `json:"species"`
	Breed     string  `json:"breed"`
	Weight    float64 `json:"weight"`
	DietNotes string  `json:"diet_notes"`
}

func (h *MasterHandler) CreatePet(c *fiber.Ctx) error {
	var req CreatePetReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Payload invalid"})
	}

	pet, err := h.petUC.RegisterPet(req.OwnerID, req.Name, req.Species, req.Breed, req.Weight, req.DietNotes)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(201).JSON(fiber.Map{"data": pet})
}

func (h *MasterHandler) GetPets(c *fiber.Ctx) error {
	ownerID := c.Query("owner_id") // FR-001: Mengambil hewan berdasarkan pemilik
	if ownerID == "" {
		return c.Status(400).JSON(fiber.Map{"error": "owner_id parameter wajib diisi"})
	}

	pets, err := h.petUC.GetPetsByOwner(ownerID)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "Gagal mengambil data hewan peliharaan"})
	}
	return c.JSON(fiber.Map{"data": pets})
}
