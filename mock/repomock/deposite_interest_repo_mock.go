package repomock

import (
	"polen/model"
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type DepositeInterestRepoMock struct {
	mock.Mock
}

func (d *DepositeInterestRepoMock) Pagging(payload dto.PageRequest) ([]dto.DepositeInterestRequest, dto.Paging, error) {
	args := d.Called(payload)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]dto.DepositeInterestRequest), args.Get(1).(dto.Paging), nil
}

// DeleteById implements DepositeInterest.
func (d *DepositeInterestRepoMock) DeleteById(id string) error {
	return d.Called(id).Error(0)
}

// FindById implements DepositeInterest.
func (d *DepositeInterestRepoMock) FindById(id string) (dto.DepositeInterestRequest, error) {
	args := d.Called(id)
	if args.Get(1) != nil {
		return dto.DepositeInterestRequest{}, args.Error(1)
	}
	return args.Get(0).(dto.DepositeInterestRequest), nil
}

// Update implements DepositeInterest.
func (d *DepositeInterestRepoMock) Update(payload dto.DepositeInterestRequest) error {
	return d.Called(payload).Error(0)
}

// Save implements DepositeIntereset.
func (d *DepositeInterestRepoMock) Save(payload model.DepositeInterest) error {
	return d.Called(payload).Error(0)
}
