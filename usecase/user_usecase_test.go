package usecase

import (
	"errors"
	"polen/mock"
	"polen/mock/repomock"
	"polen/model"
	"polen/model/dto"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserUseCaseTestSuite struct {
	suite.Suite
	urm *repomock.UserRepoMock
	uuc UserUseCase
	ctx *gin.Context
}

func (u *UserUseCaseTestSuite) SetupTest() {
	u.urm = new(repomock.UserRepoMock)
	u.ctx = &gin.Context{}
	u.uuc = NewUserUseCase(u.urm, u.ctx)
}

func TestUserUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UserUseCaseTestSuite))
}

func (u *UserUseCaseTestSuite) TestRegister_EmptyInvalid() {
	// username required
	u.urm.On("Save", model.UserCredential{Id: "1"}).Return(errors.New("username required"))
	err := u.uuc.Register(dto.AuthRequest{})
	assert.Error(u.T(), err)
	// password required
	u.urm.On("Save", model.UserCredential{Id: "1", Username: "akbar", Email: "akbar@gmail.com"}).Return(errors.New("password required"))
	err = u.uuc.Register(dto.AuthRequest{Username: "akbar", Email: "akbar@gmail.com"})
	assert.Error(u.T(), err)
	// email required
	u.urm.On("Save", model.UserCredential{Id: "1", Username: "akbar", Password: "123"}).Return(errors.New("email required"))
	err = u.uuc.Register(dto.AuthRequest{Username: "akbar", Password: "123"})
	assert.Error(u.T(), err)
	// role required
	u.urm.On("Save", model.UserCredential{Id: "1", Username: "akbar", Email: "akbar@gmail.com", Password: "123"}).Return(errors.New("role is required"))
	err = u.uuc.Register(dto.AuthRequest{Username: "akbar", Email: "akbar@gmail.com", Password: "123"})
	assert.Error(u.T(), err)
}
func (u *UserUseCaseTestSuite) TestRegisterCheckRole_Fail() {
	u.urm.On("Save", mock.MockUserCred).Return(errors.New("role you has choose isnt available"))
	err := u.uuc.Register(mock.MockAuthReq)
	assert.Error(u.T(), err)
}
func (u *UserUseCaseTestSuite) TestRegister_InvalidEmail() {
	mockAuthReq := dto.AuthRequest{
		Username: "akbar",
		Email:    "akbar@gmail",
		Password: "123",
		Role:     "peminjam",
	}
	u.urm.On("Save", mock.MockUserCred).Return(errors.New("is not valid email"))
	err := u.uuc.Register(mockAuthReq)
	assert.Error(u.T(), err)
}
func (u *UserUseCaseTestSuite) TestFindById_Success() {
	u.urm.On("FindById", mock.MockUserCred.Id).Return(mock.MockUserCred, nil)
	uc, err := u.uuc.FindById(mock.MockUserCred.Id)
	assert.Nil(u.T(), err)
	assert.NoError(u.T(), err)
	assert.Equal(u.T(), mock.MockUserCred, uc)
}
func (u *UserUseCaseTestSuite) TestFindById_Fail() {
	u.urm.On("FindById", mock.MockUserCred.Id).Return(model.UserCredential{}, errors.New("error"))
	uc, err := u.uuc.FindById(mock.MockUserCred.Id)
	assert.Error(u.T(), err)
	assert.NotNil(u.T(), err)
	assert.Equal(u.T(), model.UserCredential{}, uc)
}
