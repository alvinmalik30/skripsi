package repomock

import (
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type DepositeRepoMock struct {
	mock.Mock
}

// Update implements DepositeRepoMock.
func (d *DepositeRepoMock) Update() error {
	return d.Called().Error(0)
}

// Pagging implements DepositeRepoMock.
func (d *DepositeRepoMock) Pagging(payload dto.PageRequest) ([]dto.Deposite, dto.Paging, error) {
	args := d.Called(payload)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]dto.Deposite), args.Get(1).(dto.Paging), nil
}

// FindById implements DepositeRepoMock.
func (d *DepositeRepoMock) FindById(id string) (dto.DepositeByIdResponse, error) {
	args := d.Called(id)
	if args.Get(1) != nil {
		return dto.DepositeByIdResponse{}, args.Error(1)
	}
	return args.Get(0).(dto.DepositeByIdResponse), nil
}

// FindByUcId implements DepositeRepoMock.
func (d *DepositeRepoMock) FindByUcId(id string) (dto.DepositeByUserResponse, error) {
	args := d.Called(id)
	if args.Get(1) != nil {
		return dto.DepositeByUserResponse{}, args.Error(1)
	}
	return args.Get(0).(dto.DepositeByUserResponse), nil
}

// CreateDeposite implements DepositeRepoMock.
func (d *DepositeRepoMock) CreateDeposite(payload dto.DepositeDto) error {
	return d.Called(payload).Error(0)
}
