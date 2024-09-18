package usecasemock

import (
	"polen/model"
	"polen/model/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/mock"
)

type UserUseCaseMock struct {
	mock.Mock
}

// Paging implements UserUseCase.
func (u *UserUseCaseMock) Paging(payload dto.PageRequest, ctx *gin.Context) ([]model.UserCredential, dto.Paging, error) {
	args := u.Called(payload, ctx)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.UserCredential), args.Get(1).(dto.Paging), nil
}

// FindById implements usecase.UserUseCase.
func (u *UserUseCaseMock) FindById(id string) (model.UserCredential, error) {
	args := u.Called(id)
	if args.Get(1) != nil {
		return model.UserCredential{}, args.Error(1)
	}
	return args.Get(0).(model.UserCredential), nil
}

func (u *UserUseCaseMock) FindByUsername(username string, ctx *gin.Context) (model.UserCredential, error) {
	args := u.Called(username, ctx)
	if args.Get(1) != nil {
		return model.UserCredential{}, args.Error(1)
	}
	return args.Get(0).(model.UserCredential), nil
}

func (u *UserUseCaseMock) Register(payload dto.AuthRequest) error {
	return u.Called(payload).Error(0)
}
