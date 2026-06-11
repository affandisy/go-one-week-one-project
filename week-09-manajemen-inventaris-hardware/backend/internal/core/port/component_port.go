package port

import "github.com/affandisy/hardware-erp/internal/core/domain"

// Driven Port (Database)
type ComponentRepository interface {
	Save(component *domain.Component) error
	FindBySKU(sku string) (*domain.Component, error)
	ReserveAndCheckout(skus []string) error
}

// Driving Port (API)
type ComponentUseCase interface {
	// Menambahkan specs sebagai map[string]interface{} agar dinamis
	RegisterComponent(sku, name, category, manufacturer string, price float64, specs map[string]interface{}) (*domain.Component, error)
}
