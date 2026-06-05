package domain

type PetNutritionSummary struct {
	PetID            string
	PetName          string
	TotalLogs        int
	AverageCalories  float64
	MostUsedBrand    string
	LatestHealthNote string
}
