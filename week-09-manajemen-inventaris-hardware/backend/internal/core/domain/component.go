package domain

// Component merepresentasikan inventory fisik dan detail spesifikasinya
type Component struct {
	ID              string                 `json:"id"`
	SKU             string                 `json:"sku"`
	Name            string                 `json:"name"`
	Category        string                 `json:"category"` // CPU, GPU, Motherboard, dll.
	Manufacturer    string                 `json:"manufacturer"`
	Model           string                 `json:"model"`
	BasePrice       float64                `json:"base_price"`
	StockOnHand     int                    `json:"stock_on_hand"`
	ReservedQty     int                    `json:"reserved_qty"`
	Location        string                 `json:"location"`
	IsSerialTracked bool                   `json:"is_serial_tracked"`
	Specs           map[string]interface{} `json:"specs"` // Atribut dinamis sesuai kategori
}

// Perilaku Domain: Pengecekan stok aman dengan mempertimbangkan reservasi
func (c *Component) AvailableStock() int {
	return c.StockOnHand - c.ReservedQty
}
