package repository

import (
	"database/sql"
	"errors"
	"polen/mock"
	"polen/model"
	"polen/model/dto"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type UserRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    UserRepository
}

func (u *UserRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(u.T(), err)
	u.mockDB = db
	u.mockSQL = mock
	u.repo = NewUserRepository(u.mockDB)
}

func TestUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(UserRepoTestSuite))
}

func (u *UserRepoTestSuite) TestFindByUsername_Success() {
	rows := sqlmock.NewRows([]string{"id", "username", "role", "password"})
	rows.AddRow(mock.MockUserCred.Id, mock.MockUserCred.Username, mock.MockUserCred.Role, mock.MockUserCred.Password)
	expectedSQL := `SELECT id, username, role, password FROM user_credential WHERE username = $1`
	u.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockUserCred.Username).WillReturnRows(rows)
	uc, err := u.repo.FindByUsername(mock.MockUserCred.Username)
	assert.Nil(u.T(), err)
	assert.NoError(u.T(), err)
	assert.Equal(u.T(), mock.MockUserCred.Username, uc.Username)
}
func (u *UserRepoTestSuite) TestFindByUsername_Fail() {
	expectedSQL := `SELECT id, username, role, password FROM user_credential WHERE username = $1`
	u.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs("ismail").WillReturnError(errors.New("error"))
	uc, err := u.repo.FindByUsername("ismail")
	assert.Error(u.T(), err)
	assert.NotNil(u.T(), err)
	assert.Equal(u.T(), model.UserCredential{}, uc)
}

