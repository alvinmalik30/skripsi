package usecasemock

import (
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type TopupUseCaseMock struct {
	mock.Mock
}

// Pagging implements TopUpUseCase.
func (t *TopupUseCaseMock) Pagging(payload dto.PageRequest) ([]dto.TopUp, dto.Paging, error) {
	args := t.Called(payload)
	if args.Get(1) != nil {
		return nil, dto.Paging{}, args.Error(1)
	}
	return args.Get(0).([]dto.TopUp), args.Get(1).(dto.Paging), nil
}

// FindUploadedFile implements TopupUseCaseMock.
func (t *TopupUseCaseMock) FindUploadedFile() ([]dto.TopUp, error) {
	args := t.Called()
	if args.Get(1) != nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]dto.TopUp), nil
}

// ConfimUploadFile implements TopupUseCaseMock.
func (t *TopupUseCaseMock) ConfimUploadFile(payload dto.TopUpUser) (int, error) {
	args := t.Called(payload)
	if args.Get(1) != nil {
		return args.Int(0), args.Error(1)
	}
	return args.Int(0), nil
}

// FindById implements TopupUseCaseMock.
func (t *TopupUseCaseMock) FindById(id string) (dto.TopUpById, error) {
	args := t.Called(id)
	if args.Get(1) != nil {
		return dto.TopUpById{}, args.Error(1)
	}
	return args.Get(0).(dto.TopUpById), nil
}

// FindByIdUser implements TopupUseCaseMock.
func (t *TopupUseCaseMock) FindByIdUser(id string) (dto.TopUpByUser, error) {
	args := t.Called(id)
	if args.Get(1) != nil {
		return dto.TopUpByUser{}, args.Error(1)
	}
	return args.Get(0).(dto.TopUpByUser), nil
}

func (t *TopupUseCaseMock) UploadFile(payload dto.TopUpUser) (int, error) {
	args := t.Called(payload)
	if args.Get(1) != nil {
		return args.Int(0), args.Error(1)
	}
	return args.Int(0), nil
}

// CreateNew implements TopupUseCaseMock.
func (t *TopupUseCaseMock) CreateNew(payload dto.TopUpUser) (dto.TopUpUser, error) {
	args := t.Called(payload)
	if args.Get(1) != nil {
		return dto.TopUpUser{}, args.Error(1)
	}
	return args.Get(0).(dto.TopUpUser), nil
}
