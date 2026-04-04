package services

import (
	"github.com/affandi/belajar-bahasa/models"
	"github.com/affandi/belajar-bahasa/repositories"
)

// DTO (Data Transfer Object) agar respons JSON lebih rapi untuk Frontend
type DashboardProgressResponse struct {
	ModuleID    string  `json:"module_id"`
	ModuleTitle string  `json:"module_title"`
	LevelOrder  int     `json:"level_order"`
	IsLocked    bool    `json:"is_locked"`
	Completed   bool    `json:"completed"`
	BestScore   float64 `json:"best_score"`
	Attempts    int     `json:"attempts"`
}

type ProgressService interface {
	GetDashboardProgress(userID string) ([]DashboardProgressResponse, error)
	ResetProgress(userID string) error
}

type progressService struct {
	progRepo repositories.ProgressRepository
	modRepo  repositories.ModuleRepository
}

func NewProgressService(progRepo repositories.ProgressRepository, modRepo repositories.ModuleRepository) ProgressService {
	return &progressService{progRepo, modRepo}
}

func (s *progressService) GetDashboardProgress(userID string) ([]DashboardProgressResponse, error) {
	// 1. Ambil semua modul yang ada di sistem (Alphabet, Vocab, Sentences)
	modules, err := s.modRepo.FindAll()
	if err != nil {
		return nil, err
	}

	// 2. Ambil progres user
	userProgresses, _ := s.progRepo.GetAllByUser(userID)

	// Buat map (kamus) agar mudah dicari berdasarkan ModuleID
	progMap := make(map[string]models.UserProgress)
	for _, p := range userProgresses {
		progMap[p.ModuleID.String()] = p
	}

	// 3. Gabungkan data Modul dan Progres menjadi satu respon komprehensif
	var dashboard []DashboardProgressResponse

	for _, mod := range modules {
		modID := mod.ID.String()

		isLocked := mod.IsLocked // Status bawaan (hanya level 1 yang false)
		completed := false
		bestScore := 0.0
		attempts := 0

		// Jika user punya progres di modul ini, timpa status default-nya
		if prog, exists := progMap[modID]; exists {
			isLocked = false // Jika ada di tabel progres, berarti sudah di-unlock
			completed = prog.Completed
			bestScore = prog.BestScore
			attempts = prog.Attempts
		}

		dashboard = append(dashboard, DashboardProgressResponse{
			ModuleID:    modID,
			ModuleTitle: mod.Title,
			LevelOrder:  mod.LevelOrder,
			IsLocked:    isLocked,
			Completed:   completed,
			BestScore:   bestScore,
			Attempts:    attempts,
		})
	}

	return dashboard, nil
}

func (s *progressService) ResetProgress(userID string) error {
	return s.progRepo.DeleteByUser(userID)
}
