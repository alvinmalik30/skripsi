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

type ApplicationHCUseCaseTestSuite struct {
	suite.Suite
	ahcrm *repomock.ApplicationHCRepoMock
	ahcuc AppHandlingCostUsecase
}

func (a *ApplicationHCUseCaseTestSuite) SetupTest() {
	a.ahcrm = new(repomock.ApplicationHCRepoMock)
	a.ahcuc = NewAppHandlingCostUseCase(a.ahcrm)
}

func TestApplicationHCUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(ApplicationHCUseCaseTestSuite))
}

func (a *ApplicationHCUseCaseTestSuite) TestPagging_Success() {
	a.ahcrm.On("Pagging", mock.MockPageReq).Return(mock.MockAppHCDatas, mock.MockPaging, nil)
	ahc, p, err := a.ahcuc.Pagging(mock.MockPageReq)
	assert.Nil(a.T(), err)
	assert.Equal(a.T(), len(mock.MockAppHCDatas), len(ahc))
	assert.Equal(a.T(), mock.MockPaging.TotalRows, p.TotalRows)
}
func (a *ApplicationHCUseCaseTestSuite) TestPagging_Failed() {
	a.ahcrm.On("Pagging", mock.MockPageReq).Return(nil, dto.Paging{}, errors.New("failed"))
	ahc, p, err := a.ahcuc.Pagging(mock.MockPageReq)
	assert.Error(a.T(), err)
	assert.Equal(a.T(), 0, len(ahc))
	assert.Equal(a.T(), dto.Paging{}, p)
}
func (a *ApplicationHCUseCaseTestSuite) TestFindById_Success() {
	a.ahcrm.On("FindById", mock.MockAppHC.Id).Return(mock.MockAppHC, nil)
	ahc, err := a.ahcuc.FindById(mock.MockAppHC.Id)
	assert.Nil(a.T(), err)
	assert.Equal(a.T(), mock.MockAppHC.Id, ahc.Id)
}
func (a *ApplicationHCUseCaseTestSuite) TestFindById_Failed() {
	a.ahcrm.On("FindById", mock.MockAppHC.Id).Return(model.AppHandlingCost{}, errors.New("failed"))
	ahc, err := a.ahcuc.FindById(mock.MockAppHC.Id)
	assert.Error(a.T(), err)
	assert.Equal(a.T(), model.AppHandlingCost{}, ahc)
}
func (a *ApplicationHCUseCaseTestSuite) TestDeleteById_Success() {
	a.ahcrm.On("FindById", mock.MockAppHC.Id).Return(mock.MockAppHC, nil)
	a.ahcrm.On("DeleteById", mock.MockAppHC.Id).Return(nil)
	err := a.ahcuc.DeleteById("1")
	assert.Nil(a.T(), err)
	assert.NoError(a.T(), err)
}
func (a *ApplicationHCUseCaseTestSuite) TestDeleteById_Failed() {
	// failed find id
	a.ahcrm.On("FindById", mock.MockAppHC.Id).Return(model.AppHandlingCost{}, errors.New("failed id"))
	err := a.ahcuc.DeleteById("1")
	assert.Error(a.T(), err)
	assert.NotNil(a.T(), err)
}
func (a *ApplicationHCUseCaseTestSuite) TestDeleteById_InvalidServer() {
	// failed server error
	a.ahcrm.On("FindById", mock.MockAppHC.Id).Return(mock.MockAppHC, nil)
	a.ahcrm.On("DeleteById", mock.MockAppHC.Id).Return(errors.New("failed to delete app handling cost"))
	err := a.ahcuc.DeleteById(mock.MockAppHC.Id)
	assert.Error(a.T(), err)
	assert.NotNil(a.T(), err)
}
func (a *ApplicationHCUseCaseTestSuite) TestCreateNew_Success() {
	a.ahcrm.On("CreateNew", mock.MockAppHC).Return(nil)
	i, err := a.ahcuc.CreateNew(mock.MockAppHC)
	assert.Nil(a.T(), err)
	assert.Equal(a.T(), 201, i)
}
func (a *ApplicationHCUseCaseTestSuite) TestCreateNew_Failed() {
	// id required
	a.ahcrm.On("CreateNew", model.AppHandlingCost{}).Return(errors.New("id is required"))
	i, err := a.ahcuc.CreateNew(model.AppHandlingCost{})
	assert.Error(a.T(), err)
	assert.Equal(a.T(), 400, i)
	// name required
	a.ahcrm.On("CreateNew", model.AppHandlingCost{Id: "1"}).Return(errors.New("name is required"))
	i, err = a.ahcuc.CreateNew(model.AppHandlingCost{Id: "1"})
	assert.Error(a.T(), err)
	assert.Equal(a.T(), 400, i)
	// nominal required
	a.ahcrm.On("CreateNew", model.AppHandlingCost{Id: "1", Name: "akbar"}).Return(errors.New("nominal is required"))
	i, err = a.ahcuc.CreateNew(model.AppHandlingCost{Id: "1", Name: "akbar"})
	assert.Error(a.T(), err)
	assert.Equal(a.T(), 400, i)
	// unit required
	a.ahcrm.On("CreateNew", model.AppHandlingCost{Id: "1", Name: "akbar", Nominal: 100000}).Return(errors.New("unit is required"))
	i, err = a.ahcuc.CreateNew(model.AppHandlingCost{Id: "1", Name: "akbar", Nominal: 100000})
	assert.Error(a.T(), err)
	assert.Equal(a.T(), 400, i)
	// rupiah or percent are units required
	a.ahcrm.On("CreateNew", model.AppHandlingCost{Id: "1", Name: "akbar", Nominal: 100000, Unit: "USD"}).Return(errors.New("unit must be rupiah or percent"))
	i, err = a.ahcuc.CreateNew(model.AppHandlingCost{Id: "1", Name: "akbar", Nominal: 100000, Unit: "USD"})
	assert.Error(a.T(), err)
	assert.Equal(a.T(), 400, i)
	// server error
	a.ahcrm.On("CreateNew", mock.MockAppHC).Return(errors.New("server error"))
	i, err = a.ahcuc.CreateNew(mock.MockAppHC)
	assert.Error(a.T(), err)
	assert.Equal(a.T(), 500, i)
}
func (a *ApplicationHCUseCaseTestSuite) TestUpdate_Success() {
	a.ahcrm.On("FindById", mock.MockAppHC.Id).Return(mock.MockAppHC, nil)
	a.ahcrm.On("Update", mock.MockAppHC).Return(nil)
	err := a.ahcuc.Update(mock.MockAppHC)
	assert.Nil(a.T(), err)
}
func (a *ApplicationHCUseCaseTestSuite) TestUpdate_Failed() {
	// id required
	a.ahcrm.On("Update", model.AppHandlingCost{}).Return(errors.New("id is required"))
	err := a.ahcuc.Update(model.AppHandlingCost{})
	assert.Error(a.T(), err)
	// id db failed
	a.ahcrm.On("FindById", mock.MockAppHC.Id).Return(model.AppHandlingCost{}, errors.New("failed"))
	a.ahcrm.On("Update", mock.MockAppHC).Return(errors.New("error"))
	err = a.ahcuc.Update(mock.MockAppHC)
	assert.Error(a.T(), err)
}
func (a *ApplicationHCUseCaseTestSuite) TestUpdate_UnitPayload() {
	a.ahcrm.On("FindById", mock.MockAppHC.Id).Return(mock.MockAppHC, nil)
	a.ahcrm.On("Update", model.AppHandlingCost{Id: "1", Name: "akbar", Nominal: 100000, Unit: "USD"}).Return(errors.New("unit must be rupiah or percent"))
	err := a.ahcuc.Update(model.AppHandlingCost{Id: "1", Name: "akbar", Nominal: 100000, Unit: "USD"})
	assert.Error(a.T(), err)
}
func (a *ApplicationHCUseCaseTestSuite) TestUpdate_InvalidServer() {
	a.ahcrm.On("FindById", mock.MockAppHC.Id).Return(mock.MockAppHC, nil)
	a.ahcrm.On("Update", mock.MockAppHC).Return(errors.New("server error"))
	err := a.ahcuc.Update(mock.MockAppHC)
	assert.Error(a.T(), err)
}
