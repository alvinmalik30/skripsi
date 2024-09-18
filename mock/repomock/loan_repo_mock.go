package repomock

import (
	"polen/model"
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type LoanRepoMock struct {
	mock.Mock
}

// FindByLoanId implements LoanRepoMock.
func (l *LoanRepoMock) FindByLoanId(id string) (dto.Loan, error) {
	a := l.Called(id)
	if a.Get(1) != nil {
		return dto.Loan{}, a.Error(1)
	}
	return a.Get(0).(dto.Loan), nil
}

// UpdateLateFee implements LoanRepoMock.
func (l *LoanRepoMock) UpdateLateFee() error {
	return l.Called().Error(0)
}

// Accepted implements LoanRepoMock.
func (l *LoanRepoMock) Accepted(payload dto.InstallenmentLoanByIdResponse) error {
	return l.Called(payload).Error(0)
}

// Upload implements LoanRepoMock.
func (l *LoanRepoMock) Upload(payload dto.LoanInstallenmentResponse) error {
	return l.Called(payload).Error(0)
}

// FindById implements LoanRepoMock.
func (l *LoanRepoMock) FindById(id string) (dto.InstallenmentLoanByIdResponse, error) {
	a := l.Called(id)
	if a.Get(1) != nil {
		return dto.InstallenmentLoanByIdResponse{}, a.Error(1)
	}
	return a.Get(0).(dto.InstallenmentLoanByIdResponse), nil
}

// FindById implements LoanRepoMock.
func (l *LoanRepoMock) FindUploadedFile() ([]dto.LoanInstallenmentResponse, error) {
	a := l.Called()
	if a.Get(1) != nil {
		return nil, a.Error(1)
	}
	return a.Get(0).([]dto.LoanInstallenmentResponse), nil
}

// create implements LoanRepoMock.
func (l *LoanRepoMock) Create(loan model.Loan, installenment []model.InstallenmentLoan) error {
	return l.Called(loan, installenment).Error(0)
}
