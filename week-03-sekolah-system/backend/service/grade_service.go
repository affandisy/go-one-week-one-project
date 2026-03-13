package service

import (
	"errors"

	"github.com/affandisy/school-system-sma/models"
	"github.com/affandisy/school-system-sma/repository"
	"github.com/google/uuid"
)

type GradeService interface {
	InputGrade(req InputGradeRequest, teacherID string) error
}

type gradeService struct {
	gradeRepo repository.GradeRepository
}

func NewGradeService(gradeRepo repository.GradeRepository) GradeService {
	return &gradeService{gradeRepo}
}

type InputGradeRequest struct {
	ClassID     string  `json:"class_id"`
	StudentID   string  `json:"student_id"`
	SubjectID   string  `json:"subject_id"`
	ComponentID string  `json:"component_id"`
	Score       float64 `json:"score"`
	Notes       string  `json:"notes"`
}

func (s *gradeService) InputGrade(req InputGradeRequest, teacherID string) error {
	// 1. Validasi nilai
	if req.Score < 0 || req.Score > 100 {
		return errors.New("Nilai harus berada di antara 0 dan 100")
	}

	// 2. Validasi Otoritas Guru (Business Logic constraint)
	// Golang bertanya ke Repo: Apakah guru ini punya izin ngajar mapel ini di kelas ini?
	isAssigned, err := s.gradeRepo.CheckTeacherAssignment(teacherID, req.ClassID, req.SubjectID)
	if err != nil || !isAssigned {
		return errors.New("Akses ditolak: Anda tidak ditugaskan mengajar mata pelajaran ini di kelas tersebut")
	}

	// 3. Konversi ID
	studentUUID, _ := uuid.Parse(req.StudentID)
	subjectUUID, _ := uuid.Parse(req.SubjectID)
	compUUID, _ := uuid.Parse(req.ComponentID)
	teacherUUID, _ := uuid.Parse(teacherID)

	grade := &models.Grade{
		StudentID:   studentUUID,
		SubjectID:   subjectUUID,
		ComponentID: compUUID,
		Score:       req.Score,
		Notes:       req.Notes,
		CreatedByID: teacherUUID,
	}

	// 4. Eksekusi penyimpanan
	return s.gradeRepo.UpsertGrade(grade)
}
