package service

import (
	"github.com/affandisy/school-system-sma/models"
	"github.com/affandisy/school-system-sma/repository"
	"github.com/google/uuid"
)

type MasterService interface {
	CreateClass(req CreateClassRequest) error
	GetClasses() ([]models.Class, error)

	CreateSubject(req CreateSubjectRequest) error
	GetSubjects() ([]models.Subject, error)

	CreateStudent(req CreateStudentRequest) error
	GetStudents() ([]models.Student, error)
}

type masterService struct {
	repo repository.MasterRepository
}

func NewMasterService(repo repository.MasterRepository) MasterService {
	return &masterService{repo}
}

type CreateClassRequest struct {
	Name              string `json:"name"`
	StudyProgramID    string `json:"study_program_id"`
	AcademicYearID    string `json:"academic_year_id"`
	HomeroomTeacherID string `json:"homeroom_teacher_id"` // Opsional
}

type CreateSubjectRequest struct {
	Code           string `json:"code"`
	Name           string `json:"name"`
	StudyProgramID string `json:"study_program_id"` // Opsional (jika mapel umum)
}

type CreateStudentRequest struct {
	UserID           string `json:"user_id"` // Akun user yang sudah dibuat
	NIS              string `json:"nis"`
	ClassID          string `json:"class_id"`
	EnrollmentYearID string `json:"enrollment_year_id"`
}

func (s *masterService) CreateClass(req CreateClassRequest) error {
	spUUID, _ := uuid.Parse(req.StudyProgramID)
	ayUUID, _ := uuid.Parse(req.AcademicYearID)

	class := &models.Class{
		Name:           req.Name,
		StudyProgramID: spUUID,
		AcademicYearID: ayUUID,
	}

	if req.HomeroomTeacherID != "" {
		htUUID, _ := uuid.Parse(req.HomeroomTeacherID)
		class.HomeroomTeacherID = &htUUID
	}
	return s.repo.CreateClass(class)
}
func (s *masterService) GetClasses() ([]models.Class, error) { return s.repo.GetClasses() }

func (s *masterService) CreateSubject(req CreateSubjectRequest) error {
	subject := &models.Subject{Code: req.Code, Name: req.Name}
	if req.StudyProgramID != "" {
		spUUID, _ := uuid.Parse(req.StudyProgramID)
		subject.StudyProgramID = &spUUID
	}
	return s.repo.CreateSubject(subject)
}
func (s *masterService) GetSubjects() ([]models.Subject, error) { return s.repo.GetSubjects() }

func (s *masterService) CreateStudent(req CreateStudentRequest) error {
	userUUID, _ := uuid.Parse(req.UserID)
	classUUID, _ := uuid.Parse(req.ClassID)
	eyUUID, _ := uuid.Parse(req.EnrollmentYearID)

	student := &models.Student{
		UserID:           userUUID,
		NIS:              req.NIS,
		ClassID:          &classUUID,
		EnrollmentYearID: eyUUID,
	}
	return s.repo.CreateStudent(student)
}
func (s *masterService) GetStudents() ([]models.Student, error) { return s.repo.GetStudents() }
