package usecase

import (
	"errors"
	"polen/mock"
	"polen/mock/repomock"
	"polen/model"
	"polen/model/dto"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type LoanInterestUseCaseTestSuite struct {
	suite.Suite
	lirm *repomock.LoanInterestRepoMock
	liuc LoanInterestUseCase
}

func (l *LoanInterestUseCaseTestSuite) SetupTest() {
	l.lirm = new(repomock.LoanInterestRepoMock)
	l.liuc = NewLoanInterestUseCase(l.lirm)
}

func TestLoanInterestUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(LoanInterestUseCaseTestSuite))
}
func (l *LoanInterestUseCaseTestSuite) TestPagging_Success() {
	l.lirm.On("Pagging", mock.MockPageReq).Return(mock.MockLoanInterestDatas, mock.MockPaging, nil)
	ahc, p, err := l.liuc.Pagging(mock.MockPageReq)
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), len(mock.MockLoanInterestDatas), len(ahc))
	assert.Equal(l.T(), mock.MockPaging.TotalRows, p.TotalRows)
}
func (l *LoanInterestUseCaseTestSuite) TestPagging_Failed() {
	l.lirm.On("Pagging", mock.MockPageReq).Return(nil, dto.Paging{}, errors.New("failed"))
	ahc, p, err := l.liuc.Pagging(mock.MockPageReq)
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 0, len(ahc))
	assert.Equal(l.T(), dto.Paging{}, p)
}
func (l *LoanInterestUseCaseTestSuite) TestFindById_Success() {
	l.lirm.On("FindById", "1").Return(mock.MockLoanInterest, nil)
	ahc, err := l.liuc.FindById("1")
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), mock.MockAppHC.Id, ahc.Id)
}
func (l *LoanInterestUseCaseTestSuite) TestFindById_Failed() {
	l.lirm.On("FindById", "1").Return(model.LoanInterest{}, errors.New("failed"))
	ahc, err := l.liuc.FindById("1")
	assert.Error(l.T(), err)
	assert.Equal(l.T(), model.LoanInterest{}, ahc)
}
func (l *LoanInterestUseCaseTestSuite) TestDeleteById_Success() {
	l.lirm.On("FindById", "1").Return(mock.MockLoanInterest, nil)
	l.lirm.On("DeleteById", "1").Return(nil)
	err := l.liuc.DeleteById("1")
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LoanInterestUseCaseTestSuite) TestDeleteById_Failed() {
	// failed find id
	l.lirm.On("FindById", "1").Return(model.LoanInterest{}, errors.New("failed id"))
	err := l.liuc.DeleteById("1")
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LoanInterestUseCaseTestSuite) TestDeleteById_InvalidServer() {
	// failed server error
	l.lirm.On("FindById", "1").Return(mock.MockLoanInterest, nil)
	l.lirm.On("DeleteById", "1").Return(errors.New("failed to delete app handling cost"))
	err := l.liuc.DeleteById("1")
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LoanInterestUseCaseTestSuite) TestCreateNew_Success() {
	l.lirm.On("CreateNew", mock.MockLoanInterest).Return(nil)
	i, err := l.liuc.CreateNew(mock.MockLoanInterest)
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), 201, i)
}
func (l *LoanInterestUseCaseTestSuite) TestCreateNew_Failed() {
	// id required
	l.lirm.On("CreateNew", model.LoanInterest{}).Return(errors.New("id is required"))
	i, err := l.liuc.CreateNew(model.LoanInterest{})
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 400, i)
	// name required
	l.lirm.On("CreateNew", model.LoanInterest{Id: "1"}).Return(errors.New("duration month is required"))
	i, err = l.liuc.CreateNew(model.LoanInterest{Id: "1"})
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 400, i)
	// nominal required
	l.lirm.On("CreateNew", model.LoanInterest{Id: "1", DurationMonths: 12}).Return(errors.New("loan interest rate is required"))
	i, err = l.liuc.CreateNew(model.LoanInterest{Id: "1", DurationMonths: 12})
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 400, i)
	// server error
	l.lirm.On("CreateNew", mock.MockLoanInterest).Return(errors.New("server error"))
	i, err = l.liuc.CreateNew(mock.MockLoanInterest)
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 500, i)
}
func (l *LoanInterestUseCaseTestSuite) TestUpdate_Success() {
	l.lirm.On("FindById", "1").Return(mock.MockLoanInterest, nil)
	l.lirm.On("Update", mock.MockLoanInterest).Return(nil)
	err := l.liuc.Update(mock.MockLoanInterest)
	assert.Nil(l.T(), err)
}
func (l *LoanInterestUseCaseTestSuite) TestUpdate_Failed() {
	// id required
	l.lirm.On("Update", model.LoanInterest{}).Return(errors.New("id is required"))
	err := l.liuc.Update(model.LoanInterest{})
	assert.Error(l.T(), err)
	// id db failed
	l.lirm.On("FindById", "1").Return(model.LoanInterest{}, errors.New("failed"))
	l.lirm.On("Update", mock.MockLoanInterest).Return(errors.New("error"))
	err = l.liuc.Update(mock.MockLoanInterest)
	assert.Error(l.T(), err)
}
func (l *LoanInterestUseCaseTestSuite) TestUpdate_InvalidServer() {
	l.lirm.On("FindById", "1").Return(mock.MockLoanInterest, nil)
	l.lirm.On("Update", mock.MockLoanInterest).Return(errors.New("server error"))
	err := l.liuc.Update(mock.MockLoanInterest)
	assert.Error(l.T(), err)
}
