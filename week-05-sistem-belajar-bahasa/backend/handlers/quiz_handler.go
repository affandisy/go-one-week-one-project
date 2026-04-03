package handlers

import (
	"github.com/affandi/belajar-bahasa/services"
	"github.com/gofiber/fiber/v2"
)

type QuizHandler struct {
	service services.QuizService
}

func NewQuizHandler(service services.QuizService) *QuizHandler {
	return &QuizHandler{service}
}

func (h *QuizHandler) StartQuiz(c *fiber.Ctx) error {
	moduleID := c.Params("moduleId")

	questions, err := h.service.GenerateQuizSession(moduleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{
		"message": "Sesi kuis berhasil dibuat",
		"data": fiber.Map{
			"questions": questions,
		},
	})
}

func (h *QuizHandler) SubmitQuiz(c *fiber.Ctx) error {
	moduleID := c.Params("moduleId")
	userID := c.Locals("user_id").(string) // Didapat dari JWT Middleware

	var req struct {
		Answers []services.AnswerRequest `json:"answers"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format request tidak valid"})
	}

	result, err := h.service.EvaluateQuiz(userID, moduleID, req.Answers)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Kuis berhasil dievaluasi",
		"data":    result,
	})
}
