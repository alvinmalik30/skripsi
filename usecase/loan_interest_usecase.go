package usecase

import (
	"fmt"
	"polen/model"
	"polen/model/dto"
	"polen/repository"
)

type LoanInterestUseCase interface {
	CreateNew(payload model.LoanInterest) (int, error)
	Pagging(payload dto.PageRequest) ([]model.LoanInterest, dto.Paging, error)
	FindById(id string) (model.LoanInterest, error)
	Update(payload model.LoanInterest) error
	DeleteById(id string) error
}

type loanInterestUseCase struct {
	repo repository.LoanInterest
}

// CreateNew implements LoanInterestUseCase.
func (l *loanInterestUseCase) CreateNew(payload model.LoanInterest) (int, error) {
	if payload.Id == "" {
		return 400, fmt.Errorf("id is required")
	}
	if payload.DurationMonths == 0 {
		return 400, fmt.Errorf("duration month is required")
	}
	if payload.LoanInterestRate == 0 {
		return 400, fmt.Errorf("loan interest rate is required")
	}
	if err := l.repo.CreateNew(payload); err != nil {
		return 500, err
	}
	return 201, nil
}

// DeleteById implements LoanInterestUseCase.
func (l *loanInterestUseCase) DeleteById(id string) error {
	loan, err := l.repo.FindById(id)
	if err != nil {
		return err
	}

	err = l.repo.DeleteById(loan.Id)
	if err != nil {
		return fmt.Errorf("failed to delete loan: %v", err)
	}

	return nil
}

// FindById implements LoanInterestUseCase.
func (l *loanInterestUseCase) FindById(id string) (model.LoanInterest, error) {
	return l.repo.FindById(id)
}

// Pagging implements LoanInterestUseCase.
func (l *loanInterestUseCase) Pagging(payload dto.PageRequest) ([]model.LoanInterest, dto.Paging, error) {
	return l.repo.Pagging(payload)
}

// Update implements LoanInterestUseCase.
func (l *loanInterestUseCase) Update(payload model.LoanInterest) error {
	if payload.Id == "" {
		return fmt.Errorf("id is required")
	}

	result, err := l.FindById(payload.Id)
	if err != nil {
		return err
	}

	if payload.DurationMonths == 0 {
		payload.DurationMonths = result.DurationMonths
	}
	if payload.LoanInterestRate == 0 {
		payload.LoanInterestRate = result.LoanInterestRate
	}

	err = l.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update loan interest: %v", err)
	}

	return nil
}

func NewLoanInterestUseCase(repo repository.LoanInterest) LoanInterestUseCase {
	return &loanInterestUseCase{
		repo: repo,
	}
}
