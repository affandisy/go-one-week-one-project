package services

import (
	"errors"
	"math/rand"
	"time"

	"github.com/affandi/belajar-bahasa/models"
	"github.com/affandi/belajar-bahasa/repositories"
	"github.com/google/uuid"
)

type AnswerRequest struct {
	MaterialID string `json:"material_id"`
	Selected   string `json:"selected"`
}

type QuizResultResponse struct {
	Score              float64        `json:"score"`
	Completed          bool           `json:"completed"`
	NextModuleUnlocked *models.Module `json:"next_module_unlocked,omitempty"`
}

type QuizService interface {
	GenerateQuizSession(moduleID string) ([]models.Material, error)
	EvaluateQuiz(userID, moduleID string, answers []AnswerRequest) (*QuizResultResponse, error)
}

type quizService struct {
	matRepo  repositories.MaterialRepository
	progRepo repositories.ProgressRepository
	modRepo  repositories.ModuleRepository
}

func NewQuizService(matRepo repositories.MaterialRepository, progRepo repositories.ProgressRepository, modRepo repositories.ModuleRepository) QuizService {
	return &quizService{matRepo, progRepo, modRepo}
}

// 1. Membuat Sesi Kuis (Maksimal 10 soal acak)
func (s *quizService) GenerateQuizSession(moduleID string) ([]models.Material, error) {
	materials, err := s.matRepo.GetByModuleID(moduleID)
	if err != nil {
		return nil, errors.New("gagal mengambil materi kuis")
	}

	// Filter hanya materi yang berupa kuis
	var quizItems []models.Material
	for _, m := range materials {
		if m.ContentType != "learn_card" {
			// Sembunyikan kunci jawaban agar tidak bocor ke frontend
			m.CorrectAnswer = ""
			quizItems = append(quizItems, m)
		}
	}

	// Acak urutan pertanyaan
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(quizItems), func(i, j int) { quizItems[i], quizItems[j] = quizItems[j], quizItems[i] })

	// Batasi maksimal 10 pertanyaan sesuai PRD
	if len(quizItems) > 10 {
		quizItems = quizItems[:10]
	}

	return quizItems, nil
}

// 2. Evaluasi Jawaban & Logika Unlocking
func (s *quizService) EvaluateQuiz(userID, moduleID string, answers []AnswerRequest) (*QuizResultResponse, error) {
	// Ambil kunci jawaban asli dari database
	materials, err := s.matRepo.GetByModuleID(moduleID)
	if err != nil || len(materials) == 0 {
		return nil, errors.New("materi tidak ditemukan")
	}

	answerMap := make(map[string]string)
	for _, m := range materials {
		answerMap[m.ID.String()] = m.CorrectAnswer
	}

	// Hitung Skor
	correctCount := 0
	for _, ans := range answers {
		if correctAnswer, exists := answerMap[ans.MaterialID]; exists {
			if ans.Selected == correctAnswer {
				correctCount++
			}
		}
	}

	totalQuestions := len(answers)
	if totalQuestions == 0 {
		return nil, errors.New("tidak ada jawaban yang disubmit")
	}

	score := (float64(correctCount) / float64(totalQuestions)) * 100

	// Ambil Progres User (Atau buat baru jika belum ada)
	progress, err := s.progRepo.GetUserProgress(userID, moduleID)
	if err != nil {
		uID, _ := uuid.Parse(userID)
		mID, _ := uuid.Parse(moduleID)
		progress = &models.UserProgress{
			UserID:   uID,
			ModuleID: mID,
		}
	}

	// Update record progres
	now := time.Now()
	progress.Attempts++
	progress.LastQuizDate = &now
	if score > progress.BestScore {
		progress.BestScore = score
	}

	var nextUnlocked *models.Module

	// PRD Rules: Jika Skor >= 80% dan belum pernah komplit
	if score >= 80 && !progress.Completed {
		progress.Completed = true

		// Cari modul level selanjutnya
		currentModule, _ := s.modRepo.FindByID(moduleID)
		nextModule, _ := s.progRepo.GetNextModule(currentModule.LevelOrder)

		if nextModule != nil {
			nextUnlocked = nextModule

			// Buat record kosong untuk modul selanjutnya agar terbaca "unlocked"
			nextProg := &models.UserProgress{
				UserID:   progress.UserID,
				ModuleID: nextModule.ID,
			}
			_ = s.progRepo.SaveUserProgress(nextProg)
		}
	}

	// Simpan progres ke database
	if err := s.progRepo.SaveUserProgress(progress); err != nil {
		return nil, errors.New("gagal menyimpan progres belajar")
	}

	return &QuizResultResponse{
		Score:              score,
		Completed:          progress.Completed,
		NextModuleUnlocked: nextUnlocked,
	}, nil
}
