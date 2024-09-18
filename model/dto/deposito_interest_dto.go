package dto

type DepositeInterestRequest struct {
	Id             string  `json:"id"`
	InterestRate   float64 `json:"interest rate"`
	TaxRate        float64 `json:"tax rate"`
	DurationMounth int     `json:"duration mounth"`
}

type DepositeInterestReq struct {
	InterestRate   float64 `json:"interest rate"`
	TaxRate        float64 `json:"tax rate"`
	DurationMounth int     `json:"duration mounth"`
}
