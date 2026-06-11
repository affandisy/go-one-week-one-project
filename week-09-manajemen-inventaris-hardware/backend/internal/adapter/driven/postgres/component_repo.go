package postgres

import (
	"encoding/json"
	"errors"

	"github.com/affandisy/hardware-erp/internal/core/domain"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (r *componentRepository) FindBySKU(sku string) (*domain.Component, error) {
	var model ComponentModel
	if err := r.db.Where("sku = ?", sku).First(&model).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	var specs map[string]interface{}
	if len(model.Specs) > 0 {
		if err := json.Unmarshal(model.Specs, &specs); err != nil {
			return nil, err
		}
	}
	if specs == nil {
		specs = map[string]interface{}{}
	}

	component := &domain.Component{
		ID:              model.ID,
		SKU:             model.SKU,
		Name:            model.Name,
		Category:        model.Category,
		Manufacturer:    model.Manufacturer,
		Model:           model.Model,
		BasePrice:       model.BasePrice,
		StockOnHand:     model.StockOnHand,
		ReservedQty:     model.ReservedQty,
		Location:        model.Location,
		IsSerialTracked: model.IsSerialTracked,
		Specs:           specs,
	}

	return component, nil
}

func (r *componentRepository) ReserveAndCheckout(skus []string) error {
	// Memulai transaksi database agar bersifat Atomik (ACID)
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, sku := range skus {
			var comp ComponentModel

			// Row Locking: Clause .Clauses(clause.Locking{Strength: "UPDATE"})
			// mencegah transaksi lain membaca/mengubah baris ini sampai transaksi ini selesai
			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("sku = ?", sku).First(&comp).Error; err != nil {
				return err // Rollback jika komponen tidak ditemukan
			}

			// Validasi Stok
			if comp.StockOnHand <= 0 {
				return errors.New("stok habis untuk komponen: " + sku)
			}

			// Kurangi stok
			comp.StockOnHand -= 1

			if err := tx.Save(&comp).Error; err != nil {
				return err // Rollback jika gagal update
			}
		}
		return nil // Commit transaksi jika semua sukses
	})
}
