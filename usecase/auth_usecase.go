package usecase

import (
	"polen/model/dto"
	"polen/repository"
	"polen/utils/security"
)

type AuthUseCase interface {
	Login(payload dto.AuthLoginRequest) (dto.AuthResponse, error)
}

type authUseCase struct {
	repo repository.UserRepository
}

// Login implements AuthUseCase.
func (a *authUseCase) Login(payload dto.AuthLoginRequest) (dto.AuthResponse, error) {
	// Username di db
	user, err := a.repo.FindByUsername(payload.Username)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	// Validasi Password
	err = security.VerifyPassword(user.Password, payload.Password)
	if err != nil {
		return dto.AuthResponse{}, err
	}

	// Generate Token
	token, err := security.GenerateJwtToken(user)
	if err != nil {
		return dto.AuthResponse{}, err
	}
	return dto.AuthResponse{
		Username: user.Username,
		Token:    token,
	}, nil
}

func NewAuthUseCase(repo repository.UserRepository) AuthUseCase {
	return &authUseCase{repo: repo}
}
