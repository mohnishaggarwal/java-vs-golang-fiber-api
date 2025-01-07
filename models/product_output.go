package models

type ExpenseCategory string

const (
	VeryCheap ExpenseCategory = "VERY_CHEAP"
	Cheap     ExpenseCategory = "CHEAP"
	Expensive ExpenseCategory = "EXPENSIVE"
)

type ProductOutput struct {
	Product         *Product        `json:"product"`
	Url             string          `json:"url"`
	ExpenseCategory ExpenseCategory `json:"expenseCategory"`
}
