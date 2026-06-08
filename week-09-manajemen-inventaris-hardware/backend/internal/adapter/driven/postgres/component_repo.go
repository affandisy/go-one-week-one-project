package postgres

import (
	"encoding/json"

	"github.com/affandisy/hardware-erp/internal/core/domain"
	"gorm.io/gorm"
)

// ComponentModel adalah DTO khusus database untuk menangani tipe data JSONB
type ComponentModel struct {
	ID              string  `gorm:"type:uuid;primaryKey"`
	SKU             string  `gorm:"type:varchar(100);uniqueIndex"`
	Name            string  `gorm:"type:varchar(255);not null"`
	Category        string  `gorm:"type:varchar(50);index"`
	Manufacturer    string  `gorm:"type:varchar(100)"`
	Model           string  `gorm:"type:varchar(100)"`
	BasePrice       float64 `gorm:"type:decimal(15,2)"`
	StockOnHand     int     `gorm:"type:int;default:0"`
	ReservedQty     int     `gorm:"type:int;default:0"`
	Location        string  `gorm:"type:varchar(100)"`
	IsSerialTracked bool    `gorm:"type:boolean;default:false"`
	Specs           []byte  `gorm:"type:jsonb"` // GORM menyimpan JSONB sebagai byte array
}

func (ComponentModel) TableName() string { return "components" }

type componentRepository struct {
	db *gorm.DB
}

func NewComponentRepository(db *gorm.DB) *componentRepository {
	db.AutoMigrate(&ComponentModel{})
	return &componentRepository{db}
}

func (r *componentRepository) Save(c *domain.Component) error {
	specsBytes, err := json.Marshal(c.Specs)
	if err != nil {
		return err
	}

	model := ComponentModel{
		ID:              c.ID,
		SKU:             c.SKU,
		Name:            c.Name,
		Category:        c.Category,
		Manufacturer:    c.Manufacturer,
		Model:           c.Model,
		BasePrice:       c.BasePrice,
		StockOnHand:     c.StockOnHand,
		ReservedQty:     c.ReservedQty,
		Location:        c.Location,
		IsSerialTracked: c.IsSerialTracked,
		Specs:           specsBytes,
	}

	return r.db.Save(&model).Error
}
