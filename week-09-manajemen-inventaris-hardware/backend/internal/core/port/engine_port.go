package port

import "github.com/affandisy/hardware-erp/internal/core/domain"

// ValidationRule adalah kontrak untuk setiap aturan pengecekan hardware
type ValidationRule interface {
	// Mengembalikan boolean kompatibilitas dan pesan error jika ada pelanggaran
	Evaluate(components []domain.Component) (bool, string)
}

// AssemblyEngine bertanggung jawab merangkai komponen dan menjalankan semua ValidationRule
type AssemblyEngine interface {
	ValidateBuild(components []domain.Component) (bool, []string, int, float64)
}
