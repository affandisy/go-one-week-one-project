package models

import "time"

type StockAdjustment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"not null;index" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"` // Relasi ke data barang
	UserID    uint      `gorm:"not null" json:"user_id"`             // Siapa yang melakukan penyesuaian
	Qty       int       `gorm:"not null" json:"qty"`                 // Bisa Plus (contoh: 5) atau Minus (contoh: -2)
	Reason    string    `gorm:"type:text;not null" json:"reason"`    // Alasan: "Barang rusak dimakan tikus", dll
	Status    string    `gorm:"type:varchar(20);default:'approved'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
