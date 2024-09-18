package repomock

import (
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type SaldoRepoMock struct {
	mock.Mock
}

// FindByIdUser implements SaldoRepository.
func (s *SaldoRepoMock) FindByIdUser(id string) (dto.Saldo, error) {
	args := s.Called(id)
	if args.Get(1) != nil {
		return dto.Saldo{}, args.Error(1)
	}
	return args.Get(0).(dto.Saldo), nil
}

// Pagging implements SaldoRepository.
func (s *SaldoRepoMock) Pagging(payload dto.PageRequest) ([]dto.Saldo, dto.Paging, error) {
	args := s.Called(payload)
	if args.Get(2) != nil {
		return nil, dto.Paging{}, args.Error(2)
	}
	return args.Get(0).([]dto.Saldo), args.Get(1).(dto.Paging), nil
}
