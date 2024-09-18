package usecase

import (
	"polen/model/dto"
	"polen/repository"
)

type SaldoUsecase interface {
	FindByIdUser(id string) (dto.Saldo, error)
	Pagging(payload dto.PageRequest) ([]dto.Saldo, dto.Paging, error)
}

type saldoUsecase struct {
	repo repository.SaldoRepository
}

// FindByIdUser implements SaldoUsecase.
func (s *saldoUsecase) FindByIdUser(id string) (dto.Saldo, error) {
	return s.repo.FindByIdUser(id)
}

// Pagging implements SaldoUsecase.
func (s *saldoUsecase) Pagging(payload dto.PageRequest) ([]dto.Saldo, dto.Paging, error) {
	return s.repo.Pagging(payload)
}

func NewSaldoUsecase(repo repository.SaldoRepository) SaldoUsecase {
	return &saldoUsecase{
		repo: repo,
	}
}
