package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	PINHash   string    `gorm:"not null" json:"-"`
	Role      string    `gorm:"not null" json:"role"` // cashier, owner, admin
	CreatedAt time.Time `json:"created_at"`
}

type Supplier struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Phone     string    `json:"phone"`
	Notes     string    `json:"notes"`
	CreatedAt time.Time `json:"created_at"`
}

type Product struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	Barcode       string    `gorm:"unique;index" json:"barcode"`
	Name          string    `gorm:"not null;index" json:"name"`
	PurchasePrice float64   `gorm:"not null" json:"purchase_price"`
	SellingPrice  float64   `gorm:"not null" json:"selling_price"`
	Unit          string    `gorm:"not null" json:"unit"` // pcs, kg, dus
	Stock         int       `gorm:"not null;default:0" json:"stock"`
	MinStock      int       `gorm:"not null;default:5" json:"min_stock"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type Transaction struct {
	ID            uuid.UUID           `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ReceiptNumber string              `gorm:"unique;not null" json:"receipt_number"`
	CashierID     uuid.UUID           `gorm:"type:uuid;not null" json:"cashier_id"`
	TotalAmount   float64             `gorm:"not null" json:"total_amount"`
	Discount      float64             `gorm:"default:0" json:"discount"`
	FinalAmount   float64             `gorm:"not null" json:"final_amount"`
	PaymentMethod string              `gorm:"not null" json:"payment_method"` // cash, qris
	CashGiven     float64             `json:"cash_given"`
	ChangeAmount  float64             `json:"change_amount"`
	CreatedAt     time.Time           `json:"created_at"`
	Details       []TransactionDetail `gorm:"foreignKey:TransactionID" json:"details"`
}

type TransactionDetail struct {
	ID            uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	TransactionID uuid.UUID `gorm:"type:uuid;not null;index" json:"transaction_id"`
	ProductID     uuid.UUID `gorm:"type:uuid;not null" json:"product_id"`
	Quantity      int       `gorm:"not null" json:"quantity"`
	Price         float64   `gorm:"not null" json:"price"` // Harga jual saat transaksi terjadi
	Subtotal      float64   `gorm:"not null" json:"subtotal"`
}

type StockMovement struct {
	ID         uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey" json:"id"`
	ProductID  uuid.UUID  `gorm:"type:uuid;not null;index" json:"product_id"`
	SupplierID *uuid.UUID `gorm:"type:uuid" json:"supplier_id"` // Opsional saat barang masuk
	Type       string     `gorm:"not null" json:"type"`         // IN (Barang Masuk), OUT (Penjualan/Rusak)
	Quantity   int        `gorm:"not null" json:"quantity"`
	Note       string     `json:"note"`
	CreatedAt  time.Time  `json:"created_at"`
}
