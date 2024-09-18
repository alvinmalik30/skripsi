package repomock

import (
	"polen/model"
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type LatePaymentRepoMock struct {
	mock.Mock
}

// CreateNew implements LatePaymentFee.
func (a *LatePaymentRepoMock) CreateNew(payload model.LatePaymentFee) error {
	return a.Called(payload).Error(0)
}

// DeleteById implements LatePaymentFee.
func (a *LatePaymentRepoMock) DeleteById(id string) error {
	return a.Called(id).Error(0)
}

// FindById implements LatePaymentFee.
func (a *LatePaymentRepoMock) FindById(id string) (model.LatePaymentFee, error) {
	args := a.Called(id)
	if args.Get(1) != nil {
		return model.LatePaymentFee{}, args.Error(1)
	}
	return args.Get(0).(model.LatePaymentFee), nil
}

// Pagging implements LatePaymentFee.
func (a *LatePaymentRepoMock) Pagging(payload dto.PageRequest) ([]model.LatePaymentFee, dto.Paging, error) {
	args := a.Called(payload)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.LatePaymentFee), args.Get(1).(dto.Paging), nil
}

// Update implements LatePaymentFee.
func (a *LatePaymentRepoMock) Update(payload model.LatePaymentFee) error {
	return a.Called(payload).Error(0)
}
