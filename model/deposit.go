package model

import "time"

type Deposite struct {
	Id               string
	UserCredentialId string
	DepositeAmount   int
	InterestRate     float64
	TaxRate          float64
	DurationMounth   int
	MaturityDate     time.Time
	Status           string
	GrossProfit      int
	Tax              int
	NetProfil        int
	TotalReturn      int
}
