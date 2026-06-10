package engine

import (
	"github.com/affandisy/hardware-erp/internal/core/domain"
	"github.com/affandisy/hardware-erp/internal/core/port"
)

type defaultAssemblyEngine struct {
	rules []port.ValidationRule // Menampung array kontrak antarmuka, bukan implementasi konkret
}

// Konstruktor ini menerima variadic arguments, jadi Anda bisa memasukkan aturan sebanyak apa pun nanti
func NewAssemblyEngine(rules ...port.ValidationRule) port.AssemblyEngine {
	return &defaultAssemblyEngine{rules: rules}
}

func (e *defaultAssemblyEngine) ValidateBuild(components []domain.Component) (bool, []string, int, float64) {
	isCompatible := true
	var notes []string
	var totalPrice float64
	var totalPower float64

	// 1. Iterasi Dasar: Mengumpulkan Harga dan Estimasi Daya
	for _, c := range components {
		totalPrice += c.BasePrice
		if tdp, ok := c.Specs["tdp_w"].(float64); ok {
			totalPower += tdp
		}
		if pd, ok := c.Specs["power_draw_w"].(float64); ok {
			totalPower += pd
		}
	}

	// 2. Eksekusi Rule Engine (Dependency Inversion)
	for _, rule := range e.rules {
		valid, msg := rule.Evaluate(components)

		if !valid {
			isCompatible = false
			notes = append(notes, msg) // Kumpulkan alasan kenapa rakitan ini ditolak
		} else if msg != "" {
			notes = append(notes, msg) // Kumpulkan peringatan non-kritis
		}
	}

	if isCompatible && len(notes) == 0 {
		notes = append(notes, "Sistem tervalidasi. Semua komponen kompatibel dan siap untuk dirakit.")
	}

	return isCompatible, notes, int(totalPower), totalPrice
}
