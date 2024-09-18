package model

import "time"

type DepositeInterest struct {
	Id             string
	CreateDate     time.Time
	InterestRate   float64
	TaxRate        float64
	DurationMounth int
}

// grossProfit
// (setoran pokok*interest rate*durationDay)/365

// tax
// TaxRate*grossProfit

// netProfit
// grossProfit-tax

// totalReturn
// setoran+netProfit
