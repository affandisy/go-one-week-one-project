package domain

// AssemblyBuild merepresentasikan satu kesatuan rakitan PC
type AssemblyBuild struct {
	ID              string      `json:"id"`
	BuildName       string      `json:"build_name"`
	Components      []Component `json:"components"` // Komponen fisik yang dipilih
	TotalPowerW     int         `json:"total_power_w"`
	TotalPrice      float64     `json:"total_price"`
	IsCompatible    bool        `json:"is_compatible"`
	ValidationNotes []string    `json:"validation_notes"`
	CreatedBy       string      `json:"created_by"`
}
