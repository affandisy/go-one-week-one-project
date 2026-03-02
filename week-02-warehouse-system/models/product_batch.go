package models

import "time"

type ProductBatch struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	ProductID  uint      `gorm:"not null;index" json:"product_id"`
	Product    Product   `gorm:"foreignKey:ProductID" json:"-"`
	BatchNo    string    `gorm:"type:varchar(50);not null;index" json:"batch_no"`
	ExpiryDate time.Time `gorm:"type:date" json:"expiry_date"`
	Stock      int       `gorm:"not null;default:0" json:"stock"` // Stok khusus untuk batch ini
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
