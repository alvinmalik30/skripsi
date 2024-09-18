package usecase

import (
	"fmt"
	"polen/model"
	"polen/model/dto"
	"polen/repository"
	"polen/utils/common"
	"time"
)

type LoanUseCase interface {
	Create(payload dto.LoanRequest) (int, error)
	FindById(id string) (dto.InstallenmentLoanByIdResponse, error)
	UploadFile(iduc string, payload dto.LoanInstallenmentResponse) (int, error)
	FindUploadedFile() ([]dto.LoanInstallenmentResponse, error)
	Accepted(payload dto.InstallenmentLoanByIdResponse) error
	UpdateLateFee() error
	FindByLoanId(id string) (dto.Loan, error)
}

type loanUseCase struct {
	repo    repository.LoanRepository
	loanIr  LoanInterestUseCase
	loanHc  AppHandlingCostUsecase
	loanLPF LatePaymentFeeUsecase
}

// FindByLoanId implements LoanUseCase.
func (l *loanUseCase) FindByLoanId(id string) (dto.Loan, error) {
	return l.repo.FindByLoanId(id)
}

// Accepted implements LoanUseCase.
func (l *loanUseCase) Accepted(payload dto.InstallenmentLoanByIdResponse) error {
	return l.repo.Accepted(payload)
}

// UpdateLateFee implements LoanUseCase.
func (l *loanUseCase) UpdateLateFee() error {
	return l.repo.UpdateLateFee()
}

// FindUploadedFile implements LoanUseCase.
func (l *loanUseCase) FindUploadedFile() ([]dto.LoanInstallenmentResponse, error) {
	return l.repo.FindUploadedFile()
}

// FindById implements LoanUseCase.
func (l *loanUseCase) FindById(id string) (dto.InstallenmentLoanByIdResponse, error) {
	return l.repo.FindById(id)
}

// Create implements LoanUseCase.
func (loan *loanUseCase) Create(payload dto.LoanRequest) (int, error) {
	if payload.LoanInterestRateId == "" {
		return 400, fmt.Errorf("loan interest rate is required")
	}
	if payload.LoanHandlingCostId == "" {
		return 400, fmt.Errorf("loan handling cost is required")
	}
	if payload.LoanLatePaymentFessId == "" {
		return 400, fmt.Errorf("loan late payment fees cost is required")
	}
	if payload.LoanAmount <= 0 {
		return 400, fmt.Errorf("loan amount cost is must greather than zero")
	}
	loanIr, err := loan.loanIr.FindById(payload.LoanInterestRateId)
	if err != nil {
		return 404, fmt.Errorf("loan duration you choose arent available,")
	}
	loanHc, err := loan.loanHc.FindById(payload.LoanHandlingCostId)
	if err != nil {
		return 404, fmt.Errorf("handling cost you choose arent available,")
	}
	loanLpf, err := loan.loanLPF.FindById(payload.LoanLatePaymentFessId)
	if err != nil {
		return 404, fmt.Errorf("loan ate payment fees you choose arent available,")
	}
	// building payload
	var loanpayload model.Loan
	loanpayload.Id = common.GenerateID()
	loanpayload.UserCredentialId = payload.UserCredentialId
	loanpayload.LoanDateCreate = time.Now()
	loanpayload.LoanDuration = loanIr.DurationMonths
	loanpayload.LoanAmount = payload.LoanAmount
	loanpayload.LoanInterestRate = loanIr.LoanInterestRate
	loanpayload.AppHandlingCostUnit = loanHc.Unit
	// app handling cost
	if loanpayload.AppHandlingCostUnit == "rupiah" {
		loanpayload.AppHandlingCostNominal = loanHc.Nominal
	} else if loanpayload.AppHandlingCostUnit == "percent" {
		loanpayload.AppHandlingCostNominal = float64(payload.LoanAmount) * loanHc.Nominal
	}
	// loan rate
	loanpayload.LoanInterestNominal = int(float64(payload.LoanAmount) * loanIr.LoanInterestRate)
	loanpayload.TotalAmountOfDepth = int(float64(loanpayload.LoanAmount) + float64(loanpayload.LoanInterestNominal) + loanpayload.AppHandlingCostNominal)
	loanpayload.Status = "active loan"
	loanpayload.IsPayed = false

	var instalenmentpayload []model.InstallenmentLoan
	type paymentInstallment struct {
		paymentDeadLine time.Time
		total           int
	}
	var pis []paymentInstallment

	for i := 1; i <= loanIr.DurationMonths; i++ {
		pi := paymentInstallment{
			paymentDeadLine: time.Now().AddDate(0, i, 0),
			total:           loanpayload.TotalAmountOfDepth / loanIr.DurationMonths,
		}
		pis = append(pis, pi)
	}

	for _, v := range pis {
		// build payload installentment
		var instalenment model.InstallenmentLoan
		instalenment.Id = common.GenerateID()
		instalenment.LoanId = loanpayload.Id
		instalenment.IsPayed = false
		instalenment.PaymentInstallment = v.total
		instalenment.PaymentDeadLine = v.paymentDeadLine
		instalenment.TotalAmountOfDepth = v.total
		instalenment.LatePaymentFeesUnit = loanLpf.Unit
		instalenment.LatePaymentFeesNominal = loanLpf.Nominal
		instalenment.LatePaymentFeesTotal = 0
		instalenment.LatePaymentDays = 0
		instalenment.Status = "waiting for payment"
		instalenment.TransferConfirmRecipe = false
		instalenment.File = ""
		instalenmentpayload = append(instalenmentpayload, instalenment)
	}

	// save
	err = loan.repo.Create(loanpayload, instalenmentpayload)
	if err != nil {
		return 500, err
	}

	return 200, nil
}

func (t *loanUseCase) UploadFile(iduc string, payload dto.LoanInstallenmentResponse) (int, error) {
	// check is data available
	dataTopUp, err := t.FindById(payload.Id)
	if err != nil {
		return 500, err
	}
	if dataTopUp.UserDeatail.UserCredential.Id != iduc {
		return 403, fmt.Errorf("you arent allowed")
	}
	dataTopUp.LoanInst.File = payload.File

	err = t.repo.Upload(payload)
	if err != nil {
		return 500, err
	}
	return 200, err
}

func NewLoanUseCase(repo repository.LoanRepository, loanir LoanInterestUseCase, loanhc AppHandlingCostUsecase, loanLPF LatePaymentFeeUsecase) LoanUseCase {
	return &loanUseCase{
		repo:    repo,
		loanIr:  loanir,
		loanHc:  loanhc,
		loanLPF: loanLPF,
	}
}
