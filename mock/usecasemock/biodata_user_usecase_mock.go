package usecasemock

import (
	"polen/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type BiodataUserUseCaseMock struct {
	mock.Mock
}

// Paging implements BiodataUserUseCase.
func (bio *BiodataUserUseCaseMock) Paging(payload dto.PageRequest) ([]dto.BiodataResponse, dto.Paging, error) {
	args := bio.Called(payload)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]dto.BiodataResponse), args.Get(1).(dto.Paging), nil
}

// FindByUcId implements BiodataUserUseCaseMock.
func (bio *BiodataUserUseCaseMock) FindByUcId(id string) (dto.BiodataResponse, error) {
	args := bio.Called(id)
	if args.Get(1) != nil {
		return dto.BiodataResponse{}, args.Error(1)
	}
	return args.Get(0).(dto.BiodataResponse), nil
}

// AdminUpdate implements BiodataUserUseCaseMock.
func (bio *BiodataUserUseCaseMock) AdminUpdate(payload dto.UpdateBioRequest, ctx *gin.Context) (int, error) {
	args := bio.Called(payload, ctx)
	if args.Get(1) != nil {
		return args.Int(0), args.Error(1)
	}
	return args.Int(0), nil
}

// UserUpdate implements BiodataUserUseCaseMock.
func (bio *BiodataUserUseCaseMock) UserUpdate(payload dto.BiodataRequest, ctx *gin.Context) (int, error) {
	args := bio.Called(payload, ctx)
	if args.Get(1) != nil {
		return args.Int(0), args.Error(1)
	}
	return args.Int(0), nil
}

// FindUserUpdated implements BiodataUserUseCaseMock.
func (bio *BiodataUserUseCaseMock) FindUserUpdated() ([]dto.BiodataResponse, error) {
	args := bio.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.BiodataResponse), nil
}

// FindByUserCredential implements BiodataUserUseCaseMock.
func (bio *BiodataUserUseCaseMock) FindByUserCredential(ctx *gin.Context) (dto.BiodataResponse, error) {
	args := bio.Called(ctx)
	if args.Get(1) != nil {
		return dto.BiodataResponse{}, args.Error(1)
	}
	return args.Get(0).(dto.BiodataResponse), nil
}
