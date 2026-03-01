package models

import (
	"time"

	"gorm.io/gorm"
)

// 1. Model User Management
type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // JSON "-" agar password tidak pernah bocor ke response
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Role      string         `gorm:"type:varchar(20);not null;default:'operator'" json:"role"` // admin, manager, operator
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// 2. Model Master Barang (Inventory)
type Product struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SKU          string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"sku"`
	Name         string         `gorm:"type:varchar(150);not null" json:"name"`
	Description  string         `gorm:"type:text" json:"description"`
	Category     string         `gorm:"type:varchar(50)" json:"category"`
	Unit         string         `gorm:"type:varchar(20);not null" json:"unit"` // pcs, box, kg
	Price        float64        `gorm:"type:decimal(15,2);not null" json:"price"`
	MinStock     int            `gorm:"not null;default:0" json:"min_stock"`
	MaxStock     int            `gorm:"not null;default:0" json:"max_stock"`
	CurrentStock int            `gorm:"not null;default:0" json:"current_stock"`
	LocationID   *uint          `json:"location_id"`
	Location     Location       `gorm:"foreignKey:LocationID" json:"location"`
	IsActive     bool           `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"` // Fitur Soft Delete sesuai PRD
}

// 3. Model Transaksi Utama (Inbound/Outbound)
type Transaction struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ReferenceNo     string    `gorm:"type:varchar(50);uniqueIndex;not null" json:"reference_no"` // Nomor PO/Invoice
	TransactionDate time.Time `gorm:"not null" json:"transaction_date"`
	Type            string    `gorm:"type:varchar(20);not null" json:"type"`                   // INBOUND, OUTBOUND
	Status          string    `gorm:"type:varchar(20);not null;default:'draft'" json:"status"` // draft, approved, rejected
	// EntityName      string            `gorm:"type:varchar(100)" json:"entity_name"`                    // Nama Supplier / Customer
	PartnerID    uint              `gorm:"not null" json:"partner_id"`
	Partner      Partner           `gorm:"foreignKey:PartnerID" json:"partner"`
	Notes        string            `gorm:"type:text" json:"notes"`
	CreatedByID  uint              `gorm:"not null" json:"created_by_id"`         // Siapa yang input
	ApprovedByID *uint             `json:"approved_by_id"`                        // Siapa yang approve (bisa null)
	Items        []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"` // Relasi 1-to-Many ke detail barang
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
}

// 4. Model Detail Barang dalam Transaksi
type TransactionItem struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TransactionID uint    `gorm:"not null;index" json:"transaction_id"`
	ProductID     uint    `gorm:"not null;index" json:"product_id"`
	Product       Product `gorm:"foreignKey:ProductID" json:"product"` // Join untuk ambil info SKU/Nama
	Quantity      int     `gorm:"not null;default:1" json:"quantity"`
	UnitPrice     float64 `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	SubTotal      float64 `gorm:"type:decimal(15,2);not null" json:"sub_total"`
}
