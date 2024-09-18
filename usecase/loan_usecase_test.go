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

type LoanUseCaseTestSuite struct {
	suite.Suite
	lrm   *repomock.LoanRepoMock
	lirm  *repomock.LoanInterestRepoMock
	ahcrm *repomock.ApplicationHCRepoMock
	lprm  *repomock.LatePaymentRepoMock
	luc   LoanUseCase
	liuc  LoanInterestUseCase
	ahcuc AppHandlingCostUsecase
	lpfuc LatePaymentFeeUsecase
}

func (l *LoanUseCaseTestSuite) SetupTest() {
	l.lrm = new(repomock.LoanRepoMock)
	l.lirm = new(repomock.LoanInterestRepoMock)
	l.ahcrm = new(repomock.ApplicationHCRepoMock)
	l.lprm = new(repomock.LatePaymentRepoMock)
	l.lpfuc = NewLatePaymentFeeUseCase(l.lprm)
	l.ahcuc = NewAppHandlingCostUseCase(l.ahcrm)
	l.liuc = NewLoanInterestUseCase(l.lirm)
	l.luc = NewLoanUseCase(l.lrm, l.liuc, l.ahcuc, l.lpfuc)
}

func TestLoanUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(LoanUseCaseTestSuite))
}
func (l *LoanUseCaseTestSuite) TestFindById_Success() {
	l.lrm.On("FindById", "1").Return(mock.MockInstallLoanByIdResp, nil)
	ilbir, err := l.luc.FindById("1")
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), mock.MockInstallLoanByIdResp.LoanId, ilbir.LoanId)
}
func (l *LoanUseCaseTestSuite) TestFindById_Failed() {
	l.lrm.On("FindById", "1").Return(mock.MockInstallLoanByIdResp, errors.New("error"))
	ilbir, err := l.luc.FindById("1")
	assert.Error(l.T(), err)
	assert.Equal(l.T(), dto.InstallenmentLoanByIdResponse{}, ilbir)
}
func (l *LoanUseCaseTestSuite) TestFindUploadedFile_Success() {
	l.lrm.On("FindUploadedFile").Return(mock.MockLoanInstallRespons, nil)
	ilbir, err := l.luc.FindUploadedFile()
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), len(mock.MockLoanInstallRespons), len(ilbir))
}
func (l *LoanUseCaseTestSuite) TestFindUploadedFile_Failed() {
	l.lrm.On("FindUploadedFile").Return(mock.MockLoanInstallRespons, errors.New("error"))
	ilbir, err := l.luc.FindUploadedFile()
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 0, len(ilbir))
}
func (l *LoanUseCaseTestSuite) TestUpdateLateFee_Success() {
	l.lrm.On("UpdateLateFee").Return(nil)
	err := l.luc.UpdateLateFee()
	assert.Nil(l.T(), err)
}
func (l *LoanUseCaseTestSuite) TestUpdateLateFee_Failed() {
	l.lrm.On("UpdateLateFee").Return(errors.New("error"))
	err := l.luc.UpdateLateFee()
	assert.Error(l.T(), err)
}
func (l *LoanUseCaseTestSuite) TestAccepted_Success() {
	l.lrm.On("Accepted", mock.MockInstallLoanByIdResp).Return(nil)
	err := l.luc.Accepted(mock.MockInstallLoanByIdResp)
	assert.Nil(l.T(), err)
}
func (l *LoanUseCaseTestSuite) TestAccepted_Failed() {
	l.lrm.On("Accepted", mock.MockInstallLoanByIdResp).Return(errors.New("error"))
	err := l.luc.Accepted(mock.MockInstallLoanByIdResp)
	assert.Error(l.T(), err)
}
func (l *LoanUseCaseTestSuite) TestFindByLoanId_Success() {
	l.lrm.On("FindByLoanId", "1").Return(mock.MockDTOLoan, nil)
	l2, err := l.luc.FindByLoanId("1")
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), mock.MockDTOLoan.Id, l2.Id)
}
func (l *LoanUseCaseTestSuite) TestFindByLoanId_Failed() {
	l.lrm.On("FindByLoanId", "1").Return(mock.MockDTOLoan, errors.New("error"))
	l2, err := l.luc.FindByLoanId("1")
	assert.Error(l.T(), err)
	assert.Equal(l.T(), dto.Loan{}, l2)
}
