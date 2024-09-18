package dto

import "time"

type DepositeDto struct {
	Id             string                  `json:"id"`
	UserCredential GetAuthResponse         `json:"user credential"`
	InterestRate   DepositeInterestRequest `json:"interest rate"`
	DepositeAmount int
	MaturityDate   time.Time
	Status         bool
	GrossProfit    int
	Tax            int
	NetProfit      int
	TotalReturn    int
}

type DepositeRequest struct {
	UcId           string `json:"user id"`
	Amount         int    `json:"deposite amount"`
	InterestRateId string `json:"interest rate id"`
}

type Deposite struct {
	Id             string    `json:"id"`
	DepositeAmount int       `json:"deposite amount"`
	InterestRate   float64   `json:"interest rate"`
	TaxRate        float64   `json:"tax rate"`
	DurationMounth int       `json:"duration mounth"`
	CreateDate     time.Time `json:"create date"`
	MaturityDate   time.Time `json:"return date"`
	Status         string    `json:"is active"`
	GrossProfit    int       `json:"gross profit"`
	Tax            int       `json:"taxes"`
	NetProfit      int       `json:"net profit"`
	TotalReturn    int       `json:"total return"`
}

type DepositeByUserResponse struct {
	BioUser  BiodataResponse `json:"user biodata"`
	Deposite []Deposite      `json:"deposite"`
}

type DepositeByIdResponse struct {
	BioUser  BiodataResponse `json:"user biodata"`
	Deposite Deposite        `json:"deposite"`
}
