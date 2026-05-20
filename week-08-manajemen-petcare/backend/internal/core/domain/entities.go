package domain

import "time"

type Owner struct {
	ID    string
	Name  string
	Phone string
	Pets  []Pet
}

type Pet struct {
	ID        string
	OwnerID   string
	Name      string
	Species   string
	Breed     string
	Weight    float64
	DietNotes string
}

type NutritionLog struct {
	ID          string
	PetID       string
	LogDate     time.Time
	FoodBrand   string // Merek diet berkualitas tinggi
	Calories    int
	HealthNotes string // Transformasi fisik, bulu, pencernaan
}

type Invoice struct {
	ID          string
	OwnerID     string
	TotalAmount float64
	Status      string
	CreatedAt   time.Time
	Items       []InvoiceItem
}

type InvoiceItem struct {
	ID        string
	InvoiceID string
	PetID     string // Relasi absolut ke individu hewan
	ServiceID string
	Price     float64
}
