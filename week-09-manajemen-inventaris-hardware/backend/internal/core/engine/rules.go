package engine

import (
	"fmt"

	"github.com/affandisy/hardware-erp/internal/core/domain"
	"github.com/affandisy/hardware-erp/internal/core/port"
)

// --- 1. ATURAN SOKET (Memastikan prosesor kompatibel dengan motherboard) ---
type socketRule struct{}

func NewSocketRule() port.ValidationRule { return &socketRule{} }

func (r *socketRule) Evaluate(components []domain.Component) (bool, string) {
	var cpuSocket, moboSocket string
	var hasCPU, hasMobo bool

	for _, c := range components {
		if c.Category == "CPU" {
			if s, ok := c.Specs["socket"].(string); ok {
				cpuSocket = s
				hasCPU = true
			}
		} else if c.Category == "Motherboard" {
			if s, ok := c.Specs["socket"].(string); ok {
				moboSocket = s
				hasMobo = true
			}
		}
	}

	// Jika pengguna memasangkan Ryzen 5600 (AM4) ke papan Intel (LGA1700), sistem menolak keras
	if hasCPU && hasMobo && cpuSocket != moboSocket {
		return false, fmt.Sprintf("Bentrok Soket Kritis: CPU membutuhkan soket %s, sedangkan Motherboard menyediakan %s", cpuSocket, moboSocket)
	}
	return true, ""
}

// --- 2. ATURAN DAYA (Memastikan kapasitas catu daya memadai) ---
type powerRule struct{}

func NewPowerRule() port.ValidationRule { return &powerRule{} }

func (r *powerRule) Evaluate(components []domain.Component) (bool, string) {
	var totalPower float64
	var psuWattage float64
	var hasPSU bool

	for _, c := range components {
		// Mengakumulasi daya prosesor
		if tdp, ok := c.Specs["tdp_w"].(float64); ok {
			totalPower += tdp
		}
		// Mengakumulasi daya kartu grafis
		if pd, ok := c.Specs["power_draw_w"].(float64); ok {
			totalPower += pd
		}

		// Mencari kapasitas Power Supply
		if c.Category == "PSU" {
			if w, ok := c.Specs["wattage"].(float64); ok {
				psuWattage = w
				hasPSU = true
			}
		}
	}

	// Misalnya RTX 4060 dan komponen lain butuh 350W, tapi PSU hanya 300W
	if hasPSU && totalPower > psuWattage {
		return false, fmt.Sprintf("Kekurangan Daya: Total kebutuhan komponen %.0fW melebihi kapasitas aman PSU (%.0fW)", totalPower, psuWattage)
	}
	if !hasPSU && totalPower > 0 {
		// Mengembalikan 'true' (kompatibel) tetapi memberikan peringatan
		return true, "Peringatan: Anda belum memilih Power Supply untuk menyuplai komponen berdaya tinggi ini."
	}

	return true, ""
}
