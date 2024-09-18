package usecase

import (
	"errors"
	"polen/mock"
	"polen/mock/repomock"
	"polen/mock/usecasemock"
	"polen/model/dto"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type BiodataUserUseCaseTestSuite struct {
	suite.Suite
	biurm *repomock.BiodataUserRepoMock
	uucm  *usecasemock.UserUseCaseMock
	buuc  BiodataUserUseCase
	ctx   *gin.Context
}

func (b *BiodataUserUseCaseTestSuite) SetupTest() {
	b.biurm = new(repomock.BiodataUserRepoMock)
	b.uucm = new(usecasemock.UserUseCaseMock)
	b.ctx = &gin.Context{}
	b.buuc = NewBiodataUserUseCase(b.biurm, b.uucm, b.ctx)
}
func TestBiodataUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(BiodataUserUseCaseTestSuite))
}
func (b *BiodataUserUseCaseTestSuite) TestPaging_Success() {
	b.biurm.On("Pagging", mock.MockPageReq).Return(mock.MockBiodataResponses, mock.MockPaging, nil)
	br, p, err := b.buuc.Paging(mock.MockPageReq)
	assert.Nil(b.T(), err)
	assert.Equal(b.T(), 1, p.TotalPages)
	assert.Equal(b.T(), 1, len(br))
}
func (b *BiodataUserUseCaseTestSuite) TestPaging_Failed() {
	b.biurm.On("Pagging", mock.MockPageReq).Return(nil, dto.Paging{}, errors.New("error"))
	br, p, err := b.buuc.Paging(mock.MockPageReq)
	assert.Error(b.T(), err)
	assert.Equal(b.T(), 0, p.TotalRows)
	assert.Equal(b.T(), 0, len(br))
}
func (b *BiodataUserUseCaseTestSuite) TestPaging_DefaultInputPage() {
	b.biurm.On("Pagging", dto.PageRequest{Page: 1, Size: 5}).Return(mock.MockBiodataResponses, mock.MockPaging, nil)
	br, p, err := b.buuc.Paging(dto.PageRequest{Page: -1, Size: 5})
	assert.Nil(b.T(), err)
	assert.Equal(b.T(), 1, p.TotalPages)
	assert.Equal(b.T(), 1, len(br))
}
func (b *BiodataUserUseCaseTestSuite) TestFindByUcId_Success() {
	b.biurm.On("FindByUcId", mock.MockUserCred.Id).Return(mock.MockBiodataResponse, nil)
	br, err := b.buuc.FindByUcId(mock.MockUserCred.Id)
	assert.Nil(b.T(), err)
	assert.Equal(b.T(), mock.MockUserCred.Id, br.UserCredential.Id)
}
func (b *BiodataUserUseCaseTestSuite) TestFindByUcId_Failed() {
	b.biurm.On("FindByUcId", mock.MockUserCred.Id).Return(mock.MockBiodataResponse, errors.New("error"))
	br, err := b.buuc.FindByUcId(mock.MockUserCred.Id)
	assert.Error(b.T(), err)
	assert.Equal(b.T(), dto.BiodataResponse{}, br)
}
func (b *BiodataUserUseCaseTestSuite) TestAdminUpdate_InfoRequired() {
	b.biurm.On("FindByUcId", mock.MockUpdateBioReq.UserCredentialId).Return(mock.MockBiodataResponse, nil)
	b.biurm.On("AdminUpdate", mock.MockBiodata).Return(nil)
	i, err := b.buuc.AdminUpdate(dto.UpdateBioRequest{UserCredentialId: "1"}, b.ctx)
	assert.Error(b.T(), err)
	assert.Equal(b.T(), 400, i)
}
func (b *BiodataUserUseCaseTestSuite) TestAdminUpdate_InvalidUCId() {
	b.biurm.On("FindByUcId", mock.MockUpdateBioReq.UserCredentialId).Return(mock.MockBiodataResponse, errors.New("error"))
	b.biurm.On("AdminUpdate", mock.MockBiodata).Return(nil)
	i, err := b.buuc.AdminUpdate(mock.MockUpdateBioReq, b.ctx)
	assert.Error(b.T(), err)
	assert.Equal(b.T(), 500, i)
}
func (b *BiodataUserUseCaseTestSuite) TestFindUserUpdated_Success() {
	b.biurm.On("FindUserUpdated").Return(mock.MockBiodataResponses, nil)
	br, err := b.buuc.FindUserUpdated()
	assert.Nil(b.T(), err)
	assert.Equal(b.T(), mock.MockBiodataResponses, br)
}
func (b *BiodataUserUseCaseTestSuite) TestFindUserUpdated_Failed() {
	b.biurm.On("FindUserUpdated").Return(mock.MockBiodataResponses, errors.New("error"))
	br, err := b.buuc.FindUserUpdated()
	assert.Error(b.T(), err)
	assert.Nil(b.T(), br)
}
