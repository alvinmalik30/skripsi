package model

type LoanInterest struct {
	Id               string  `json:"id"`
	DurationMonths   int     `json:"duration mounths"`
	LoanInterestRate float64 `json:"interest rate"`
}
