package usecasemock

import (
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type DepositeInterestUseCaseMock struct {
	mock.Mock
}

func (d *DepositeInterestUseCaseMock) Pagging(payload dto.PageRequest) ([]dto.DepositeInterestRequest, dto.Paging, error) {
	args := d.Called(payload)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]dto.DepositeInterestRequest), args.Get(1).(dto.Paging), nil
}

// DeleteById implements DepositeInterestUseCase.
func (d *DepositeInterestUseCaseMock) DeleteById(id string) error {
	return d.Called(id).Error(0)
}

// FindById implements DepositeInterestUseCase.
func (d *DepositeInterestUseCaseMock) FindById(id string) (dto.DepositeInterestRequest, error) {
	args := d.Called(id)
	if args.Get(1) != nil {
		return dto.DepositeInterestRequest{}, args.Error(1)
	}
	return args.Get(0).(dto.DepositeInterestRequest), nil
}

// Update implements DepositeInterestUseCase.
func (d *DepositeInterestUseCaseMock) Update(payload dto.DepositeInterestRequest) error {
	return d.Called(payload).Error(0)
}

// CreateNew implements DepositeInteresetUseCase.
func (d *DepositeInterestUseCaseMock) CreateNew(payload dto.DepositeInterestRequest) (int, error) {
	args := d.Called(payload)
	if args.Get(1) != nil {
		return args.Int(0), args.Error(1)
	}
	return args.Int(0), nil
}
