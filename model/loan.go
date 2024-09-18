package model

import "time"

type Loan struct {
	Id                     string
	UserCredentialId       string
	LoanDateCreate         time.Time
	LoanAmount             int
	LoanDuration           int
	LoanInterestRate       float64
	LoanInterestNominal    int
	AppHandlingCostNominal float64
	AppHandlingCostUnit    string
	TotalAmountOfDepth     int
	IsPayed                bool
	Status                 string
}

type InstallenmentLoan struct {
	Id                     string
	LoanId                 string
	IsPayed                bool
	PaymentInstallment     int
	PaymentDeadLine        time.Time
	TotalAmountOfDepth     int
	LatePaymentFeesNominal float64
	LatePaymentFeesUnit    string
	LatePaymentDays        int
	LatePaymentFeesTotal   int
	PaymentDate            time.Time
	Status                 string
	TransferConfirmRecipe  bool
	File                   string
}
