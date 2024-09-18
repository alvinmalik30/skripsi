package repomock

import (
	"polen/model"
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type TopupRepoMock struct {
	mock.Mock
}

func (t *TopupRepoMock) Save(payload model.TopUp) error {
	return t.Called(payload).Error(0)
}
func (t *TopupRepoMock) FindByIdUser(id string) (dto.TopUpByUser, error) {
	args := t.Called(id)
	if args.Get(1) != nil {
		return dto.TopUpByUser{}, args.Error(1)
	}
	return args.Get(0).(dto.TopUpByUser), nil
}
func (t *TopupRepoMock) FindById(id string) (dto.TopUpById, error) {
	args := t.Called(id)
	if args.Get(1) != nil {
		return dto.TopUpById{}, args.Error(1)
	}
	return args.Get(0).(dto.TopUpById), nil
}
func (t *TopupRepoMock) Upload(payload model.TopUp) error {
	return t.Called(payload).Error(0)
}
func (t *TopupRepoMock) ConfimUpload(payload model.TopUp) error {
	return t.Called(payload).Error(0)
}
func (t *TopupRepoMock) NotConfimUpload(payload model.TopUp) error {
	return t.Called(payload).Error(0)
}
func (t *TopupRepoMock) FindUploadedFile() ([]dto.TopUp, error) {
	args := t.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.TopUp), nil
}
func (t *TopupRepoMock) Pagging(payload dto.PageRequest) ([]dto.TopUp, dto.Paging, error) {
	args := t.Called(payload)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]dto.TopUp), args.Get(1).(dto.Paging), nil
}
