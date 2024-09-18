package dto

type LoanInterest struct {
	DurationMonths   int     `json:"duration mounths"`
	LoanInterestRate float64 `json:"interest rate"`
}
