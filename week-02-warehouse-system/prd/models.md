type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`
	Email     string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"email"`
	Password  string         `gorm:"type:varchar(255);not null" json:"-"` // JSON "-" agar password tidak pernah bocor ke response
	AvatarURL *string        `gorm:"type:varchar(255)" json:"avatar_url"`
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Role      string         `gorm:"type:varchar(20);not null;default:'operator'" json:"role"` // admin, manager, operator
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type Product struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	SKU          string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"sku"`
	Name         string         `gorm:"type:varchar(150);not null" json:"name"`
	ImageURL     *string        `gorm:"type:varchar(255)" json:"image_url"`
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

type Transaction struct {
	ID                uint              `gorm:"primaryKey" json:"id"`
	ReferenceNo       string            `gorm:"type:varchar(50);uniqueIndex;not null" json:"reference_no"`
	TransactionDate   time.Time         `gorm:"not null" json:"transaction_date"`
	Type              string            `gorm:"type:varchar(20);not null" json:"type"`                   // INBOUND, OUTBOUND
	Status            string            `gorm:"type:varchar(20);not null;default:'draft'" json:"status"` // draft, approved, rejected
	PartnerID         *uint             `gorm:"not null" json:"partner_id"`
	Partner           Partner           `gorm:"foreignKey:PartnerID" json:"partner"`
	Notes             string            `gorm:"type:text" json:"notes"`
	WarehouseID       uint              `gorm:"not null" json:"warehouse_id"`
	Warehouse         Warehouse         `gorm:"foreignKey:WarehouseID" json:"warehouse"`
	TargetWarehouseID *uint             `json:"target_warehouse_id"`
	TargetWarehouse   Warehouse         `gorm:"foreignKey:TargetWarehouseID" json:"target_warehouse"`
	CreatedByID       uint              `gorm:"not null" json:"created_by_id"`
	ApprovedByID      *uint             `json:"approved_by_id"`
	Items             []TransactionItem `gorm:"foreignKey:TransactionID" json:"items"`
	CreatedAt         time.Time         `json:"created_at"`
	UpdatedAt         time.Time         `json:"updated_at"`
}

type TransactionItem struct {
	ID            uint    `gorm:"primaryKey" json:"id"`
	TransactionID uint    `gorm:"not null;index" json:"transaction_id"`
	ProductID     uint    `gorm:"not null;index" json:"product_id"`
	Product       Product `gorm:"foreignKey:ProductID" json:"product"` // Join untuk ambil info SKU/Nama
	Quantity      int     `gorm:"not null;default:1" json:"quantity"`
	UnitPrice     float64 `gorm:"type:decimal(15,2);not null" json:"unit_price"`
	SubTotal      float64 `gorm:"type:decimal(15,2);not null" json:"sub_total"`
}

type Location struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Code        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"code"` // Contoh: RAK-A1-01
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`            // Contoh: Rak Makanan Ringan
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

type Partner struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(150);not null" json:"name"`
	Type      string         `gorm:"type:varchar(20);not null" json:"type"` // "SUPPLIER" atau "CUSTOMER"
	Phone     string         `gorm:"type:varchar(20)" json:"phone"`
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

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

type Warehouse struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Code      string         `gorm:"type:varchar(20);uniqueIndex;not null" json:"code"` // Contoh: JKT-01
	Name      string         `gorm:"type:varchar(100);not null" json:"name"`            // Contoh: Gudang Utama Jakarta
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}