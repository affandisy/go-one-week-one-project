package service

import (
	"errors"
	"sort"

	"github.com/affandisy/school-system-sma/models"
	"github.com/affandisy/school-system-sma/repository"
	"github.com/google/uuid"
)

type RecommendationService interface {
	CalculateRanking(classID string, semester int, academicYearID string) ([]models.StudentRecommendation, error)
}

type recommendationService struct {
	repo repository.RecommendationRepository
}

func NewRecommendationService(repo repository.RecommendationRepository) RecommendationService {
	return &recommendationService{repo}
}

func (s *recommendationService) CalculateRanking(classID string, semester int, academicYearID string) ([]models.StudentRecommendation, error) {
	// 1. Ambil Bobot Kriteria (Misal: Akademik 0.7, Kehadiran 0.3)
	criteria, err := s.repo.GetCriteria()
	if err != nil || len(criteria) == 0 {
		return nil, errors.New("Kriteria rekomendasi belum diatur")
	}

	var weightAkademik, weightKehadiran float64
	for _, c := range criteria {
		if c.Name == "Nilai Akademik" {
			weightAkademik = c.Weight
		}
		if c.Name == "Kehadiran" {
			weightKehadiran = c.Weight
		}
	}

	// 2. Ambil Daftar Siswa di Kelas tersebut
	students, err := s.repo.GetStudentsByClass(classID)
	if err != nil || len(students) == 0 {
		return nil, errors.New("Tidak ada siswa ditemukan di kelas tersebut")
	}

	var recommendations []models.StudentRecommendation
	ayUUID, _ := uuid.Parse(academicYearID)

	// 3. Kalkulasi per Siswa (Golang Cerdas)
	for _, student := range students {
		// Ambil data mentah
		avgGrade, _ := s.repo.GetAverageGrade(student.ID.String())
		attendanceRate, _ := s.repo.GetAttendanceRate(student.UserID.String())

		// Hitung skor akhir menggunakan rumus pembobotan (Simple Additive Weighting)
		finalScore := (avgGrade * weightAkademik) + (attendanceRate * weightKehadiran)

		recommendations = append(recommendations, models.StudentRecommendation{
			StudentID:      student.ID,
			Semester:       semester,
			AcademicYearID: ayUUID,
			Score:          finalScore,
			// Rank akan diisi nanti
		})
	}

	// 4. Urutkan berdasarkan Skor Tertinggi (Descending)
	sort.Slice(recommendations, func(i, j int) bool {
		return recommendations[i].Score > recommendations[j].Score
	})

	// 5. Berikan Peringkat (Rank)
	for i := range recommendations {
		recommendations[i].Rank = i + 1
	}

	// 6. Simpan hasil kalkulasi ke Database
	// Idealnya hapus dulu kalkulasi yang lama jika ada (agar tidak dobel)
	err = s.repo.SaveRecommendations(recommendations)
	if err != nil {
		return nil, errors.New("Gagal menyimpan hasil rekomendasi")
	}

	return recommendations, nil
}