func (u *UserRepoTestSuite) TestSave_Success() {
	u.mockSQL.ExpectBegin()
	u.mockSQL.ExpectExec(`INSERT INTO user_credential`).WillReturnResult(sqlmock.NewResult(1, 1))
	u.mockSQL.ExpectExec(`INSERT INTO biodata`).WillReturnResult(sqlmock.NewResult(1, 1))
	u.mockSQL.ExpectCommit()
	err := u.repo.Save(mock.MockUserCred, mock.MockBiodata.Id)
	assert.Nil(u.T(), err)
	assert.NoError(u.T(), err)
}
func (u *UserRepoTestSuite) TestSave_Failed() {
	// Begin failed
	u.mockSQL.ExpectBegin().WillReturnError(errors.New("begin failed"))
	err := u.repo.Save(mock.MockUserCred, mock.MockBiodata.Id)
	assert.Error(u.T(), err)
	// Insert user_credential failed
	u.mockSQL.ExpectBegin()
	u.mockSQL.ExpectExec(`INSERT INTO user_credential`).WillReturnError(errors.New("insert failed"))
	err = u.repo.Save(mock.MockUserCred, mock.MockBiodata.Id)
	assert.Error(u.T(), err)
	// Insert biodata failed
	u.mockSQL.ExpectBegin()
	u.mockSQL.ExpectExec(`INSERT INTO user_credential`).WillReturnResult(sqlmock.NewResult(1, 1))
	u.mockSQL.ExpectExec(`INSERT INTO biodata`).WillReturnError(errors.New("insert failed"))
	err = u.repo.Save(mock.MockUserCred, mock.MockBiodata.Id)
	assert.Error(u.T(), err)
}
func (u *UserRepoTestSuite) TestFindById_Success() {
	rows := sqlmock.NewRows([]string{"id", "username", "email", "role", "virtual_account_number", "is_active"})
	rows.AddRow(mock.MockUserCred.Id, mock.MockUserCred.Username, mock.MockUserCred.Email, mock.MockUserCred.Role, mock.MockUserCred.VANumber, mock.MockUserCred.IsActive)
	expectedSQL := `SELECT id, username, email, role, virtual_account_number, is_active FROM user_credential WHERE id =$1`
	u.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockUserCred.Id).WillReturnRows(rows)
	uc, err := u.repo.FindById(mock.MockUserCred.Id)
	assert.Nil(u.T(), err)
	assert.NoError(u.T(), err)
	assert.Equal(u.T(), mock.MockUserCred.Id, uc.Id)
}
func (u *UserRepoTestSuite) TestFindById_Fail() {
	expectedSQL := `SELECT id, username, email, role, virtual_account_number, is_active FROM user_credential WHERE id =$1`
	u.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockUserCred.Id).WillReturnError(errors.New("error"))
	uc, err := u.repo.FindById(mock.MockUserCred.Id)
	assert.Error(u.T(), err)
	assert.NotNil(u.T(), err)
	assert.Equal(u.T(), model.UserCredential{}, uc)
}
func (u *UserRepoTestSuite) TestSaldo_Success() {
	u.mockSQL.ExpectBegin()
	u.mockSQL.ExpectExec(`INSERT INTO user_credential`).WillReturnResult(sqlmock.NewResult(1, 1))
	u.mockSQL.ExpectExec(`INSERT INTO biodata`).WillReturnResult(sqlmock.NewResult(1, 1))
	u.mockSQL.ExpectExec(`INSERT INTO saldo`).WillReturnResult(sqlmock.NewResult(1, 1))
	u.mockSQL.ExpectCommit()
	err := u.repo.Saldo(mock.MockUserCred, mock.MockSaldo.Id, mock.MockBiodata.Id)
	assert.Nil(u.T(), err)
	assert.NoError(u.T(), err)
}
func (u *UserRepoTestSuite) TestSaldo_Fail() {
	// Begin Failed
	u.mockSQL.ExpectBegin().WillReturnError(errors.New("begin failed"))
	err := u.repo.Saldo(mock.MockUserCred, mock.MockSaldo.Id, mock.MockBiodata.Id)
	assert.Error(u.T(), err)
	// User_Credential insert failed
	u.mockSQL.ExpectBegin()
	u.mockSQL.ExpectExec(`INSERT INTO user_credential`).WillReturnError(errors.New("insert failed"))
	err = u.repo.Saldo(mock.MockUserCred, mock.MockSaldo.Id, mock.MockBiodata.Id)
	assert.Error(u.T(), err)
	// Biodata insert failed
	u.mockSQL.ExpectBegin()
	u.mockSQL.ExpectExec(`INSERT INTO user_credential`).WillReturnResult(sqlmock.NewResult(1, 1))
	u.mockSQL.ExpectExec(`INSERT INTO biodata`).WillReturnError(errors.New("insert failed"))
	err = u.repo.Saldo(mock.MockUserCred, mock.MockSaldo.Id, mock.MockBiodata.Id)
	assert.Error(u.T(), err)
	// Saldo insert failed
	u.mockSQL.ExpectBegin()
	u.mockSQL.ExpectExec(`INSERT INTO user_credential`).WillReturnResult(sqlmock.NewResult(1, 1))
	u.mockSQL.ExpectExec(`INSERT INTO biodata`).WillReturnResult(sqlmock.NewResult(1, 1))
	u.mockSQL.ExpectExec(`INSERT INTO saldo`).WillReturnError(errors.New("insert failed"))
	err = u.repo.Saldo(mock.MockUserCred, mock.MockSaldo.Id, mock.MockBiodata.Id)
	assert.Error(u.T(), err)
}
func (u *UserRepoTestSuite) TestPagging_Success() {
	rows := sqlmock.NewRows([]string{"id", "username", "email", "role", "virtual_account_number", "is_active"})
	for _, userCred := range mock.MockUserCreds {
		rows.AddRow(userCred.Id, userCred.Username, userCred.Email, userCred.Role, userCred.VANumber, userCred.IsActive)
	}
	expectedSQL := `SELECT id, username, email, role, virtual_account_number, is_active FROM user_credential LIMIT $2 OFFSET $1`
	u.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	u.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM user_credential`)).WillReturnRows(rowCount)

	uc, p, err := u.repo.Pagging(mock.MockPageReq)
	assert.Nil(u.T(), err)
	assert.Equal(u.T(), 1, len(uc))
	assert.Equal(u.T(), 1, p.TotalRows)
}
func (u *UserRepoTestSuite) TestPagging_Fail() {
	// error select paging
	expectedSQL := `SELECT id, username, email, role, virtual_account_number, is_active FROM user_credential LIMIT $2 OFFSET $1`
	u.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("failed"))
	uc, p, err := u.repo.Pagging(dto.PageRequest{})
	assert.Error(u.T(), err)
	assert.Nil(u.T(), uc)
	assert.Equal(u.T(), 0, p.TotalRows)
	// error select count
	rows := sqlmock.NewRows([]string{"id", "username", "email", "role", "virtual_account_number", "is_active"})
	for _, userCred := range mock.MockUserCreds {
		rows.AddRow(userCred.Id, userCred.Username, userCred.Email, userCred.Role, userCred.VANumber, userCred.IsActive)
	}
	expectedSQL = `SELECT id, username, email, role, virtual_account_number, is_active FROM user_credential LIMIT $2 OFFSET $1`
	u.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)
	u.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM user_credential`)).WillReturnError(errors.New("failed"))
	uc, p, err = u.repo.Pagging(mock.MockPageReq)
	assert.Error(u.T(), err)
	assert.Nil(u.T(), uc)
	assert.Equal(u.T(), 0, p.TotalRows)
}
