package domain

type CategorySummary struct {
	CategoryName string
	Color        string
	Amount       float64
	Percentage   float64
}

type MonthlyReport struct {
	TotalIncome        float64
	TotalExpense       float64
	NetBalance         float64
	VsLastMonthPercent float64
	ExpenseBreakdown   []CategorySummary
}
