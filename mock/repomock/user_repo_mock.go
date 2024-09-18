package repomock

import (
	"polen/model"
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type UserRepoMock struct {
	mock.Mock
}

func (u *UserRepoMock) Pagging(payload dto.PageRequest) ([]model.UserCredential, dto.Paging, error) {
	args := u.Called(payload)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]model.UserCredential), args.Get(1).(dto.Paging), nil
}
func (u *UserRepoMock) Saldo(payload model.UserCredential, idsaldo string, bioId string) error {
	return u.Called(payload, idsaldo, bioId).Error(0)
}
func (u *UserRepoMock) FindById(id string) (model.UserCredential, error) {
	args := u.Called(id)
	if args.Get(1) != nil {
		return model.UserCredential{}, args.Error(1)
	}
	return args.Get(0).(model.UserCredential), nil
}
func (u *UserRepoMock) FindByUsername(username string) (model.UserCredential, error) {
	args := u.Called(username)
	if args.Get(1) != nil {
		return model.UserCredential{}, args.Error(1)
	}
	return args.Get(0).(model.UserCredential), nil
}
func (u *UserRepoMock) Save(payload model.UserCredential, bioId string) error {
	return u.Called(payload, bioId).Error(0)
}
