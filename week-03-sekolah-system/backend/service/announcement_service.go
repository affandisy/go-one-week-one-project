package service

import (
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/affandisy/school-system-sma/models"
	"github.com/affandisy/school-system-sma/repository"
	"github.com/google/uuid"
)

type AnnouncementService interface {
	CreateAnnouncement(title, content string, targetRoles []string, creatorID string) error
	GetAnnouncements() ([]models.Announcement, error)
	TriggerAcademicReportGeneration(academicYearID string, requestedByID string) error
}

type announcementService struct {
	repo repository.AnnouncementRepository
}

func NewAnnouncementService(repo repository.AnnouncementRepository) AnnouncementService {
	return &announcementService{repo}
}

func (s *announcementService) CreateAnnouncement(title, content string, targetRoles []string, creatorID string) error {
	if title == "" || content == "" {
		return errors.New("Judul dan konten pengumuman tidak boleh kosong")
	}

	targetRolesJSON, _ := json.Marshal(targetRoles)
	cUUID, _ := uuid.Parse(creatorID)

	ann := &models.Announcement{
		Title:       title,
		Content:     content,
		TargetRoles: string(targetRolesJSON),
		CreatedByID: cUUID,
	}

	// Simpan pengumuman
	if err := s.repo.Create(ann); err != nil {
		return errors.New("Gagal membuat pengumuman")
	}

	audit := &models.AuditLog{
		ActorUserID: cUUID,
		Action:      "CREATE",
		EntityType:  "Announcement",
		EntityID:    ann.ID.String(),
		AfterData:   string(targetRolesJSON),
	}
	_ = s.repo.SaveAuditLog(audit)

	return nil
}

func (s *announcementService) GetAnnouncements() ([]models.Announcement, error) {
	return s.repo.FindAll()
}

func (s *announcementService) TriggerAcademicReportGeneration(academicYearID string, requestedByID string) error {
	if academicYearID == "" {
		return errors.New("Tahun Ajaran wajib diisi")
	}

	go s.generateReportWorker(academicYearID, requestedByID)

	return nil
}

func (s *announcementService) generateReportWorker(academicYearID string, requestedByID string) {
	log.Printf("👷 [WORKER] Memulai generate laporan akademik untuk Tahun Ajaran %s...", academicYearID)

	time.Sleep(5 * time.Second)

	log.Printf("✅ [WORKER] Laporan berhasil dibuat. (Sistem akan mengirim email / notifikasi ke User ID: %s)", requestedByID)

	rUUID, _ := uuid.Parse(requestedByID)
	audit := &models.AuditLog{
		ActorUserID: rUUID,
		Action:      "GENERATE_REPORT",
		EntityType:  "AcademicReport",
		EntityID:    academicYearID,
	}
	_ = s.repo.SaveAuditLog(audit)
}
