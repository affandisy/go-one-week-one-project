package models

import "time"

type WarehouseStock struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	WarehouseID uint      `gorm:"not null;uniqueIndex:idx_warehouse_product" json:"warehouse_id"`
	ProductID   uint      `gorm:"not null;uniqueIndex:idx_warehouse_product" json:"product_id"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	UpdatedAt   time.Time `json:"updated_at"`

	// Relasi
	Warehouse Warehouse `gorm:"foreignKey:WarehouseID" json:"warehouse"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product"`
}
