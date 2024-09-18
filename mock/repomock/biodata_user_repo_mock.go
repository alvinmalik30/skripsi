package repomock

import (
	"polen/model"
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type BiodataUserRepoMock struct {
	mock.Mock
}

// Pagging implements BiodataUser.
func (bio *BiodataUserRepoMock) Pagging(payload dto.PageRequest) ([]dto.BiodataResponse, dto.Paging, error) {
	args := bio.Called(payload)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]dto.BiodataResponse), args.Get(1).(dto.Paging), nil
}

// AdminUpdate implements BiodataUser.
func (bio *BiodataUserRepoMock) AdminUpdate(payload model.BiodataUser) error {
	return bio.Called(payload).Error(0)
}

// UserUpdate implements BiodataUser.
func (bio *BiodataUserRepoMock) UserUpdate(payload model.BiodataUser) error {
	return bio.Called(payload).Error(0)
}

// FindByUcId implements BiodataUser.
func (bio *BiodataUserRepoMock) FindByUcId(id string) (dto.BiodataResponse, error) {
	args := bio.Called(id)
	if args.Get(1) != nil {
		return dto.BiodataResponse{}, args.Error(1)
	}
	return args.Get(0).(dto.BiodataResponse), nil
}

// FindByUcId implements BiodataUser.
func (bio *BiodataUserRepoMock) FindUserUpdated() ([]dto.BiodataResponse, error) {
	args := bio.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.BiodataResponse), nil
}
