package repository

import (
	"github.com/affandisy/school-system-sma/models"
	"gorm.io/gorm"
)

type AnnouncementRepository interface {
	Create(ann *models.Announcement) error
	FindAll() ([]models.Announcement, error)
	SaveAuditLog(log *models.AuditLog) error
}

type announcementRepository struct {
	db *gorm.DB
}

func NewAnnouncementRepository(db *gorm.DB) AnnouncementRepository {
	return &announcementRepository{db}
}

func (r *announcementRepository) Create(ann *models.Announcement) error {
	return r.db.Create(ann).Error
}

func (r *announcementRepository) FindAll() ([]models.Announcement, error) {
	var anns []models.Announcement
	// Ambil pengumuman beserta nama pembuatnya, urutkan dari yang terbaru
	err := r.db.Preload("CreatedBy").Order("created_at DESC").Find(&anns).Error
	return anns, err
}

func (r *announcementRepository) SaveAuditLog(log *models.AuditLog) error {
	return r.db.Create(log).Error
}
