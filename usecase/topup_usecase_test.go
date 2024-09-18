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

type TopupUseCaseTestSuite struct {
	suite.Suite
	trm  *repomock.TopupRepoMock
	urm  *repomock.UserRepoMock
	uuc  UserUseCase
	tuuc TopUpUseCase
}

func (t *TopupUseCaseTestSuite) SetupTest() {
	t.trm = new(repomock.TopupRepoMock)
	t.urm = new(repomock.UserRepoMock)
	t.tuuc = NewTopUpUseCase(t.trm, t.uuc)
}

func TestTopupUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(TopupUseCaseTestSuite))
}

func (t *TopupUseCaseTestSuite) TestPagging_Success() {
	t.trm.On("Pagging", mock.MockPageReq).Return(mock.MockListTopUp, mock.MockPaging, nil)
	ahc, p, err := t.tuuc.Pagging(mock.MockPageReq)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), len(mock.MockListTopUp), len(ahc))
	assert.Equal(t.T(), mock.MockPaging.TotalRows, p.TotalRows)
}
func (t *TopupUseCaseTestSuite) TestPagging_Failed() {
	t.trm.On("Pagging", mock.MockPageReq).Return(nil, dto.Paging{}, errors.New("failed"))
	ahc, p, err := t.tuuc.Pagging(mock.MockPageReq)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), 0, len(ahc))
	assert.Equal(t.T(), dto.Paging{}, p)
}
func (t *TopupUseCaseTestSuite) TestPagging_DefaultPage() {
	t.trm.On("Pagging", mock.MockPageReq).Return(mock.MockListTopUp, mock.MockPaging, nil)
	ahc, p, err := t.tuuc.Pagging(dto.PageRequest{Page: -1, Size: 5})
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), len(mock.MockListTopUp), len(ahc))
	assert.Equal(t.T(), mock.MockPaging.TotalRows, p.TotalRows)
}
func (t *TopupUseCaseTestSuite) TestFindUploadedFile_Success() {
	t.trm.On("FindUploadedFile").Return(mock.MockListTopUp, nil)
	tu, err := t.tuuc.FindUploadedFile()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), len(mock.MockListTopUp), len(tu))
}
func (t *TopupUseCaseTestSuite) TestFindUploadedFile_Failed() {
	t.trm.On("FindUploadedFile").Return(mock.MockListTopUp, errors.New("error"))
	tu, err := t.tuuc.FindUploadedFile()
	assert.Error(t.T(), err)
	assert.Equal(t.T(), 0, len(tu))
}
func (t *TopupUseCaseTestSuite) TestFindById_Success() {
	t.trm.On("FindById", "1").Return(mock.MockTopUpId, nil)
	tu, err := t.tuuc.FindById("1")
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), mock.MockTopUpId.TopUp.Id, tu.TopUp.Id)
}
func (t *TopupUseCaseTestSuite) TestFindByIde_Failed() {
	t.trm.On("FindById", "1").Return(dto.TopUpById{}, errors.New("error"))
	tu, err := t.tuuc.FindById("1")
	assert.Error(t.T(), err)
	assert.Equal(t.T(), dto.TopUpById{}, tu)
}
func (t *TopupUseCaseTestSuite) TestFindByIdUser_Success() {
	t.trm.On("FindByIdUser", "1").Return(mock.MockTopUpByUser, nil)
	tu, err := t.tuuc.FindByIdUser("1")
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), mock.MockTopUpByUser.UserCredential.Id, tu.UserCredential.Id)
}
func (t *TopupUseCaseTestSuite) TestFindByIdUser_Failed() {
	t.trm.On("FindByIdUser", "1").Return(dto.TopUpByUser{}, errors.New("error"))
	tu, err := t.tuuc.FindByIdUser("1")
	assert.Error(t.T(), err)
	assert.Equal(t.T(), dto.TopUpByUser{}, tu)
}
