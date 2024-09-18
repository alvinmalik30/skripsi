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

type LatePaymentUseCaseTestSuite struct {
	suite.Suite
	lprm  *repomock.LatePaymentRepoMock
	lpfuc LatePaymentFeeUsecase
}

func (l *LatePaymentUseCaseTestSuite) SetupTest() {
	l.lprm = new(repomock.LatePaymentRepoMock)
	l.lpfuc = NewLatePaymentFeeUseCase(l.lprm)
}

func TestLatePaymentUseCaseTestSuite(t *testing.T) {
	suite.Run(t, new(LatePaymentUseCaseTestSuite))
}
func (l *LatePaymentUseCaseTestSuite) TestPagging_Success() {
	l.lprm.On("Pagging", mock.MockPageReq).Return(mock.MockLatePFDatas, mock.MockPaging, nil)
	ahc, p, err := l.lpfuc.Pagging(mock.MockPageReq)
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), len(mock.MockLatePFDatas), len(ahc))
	assert.Equal(l.T(), mock.MockPaging.TotalRows, p.TotalRows)
}
func (l *LatePaymentUseCaseTestSuite) TestPagging_Failed() {
	l.lprm.On("Pagging", mock.MockPageReq).Return(nil, dto.Paging{}, errors.New("failed"))
	ahc, p, err := l.lpfuc.Pagging(mock.MockPageReq)
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 0, len(ahc))
	assert.Equal(l.T(), dto.Paging{}, p)
}
func (l *LatePaymentUseCaseTestSuite) TestDeleteById_Success() {
	l.lprm.On("FindById", "1").Return(mock.MockLatePF, nil)
	l.lprm.On("DeleteById", "1").Return(nil)
	err := l.lpfuc.DeleteById("1")
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LatePaymentUseCaseTestSuite) TestDeleteById_Failed() {
	// failed find id
	l.lprm.On("FindById", "1").Return(model.LatePaymentFee{}, errors.New("failed id"))
	err := l.lpfuc.DeleteById("1")
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LatePaymentUseCaseTestSuite) TestDeleteById_InvalidServer() {
	// failed server error
	l.lprm.On("FindById", "1").Return(mock.MockLatePF, nil)
	l.lprm.On("DeleteById", "1").Return(errors.New("failed to delete app handling cost"))
	err := l.lpfuc.DeleteById("1")
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LatePaymentUseCaseTestSuite) TestFindById_Success() {
	l.lprm.On("FindById", "1").Return(mock.MockLatePF, nil)
	ahc, err := l.lpfuc.FindById("1")
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), mock.MockLatePF.Id, ahc.Id)
}
func (l *LatePaymentUseCaseTestSuite) TestFindById_Failed() {
	l.lprm.On("FindById", "1").Return(model.LatePaymentFee{}, errors.New("failed"))
	ahc, err := l.lpfuc.FindById("1")
	assert.Error(l.T(), err)
	assert.Equal(l.T(), model.LatePaymentFee{}, ahc)
}
func (l *LatePaymentUseCaseTestSuite) TestCreateNew_Success() {
	l.lprm.On("CreateNew", mock.MockLatePF).Return(nil)
	i, err := l.lpfuc.CreateNew(mock.MockLatePF)
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), 201, i)
}
func (l *LatePaymentUseCaseTestSuite) TestCreateNew_Failed() {
	// id required
	l.lprm.On("CreateNew", model.LatePaymentFee{}).Return(errors.New("id is required"))
	i, err := l.lpfuc.CreateNew(model.LatePaymentFee{})
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 400, i)
	// name required
	l.lprm.On("CreateNew", model.LatePaymentFee{Id: "1"}).Return(errors.New("name is required"))
	i, err = l.lpfuc.CreateNew(model.LatePaymentFee{Id: "1"})
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 400, i)
	// nominal required
	l.lprm.On("CreateNew", model.LatePaymentFee{Id: "1", Name: "akbar"}).Return(errors.New("nominal is required"))
	i, err = l.lpfuc.CreateNew(model.LatePaymentFee{Id: "1", Name: "akbar"})
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 400, i)
	// unit required
	l.lprm.On("CreateNew", model.LatePaymentFee{Id: "1", Name: "akbar", Nominal: 100000}).Return(errors.New("unit is required"))
	i, err = l.lpfuc.CreateNew(model.LatePaymentFee{Id: "1", Name: "akbar", Nominal: 100000})
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 400, i)
	// rupiah or percent are units required
	l.lprm.On("CreateNew", model.LatePaymentFee{Id: "1", Name: "akbar", Nominal: 100000, Unit: "USD"}).Return(errors.New("unit must be rupiah or percent"))
	i, err = l.lpfuc.CreateNew(model.LatePaymentFee{Id: "1", Name: "akbar", Nominal: 100000, Unit: "USD"})
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 400, i)
	// server error
	l.lprm.On("CreateNew", mock.MockLatePF).Return(errors.New("server error"))
	i, err = l.lpfuc.CreateNew(mock.MockLatePF)
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 500, i)
}
func (l *LatePaymentUseCaseTestSuite) TestUpdate_Success() {
	l.lprm.On("FindById", "1").Return(mock.MockLatePF, nil)
	l.lprm.On("Update", mock.MockLatePF).Return(nil)
	err := l.lpfuc.Update(mock.MockLatePF)
	assert.Nil(l.T(), err)
}
func (l *LatePaymentUseCaseTestSuite) TestUpdate_Failed() {
	// id required
	l.lprm.On("Update", model.LatePaymentFee{}).Return(errors.New("id is required"))
	err := l.lpfuc.Update(model.LatePaymentFee{})
	assert.Error(l.T(), err)
	// id db failed
	l.lprm.On("FindById", "1").Return(model.LatePaymentFee{}, errors.New("failed"))
	l.lprm.On("Update", mock.MockLatePF).Return(errors.New("error"))
	err = l.lpfuc.Update(mock.MockLatePF)
	assert.Error(l.T(), err)
}
func (l *LatePaymentUseCaseTestSuite) TestUpdate_UnitPayload() {
	l.lprm.On("FindById", "1").Return(mock.MockLatePF, nil)
	l.lprm.On("Update", model.LatePaymentFee{Id: "1", Name: "akbar", Nominal: 100000, Unit: "USD"}).Return(errors.New("unit must be rupiah or percent"))
	err := l.lpfuc.Update(model.LatePaymentFee{Id: "1", Name: "akbar", Nominal: 100000, Unit: "USD"})
	assert.Error(l.T(), err)
}
func (l *LatePaymentUseCaseTestSuite) TestUpdate_InvalidServer() {
	l.lprm.On("FindById", "1").Return(mock.MockLatePF, nil)
	l.lprm.On("Update", mock.MockLatePF).Return(errors.New("server error"))
	err := l.lpfuc.Update(mock.MockLatePF)
	assert.Error(l.T(), err)
}
