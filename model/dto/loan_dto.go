package dto

import "time"

type LoanRequest struct {
	UserCredentialId      string
	LoanInterestRateId    string `json:"loan interest rate id"`
	LoanHandlingCostId    string `json:"loan handling cost id"`
	LoanLatePaymentFessId string `json:"loan late payment fess id"`
	LoanAmount            int    `json:"loan amount"`
}

type LoanReq struct {
	LoanInterestRateId    string `json:"loan interest rate id"`
	LoanHandlingCostId    string `json:"loan handling cost id"`
	LoanLatePaymentFessId string `json:"loan late payment fess id"`
	LoanAmount            int    `json:"loan amount"`
}

type Confirm struct {
	Id      string `json:"id"`
	IsPayed bool   `json:"is payed"`
	Status  string `json:"status"`
}

type LoanInstallenmentResponse struct {
	Id                    string      `json:"id"`
	IsPayed               bool        `json:"is payed"`
	PaymentInstallment    int         `json:"payment instalenment"`
	PaymentDeadLine       time.Time   `json:"payment deadline"`
	TotalAmountOfDepth    int         `json:"total amount"`
	LatePayment           LatePayment `json:"late payment"`
	PaymentDate           time.Time   `json:"payment date"`
	Status                string      `json:"status"`
	TransferConfirmRecipe bool        `json:"sending transfer recipe"`
	File                  string      `json:"file path"`
}

type LatePayment struct {
	LatePaymentFees      string `json:"late paymnet fees"`
	LatePaymentDays      int    `json:"late payment days"`
	LatePaymentFeesTotal int    `json:"late payment fees total"`
}

type InstallenmentLoanByIdResponse struct {
	UserDeatail BiodataResponse           `json:"user detail"`
	LoanId      string                    `json:"loan id"`
	LoanInst    LoanInstallenmentResponse `json:"loan installenment"`
}

type InstallenmentLoanByLoanIdResponse struct {
	UserDeatail BiodataResponse             `json:"user detail"`
	LoanId      string                      `json:"loan id"`
	LoanInst    []LoanInstallenmentResponse `json:"loan installenment"`
}

type Loan struct {
	Id                  string                      `json:"id"`
	UserCredentialId    string                      `json:"user credential id"`
	LoanDateCreate      time.Time                   `json:"loadn date create"`
	LoanAmount          int                         `json:"loan amount"`
	LoanDuration        int                         `json:"loan duration"`
	LoanInterestRate    float64                     `json:"loan interest rate"`
	LoanInterestNominal int                         `json:"loan interest nominal"`
	AppHandlingCost     string                      `json:"handling cost application"`
	TotalAmountOfDepth  int                         `json:"total amount"`
	IsPayed             bool                        `json:"is payed"`
	Status              string                      `json:"status"`
	Installment         []LoanInstallenmentResponse `json:"installment"`
}
