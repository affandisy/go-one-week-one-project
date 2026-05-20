package postgres

import (
	"time"

	"github.com/affandisy/petcare-system/internal/core/domain"
	"github.com/affandisy/petcare-system/internal/core/port"
	"gorm.io/gorm"
)

// DTO khusus GORM untuk skema Database (FR-004)
type InvoiceModel struct {
	ID          string `gorm:"type:uuid;primaryKey"`
	OwnerID     string `gorm:"type:uuid;index"`
	TotalAmount float64
	Status      string
	CreatedAt   time.Time
	Items       []InvoiceItemModel `gorm:"foreignKey:InvoiceID"`
}

func (InvoiceModel) TableName() string { return "invoices" }

type InvoiceItemModel struct {
	ID        string `gorm:"type:uuid;primaryKey"`
	InvoiceID string `gorm:"type:uuid;index"`
	PetID     string `gorm:"type:uuid;index"` // Memisahkan data transaksi ke masing-masing hewan
	ServiceID string `gorm:"type:uuid"`
	Price     float64
}

func (InvoiceItemModel) TableName() string { return "invoice_items" }

type billingRepository struct {
	db *gorm.DB
}

func NewBillingRepository(db *gorm.DB) port.BillingRepository {
	db.AutoMigrate(&InvoiceModel{}, &InvoiceItemModel{}) // Migrasi otomatis
	return &billingRepository{db}
}

func (r *billingRepository) SaveInvoice(invoice *domain.Invoice) error {
	// Mapper Domain -> Model DB
	var itemsModel []InvoiceItemModel
	for _, item := range invoice.Items {
		itemsModel = append(itemsModel, InvoiceItemModel{
			ID:        item.ID,
			InvoiceID: item.InvoiceID,
			PetID:     item.PetID,
			ServiceID: item.ServiceID,
			Price:     item.Price,
		})
	}

	model := InvoiceModel{
		ID:          invoice.ID,
		OwnerID:     invoice.OwnerID,
		TotalAmount: invoice.TotalAmount,
		Status:      invoice.Status,
		CreatedAt:   invoice.CreatedAt,
		Items:       itemsModel,
	}

	return r.db.Create(&model).Error
}
