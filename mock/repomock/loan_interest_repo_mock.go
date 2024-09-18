package repomock

import (
	"polen/model"
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type LoanInterestRepoMock struct {
	mock.Mock
}

// Pagging implements LoanInterest.
func (l *LoanInterestRepoMock) Pagging(payload dto.PageRequest) ([]model.LoanInterest, dto.Paging, error) {
	args := l.Called(payload)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.LoanInterest), args.Get(1).(dto.Paging), nil
}

// CreateNew implements LoanInterest.
func (l *LoanInterestRepoMock) CreateNew(payload model.LoanInterest) error {
	return l.Called(payload).Error(0)
}

// DeleteById implements LoanInterest.
func (l *LoanInterestRepoMock) DeleteById(id string) error {
	return l.Called(id).Error(0)
}

// FindById implements LoanInterest.
func (l *LoanInterestRepoMock) FindById(id string) (model.LoanInterest, error) {
	args := l.Called(id)
	if args.Get(1) != nil {
		return model.LoanInterest{}, args.Error(1)
	}
	return args.Get(0).(model.LoanInterest), nil
}

// Update implements LoanInterest.
func (l *LoanInterestRepoMock) Update(payload model.LoanInterest) error {
	return l.Called(payload).Error(0)
}
