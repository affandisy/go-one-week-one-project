package rest

import (
	"github.com/affandisy/petcare-system/internal/core/port"
	"github.com/gofiber/fiber/v2"
)

type NutritionHandler struct {
	useCase port.NutritionUseCase
}

func NewNutritionHandler(useCase port.NutritionUseCase) *NutritionHandler {
	return &NutritionHandler{useCase}
}

// Request Payload DTO
type CreateNutritionLogReq struct {
	PetID       string `json:"pet_id"`
	FoodBrand   string `json:"food_brand"`   // Contoh payload: "Orijen Cat & Kitten"
	Calories    int    `json:"calories"`     // Contoh payload: 250
	HealthNotes string `json:"health_notes"` // Contoh payload: "Transformasi bulu kucing lokal terlihat jauh lebih lebat dan pencernaan membaik setelah transisi diet."
}

func (h *NutritionHandler) CreateLog(c *fiber.Ctx) error {
	var req CreateNutritionLogReq
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	if err := h.useCase.RecordDiet(req.PetID, req.FoodBrand, req.Calories, req.HealthNotes); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Rekam gizi berhasil dicatat"})
}

func (h *NutritionHandler) GetLogs(c *fiber.Ctx) error {
	petID := c.Query("pet_id")
	if petID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Parameter pet_id wajib disertakan"})
	}

	logs, err := h.useCase.GetLogsByPet(petID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil riwayat gizi"})
	}

	return c.JSON(fiber.Map{"data": logs})
}
