package usecasemock

import (
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type AuthUseCaseMock struct {
	mock.Mock
}

// Login implements AuthUseCase.
func (a *AuthUseCaseMock) Login(payload dto.AuthLoginRequest) (dto.AuthResponse, error) {
	a2 := a.Called(payload)
	if a2.Get(1) != nil {
		return dto.AuthResponse{}, a2.Error(1)
	}
	return a2.Get(0).(dto.AuthResponse), nil
}
