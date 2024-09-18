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

type DepositeUseCaseTestSuite struct {
	suite.Suite
	drm  *repomock.DepositeRepoMock
	dirm *repomock.DepositeInterestRepoMock
	srm  *repomock.SaldoRepoMock
	duc  DepositeUseCase
	diuc DepositeInterestUseCase
	suc  SaldoUsecase
}

func (d *DepositeUseCaseTestSuite) SetupTest() {
	d.drm = new(repomock.DepositeRepoMock)
	d.suc = NewSaldoUsecase(d.srm)
	d.diuc = NewDepositeInterestUseCase(d.dirm)
	d.duc = NewDepositeUseCase(d.drm, d.diuc, d.suc)
}

func TestDepositeUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(DepositeUseCaseTestSuite))
}
func (d *DepositeUseCaseTestSuite) TestUpdate_Success() {
	d.drm.On("Update").Return(nil)
	err := d.duc.Update()
	assert.Nil(d.T(), err)
}
func (d *DepositeUseCaseTestSuite) TestUpdate_Failed() {
	d.drm.On("Update").Return(errors.New("failed"))
	err := d.duc.Update()
	assert.Error(d.T(), err)
}
func (d *DepositeUseCaseTestSuite) TestPagging_Success() {
	d.drm.On("Pagging", mock.MockPageReq).Return(mock.MockDepositeDto, mock.MockPaging, nil)
	ahc, p, err := d.duc.Pagging(mock.MockPageReq)
	assert.Nil(d.T(), err)
	assert.Equal(d.T(), len(mock.MockDepositeDto), len(ahc))
	assert.Equal(d.T(), mock.MockPaging.TotalRows, p.TotalRows)
}
func (d *DepositeUseCaseTestSuite) TestPagging_Failed() {
	d.drm.On("Pagging", mock.MockPageReq).Return(nil, dto.Paging{}, errors.New("failed"))
	ahc, p, err := d.duc.Pagging(mock.MockPageReq)
	assert.Error(d.T(), err)
	assert.Equal(d.T(), 0, len(ahc))
	assert.Equal(d.T(), dto.Paging{}, p)
}
func (d *DepositeUseCaseTestSuite) TestPagging_DefaultPage() {
	d.drm.On("Pagging", mock.MockPageReq).Return(mock.MockDepositeDto, mock.MockPaging, nil)
	ahc, p, err := d.duc.Pagging(dto.PageRequest{Page: -1, Size: 5})
	assert.Nil(d.T(), err)
	assert.Equal(d.T(), len(mock.MockDepositeDto), len(ahc))
	assert.Equal(d.T(), mock.MockPaging.TotalRows, p.TotalRows)
}
func (d *DepositeUseCaseTestSuite) TestFindById_Success() {
	d.drm.On("FindById", "1").Return(mock.MockDepositeByIdResponse, nil)
	i, dbir, err := d.duc.FindById("1")
	assert.Nil(d.T(), err)
	assert.Equal(d.T(), mock.MockDepositeByIdResponse.Deposite.Id, dbir.Deposite.Id)
	assert.Equal(d.T(), 200, i)
}
func (d *DepositeUseCaseTestSuite) TestFindById_Failed() {
	d.drm.On("FindById", "1").Return(mock.MockDepositeByIdResponse, errors.New("failed"))
	i, dbir, err := d.duc.FindById("1")
	assert.Error(d.T(), err)
	assert.Equal(d.T(), dto.DepositeByIdResponse{}, dbir)
	assert.Equal(d.T(), 500, i)
}
func (d *DepositeUseCaseTestSuite) TestFindById_IdRequired() {
	d.drm.On("FindById", "1").Return(mock.MockDepositeByIdResponse, nil)
	i, dbir, err := d.duc.FindById("")
	assert.Error(d.T(), err)
	assert.Equal(d.T(), dto.DepositeByIdResponse{}, dbir)
	assert.Equal(d.T(), 400, i)
}
func (d *DepositeUseCaseTestSuite) TestFindByUcId_Success() {
	d.drm.On("FindByUcId", "1").Return(mock.MockDepositeByUserResponse, nil)
	i, dbir, err := d.duc.FindByUcId("1")
	assert.Nil(d.T(), err)
	assert.Equal(d.T(), mock.MockDepositeByUserResponse.BioUser.Id, dbir.BioUser.Id)
	assert.Equal(d.T(), 200, i)
}
func (d *DepositeUseCaseTestSuite) TestFindByUcId_Failed() {
	d.drm.On("FindByUcId", "1").Return(mock.MockDepositeByUserResponse, errors.New("failed"))
	i, dbir, err := d.duc.FindByUcId("1")
	assert.Error(d.T(), err)
	assert.Equal(d.T(), dto.DepositeByUserResponse{}, dbir)
	assert.Equal(d.T(), 500, i)
}
func (d *DepositeUseCaseTestSuite) TestFindByUcId_IdRequired() {
	d.drm.On("FindByUcId", "1").Return(mock.MockDepositeByUserResponse, nil)
	i, dbir, err := d.duc.FindByUcId("")
	assert.Error(d.T(), err)
	assert.Equal(d.T(), dto.DepositeByUserResponse{}, dbir)
	assert.Equal(d.T(), 400, i)
}
func (d *DepositeUseCaseTestSuite) TestCreateDeposite_EmptyInvalid() {
	// interest rate is required
	d.drm.On("CreateDeposite", mock.MockDeposite).Return(nil)
	i, err := d.duc.CreateDeposite(dto.DepositeDto{})
	assert.Error(d.T(), err)
	assert.Equal(d.T(), 400, i)
	// deposite amount must greather than zero
	d.drm.On("CreateDeposite", mock.MockDeposite).Return(nil)
	i, err = d.duc.CreateDeposite(dto.DepositeDto{InterestRate: dto.DepositeInterestRequest{Id: "1"}, DepositeAmount: 0})
	assert.Error(d.T(), err)
	assert.Equal(d.T(), 400, i)
}
