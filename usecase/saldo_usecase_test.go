package usecase

import (
	"errors"
	"polen/mock"
	"polen/mock/repomock"
	"polen/model/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type SaldoUseCaseTestSuite struct {
	suite.Suite
	srm *repomock.SaldoRepoMock
	suc SaldoUsecase
}

func (s *SaldoUseCaseTestSuite) SetupTest() {
	s.srm = new(repomock.SaldoRepoMock)
	s.suc = NewSaldoUsecase(s.srm)
}

func TestSaldoUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(SaldoUseCaseTestSuite))
}
func (s *SaldoUseCaseTestSuite) TestPagging_Success() {
	s.srm.On("Pagging", mock.MockPageReq).Return(mock.MockSaldoDatas, mock.MockPaging, nil)
	ahc, p, err := s.suc.Pagging(mock.MockPageReq)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), len(mock.MockSaldoDatas), len(ahc))
	assert.Equal(s.T(), mock.MockPaging.TotalRows, p.TotalRows)
}
func (s *SaldoUseCaseTestSuite) TestPagging_Failed() {
	s.srm.On("Pagging", mock.MockPageReq).Return(nil, dto.Paging{}, errors.New("failed"))
	ahc, p, err := s.suc.Pagging(mock.MockPageReq)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 0, len(ahc))
	assert.Equal(s.T(), dto.Paging{}, p)
}
func (s *SaldoUseCaseTestSuite) TestFindById_Success() {
	s.srm.On("FindByIdUser", "1").Return(mock.MockSaldoData, nil)
	ahc, err := s.suc.FindByIdUser("1")
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), mock.MockSaldoData.Id, ahc.Id)
}
func (s *SaldoUseCaseTestSuite) TestFindById_Failed() {
	s.srm.On("FindByIdUser", "1").Return(dto.Saldo{}, errors.New("failed"))
	ahc, err := s.suc.FindByIdUser("1")
	assert.Error(s.T(), err)
	assert.Equal(s.T(), dto.Saldo{}, ahc)
}
