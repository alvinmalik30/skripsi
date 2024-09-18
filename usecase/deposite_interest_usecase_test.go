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

type DepositeInteresetUseCaseTestSuite struct {
	suite.Suite
	dirm *repomock.DepositeInterestRepoMock
	diuc DepositeInterestUseCase
}

func (d *DepositeInteresetUseCaseTestSuite) SetupTest() {
	d.dirm = new(repomock.DepositeInterestRepoMock)
	d.diuc = NewDepositeInterestUseCase(d.dirm)
}
func TestDepositeInteresetUseCaseSuite(t *testing.T) {
	suite.Run(t, new(DepositeInteresetUseCaseTestSuite))
}
func (d *DepositeInteresetUseCaseTestSuite) TestDeleteById_Success() {
	d.dirm.On("FindById", mock.MockDepositeInterestReq.Id).Return(mock.MockDepositeInterestReq, nil)
	d.dirm.On("DeleteById", mock.MockDepositeInterestReq.Id).Return(nil)
	err := d.diuc.DeleteById(mock.MockDepositeInterestReq.Id)
	assert.Nil(d.T(), err)
}
func (d *DepositeInteresetUseCaseTestSuite) TestDeleteById_Failed() {
	d.dirm.On("FindById", mock.MockDepositeInterestReq.Id).Return(mock.MockDepositeInterestReq, nil)
	d.dirm.On("DeleteById", mock.MockDepositeInterestReq.Id).Return(errors.New("failed to delete deposite"))
	err := d.diuc.DeleteById(mock.MockDepositeInterestReq.Id)
	assert.Error(d.T(), err)
}
func (d *DepositeInteresetUseCaseTestSuite) TestDeleteById_IdDepositeInvalid() {
	d.dirm.On("FindById", mock.MockDepositeInterestReq.Id).Return(mock.MockDepositeInterest, errors.New("error"))
	d.dirm.On("DeleteById", mock.MockDepositeInterestReq.Id).Return(nil)
	err := d.diuc.DeleteById(mock.MockDepositeInterestReq.Id)
	assert.Error(d.T(), err)
}
func (d *DepositeInteresetUseCaseTestSuite) TestFindById_Success() {
	d.dirm.On("FindById", mock.MockDepositeInterestReq.Id).Return(mock.MockDepositeInterestReq, nil)
	di, err := d.diuc.FindById(mock.MockDepositeInterestReq.Id)
	assert.Nil(d.T(), err)
	assert.Equal(d.T(), mock.MockDepositeInterestReq.Id, di.Id)
}
func (d *DepositeInteresetUseCaseTestSuite) TestFindById_Failed() {
	d.dirm.On("FindById", mock.MockDepositeInterestReq.Id).Return(mock.MockDepositeInterestReq, errors.New("error"))
	di, err := d.diuc.FindById(mock.MockDepositeInterestReq.Id)
	assert.Error(d.T(), err)
	assert.Equal(d.T(), dto.DepositeInterestRequest{}, di)
}
func (d *DepositeInteresetUseCaseTestSuite) TestUpdate_Success() {
	d.dirm.On("FindById", mock.MockDepositeInterestReq.Id).Return(mock.MockDepositeInterestReq, nil)
	d.dirm.On("Update", mock.MockDepositeInterestReq).Return(nil)
	err := d.diuc.Update(mock.MockDepositeInterestReq)
	assert.Nil(d.T(), err)
}
func (d *DepositeInteresetUseCaseTestSuite) TestUpdate_Failed() {
	d.dirm.On("FindById", mock.MockDepositeInterestReq.Id).Return(mock.MockDepositeInterestReq, nil)
	d.dirm.On("Update", mock.MockDepositeInterestReq).Return(errors.New("failed to update deposite"))
	err := d.diuc.Update(mock.MockDepositeInterestReq)
	assert.Error(d.T(), err)
}
func (d *DepositeInteresetUseCaseTestSuite) TestUpdate_CheckGTZero() {
	d.dirm.On("FindById", mock.MockDepositeInterestReq.Id).Return(mock.MockDepositeInterestReq, nil)
	d.dirm.On("Update", mock.MockDepositeInterestReq).Return(nil)
	err := d.diuc.Update(dto.DepositeInterestRequest{Id: "1"})
	assert.Nil(d.T(), err)
}
func (d *DepositeInteresetUseCaseTestSuite) TestUpdate_EmptyInvalid() {
	// Id required
	d.dirm.On("FindById", mock.MockDepositeInterestReq.Id).Return(mock.MockDepositeInterestReq, nil)
	d.dirm.On("Update", dto.DepositeInterestRequest{}).Return(errors.New("id is required"))
	err := d.diuc.Update(dto.DepositeInterestRequest{})
	assert.Error(d.T(), err)
}
func (d *DepositeInteresetUseCaseTestSuite) TestUpdate_IdDepositeInvalid() {
	d.dirm.On("FindById", mock.MockDepositeInterestReq.Id).Return(mock.MockDepositeInterestReq, errors.New("error"))
	d.dirm.On("Update", mock.MockDepositeInterestReq).Return(nil)
	err := d.diuc.Update(mock.MockDepositeInterestReq)
	assert.Error(d.T(), err)
}
func (d *DepositeInteresetUseCaseTestSuite) TestCreateNew_EmptyInvalid() {
	// id required
	d.dirm.On("Save", model.DepositeInterest{}).Return(errors.New("id is required"))
	i, err := d.diuc.CreateNew(dto.DepositeInterestRequest{})
	assert.Error(d.T(), err)
	assert.Equal(d.T(), 400, i)
	// interest rate required
	d.dirm.On("Save", model.DepositeInterest{Id: "1"}).Return(errors.New("interest rate is required"))
	i, err = d.diuc.CreateNew(dto.DepositeInterestRequest{Id: "1"})
	assert.Error(d.T(), err)
	assert.Equal(d.T(), 400, i)
	// tax is required
	d.dirm.On("Save", model.DepositeInterest{Id: "1", InterestRate: 1}).Return(errors.New("tax is required"))
	i, err = d.diuc.CreateNew(dto.DepositeInterestRequest{Id: "1", InterestRate: 1})
	assert.Error(d.T(), err)
	assert.Equal(d.T(), 400, i)
	// duration month required
	d.dirm.On("Save", model.DepositeInterest{Id: "1", InterestRate: 1, TaxRate: 1}).Return(errors.New("duration month is required"))
	i, err = d.diuc.CreateNew(dto.DepositeInterestRequest{Id: "1", InterestRate: 1, TaxRate: 1})
	assert.Error(d.T(), err)
	assert.Equal(d.T(), 400, i)
}
func (d *DepositeInteresetUseCaseTestSuite) TestPagging_Success() {
	d.dirm.On("Pagging", mock.MockPageReq).Return(mock.MockDeposites, mock.MockPaging, nil)
	ahc, p, err := d.diuc.Pagging(mock.MockPageReq)
	assert.Nil(d.T(), err)
	assert.Equal(d.T(), len(mock.MockDeposites), len(ahc))
	assert.Equal(d.T(), mock.MockPaging.TotalRows, p.TotalRows)
}
func (d *DepositeInteresetUseCaseTestSuite) TestPagging_Failed() {
	d.dirm.On("Pagging", mock.MockPageReq).Return(nil, dto.Paging{}, errors.New("failed"))
	ahc, p, err := d.diuc.Pagging(mock.MockPageReq)
	assert.Error(d.T(), err)
	assert.Equal(d.T(), 0, len(ahc))
	assert.Equal(d.T(), dto.Paging{}, p)
}
func (d *DepositeInteresetUseCaseTestSuite) TestPagging_DefaultPage() {
	d.dirm.On("Pagging", mock.MockPageReq).Return(mock.MockDeposites, mock.MockPaging, nil)
	ahc, p, err := d.diuc.Pagging(dto.PageRequest{Page: -1, Size: 5})
	assert.Nil(d.T(), err)
	assert.Equal(d.T(), len(mock.MockDeposites), len(ahc))
	assert.Equal(d.T(), mock.MockPaging.TotalRows, p.TotalRows)
}
