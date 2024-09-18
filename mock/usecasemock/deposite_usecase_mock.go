package usecasemock

import (
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type DepositeUseCaseMock struct {
	mock.Mock
}

// Update implements DepositeUseCase.
func (d *DepositeUseCaseMock) Update() error {
	return d.Called().Error(0)
}

// Pagging implements DepositeUseCaseMock.
func (d *DepositeUseCaseMock) Pagging(payload dto.PageRequest) ([]dto.Deposite, dto.Paging, error) {
	args := d.Called(payload)
	if args.Get(1) != nil {
		return nil, dto.Paging{}, args.Error(1)
	}
	return args.Get(0).([]dto.Deposite), args.Get(1).(dto.Paging), nil
}

// FindById implements DepositeUseCaseMock.
func (d *DepositeUseCaseMock) FindById(id string) (int, dto.DepositeByIdResponse, error) {
	args := d.Called(id)
	if args.Get(2) != nil {
		return args.Int(0), dto.DepositeByIdResponse{}, args.Error(1)
	}
	return args.Int(0), args.Get(1).(dto.DepositeByIdResponse), nil
}

// FindByUcId implements DepositeUseCaseMock.
func (d *DepositeUseCaseMock) FindByUcId(id string) (int, dto.DepositeByUserResponse, error) {
	args := d.Called(id)
	if args.Get(2) != nil {
		return args.Int(0), dto.DepositeByUserResponse{}, args.Error(1)
	}
	return args.Int(0), args.Get(1).(dto.DepositeByUserResponse), nil
}

// CreateDeposite implements DepositeUseCaseMock.
func (d *DepositeUseCaseMock) CreateDeposite(payload dto.DepositeDto) (int, error) {
	args := d.Called(payload)
	if args.Get(1) != nil {
		return args.Int(0), args.Error(1)
	}
	return args.Int(0), nil
}
