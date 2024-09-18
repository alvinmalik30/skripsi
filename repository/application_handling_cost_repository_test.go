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

type ApplicationHCRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    AppHandlingCost
}

func (a *ApplicationHCRepoTestSuite) SetupTest() {
	db, s, err := sqlmock.New()
	assert.NoError(a.T(), err)
	a.mockDB = db
	a.mockSQL = s
	a.repo = NewAppHandlingCostRepository(a.mockDB)
}

func TestApplicationHCRepoTestSuite(t *testing.T) {
	suite.Run(t, new(ApplicationHCRepoTestSuite))
}
func (a *ApplicationHCRepoTestSuite) TestUpdate_Success() {
	expectedQuery := `UPDATE application_handling_cost SET name = $2, nominal = $3, unit =$4 WHERE id = $1;`
	a.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnResult(sqlmock.NewResult(1, 1))
	err := a.repo.Update(mock.MockAppHC)
	assert.Nil(a.T(), err)
	assert.NoError(a.T(), err)
}
func (a *ApplicationHCRepoTestSuite) TestUpdate_Failed() {
	expectedQuery := `UPDATE application_handling_cost SET name = $2, nominal = $3, unit =$4 WHERE id = $1;`
	a.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error update"))
	err := a.repo.Update(mock.MockAppHC)
	assert.Error(a.T(), err)
	assert.NotNil(a.T(), err)
}
func (a *ApplicationHCRepoTestSuite) TestFindById_Success() {
	r := sqlmock.NewRows([]string{"id", "name", "nominal", "unit"})
	r.AddRow(mock.MockAppHC.Id, mock.MockAppHC.Name, mock.MockAppHC.Nominal, mock.MockAppHC.Unit)
	expectedQuery := `SELECT id, name, nominal, unit FROM application_handling_cost WHERE id = $1`
	a.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(r)
	ahc, err := a.repo.FindById(mock.MockAppHC.Id)
	assert.Nil(a.T(), err)
	assert.Equal(a.T(), mock.MockAppHC.Id, ahc.Id)
}
func (a *ApplicationHCRepoTestSuite) TestFindById_Failed() {
	expectedQuery := `SELECT id, name, nominal, unit FROM application_handling_cost WHERE id = $1`
	a.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error"))
	ahc, err := a.repo.FindById(mock.MockAppHC.Id)
	assert.Error(a.T(), err)
	assert.Equal(a.T(), model.AppHandlingCost{}, ahc)
}
func (a *ApplicationHCRepoTestSuite) TestDeleteById_Success() {
	expectedQuery := `DELETE FROM application_handling_cost WHERE id = $1`
	a.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnResult(sqlmock.NewResult(1, 1))
	err := a.repo.DeleteById(mock.MockAppHC.Id)
	assert.Nil(a.T(), err)
	assert.NoError(a.T(), err)
}
func (a *ApplicationHCRepoTestSuite) TestDeleteById_Failed() {
	expectedQuery := `DELETE FROM application_handling_cost WHERE id = $1`
	a.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error delete"))
	err := a.repo.DeleteById(mock.MockAppHC.Id)
	assert.Error(a.T(), err)
	assert.NotNil(a.T(), err)
}
func (a *ApplicationHCRepoTestSuite) TestCreateNew_Success() {
	expectedQuery := `INSERT INTO application_handling_cost`
	a.mockSQL.ExpectExec(expectedQuery).WillReturnResult(sqlmock.NewResult(1, 1))
	err := a.repo.CreateNew(mock.MockAppHC)
	assert.Nil(a.T(), err)
	assert.NoError(a.T(), err)
}
func (a *ApplicationHCRepoTestSuite) TestCreateNew_Failed() {
	expectedQuery := `INSERT INTO application_handling_cost`
	a.mockSQL.ExpectExec(expectedQuery).WillReturnError(errors.New("error crated"))
	err := a.repo.CreateNew(mock.MockAppHC)
	assert.Error(a.T(), err)
	assert.NotNil(a.T(), err)
}
func (a *ApplicationHCRepoTestSuite) TestPagging_Success() {
	rows := sqlmock.NewRows([]string{"id", "name", "nominal", "unit"})
	for _, hc := range mock.MockAppHCDatas {
		rows.AddRow(hc.Id, hc.Name, hc.Nominal, hc.Unit)
	}
	expectedQuery := `SELECT id, name, nominal, unit FROM application_handling_cost LIMIT $2 OFFSET $1;`
	a.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	a.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM application_handling_cost`)).WillReturnRows(rowCount)

	uc, p, err := a.repo.Pagging(mock.MockPageReq)
	assert.Nil(a.T(), err)
	assert.Equal(a.T(), 1, len(uc))
	assert.Equal(a.T(), 1, p.TotalRows)
}
func (a *ApplicationHCRepoTestSuite) TestPagging_Fail() {
	// error select paging
	expectedQuery := `SELECT id, name, nominal, unit FROM application_handling_cost LIMIT $2 OFFSET $1;`
	a.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnError(errors.New("failed"))
	uc, p, err := a.repo.Pagging(dto.PageRequest{})
	assert.Error(a.T(), err)
	assert.Nil(a.T(), uc)
	assert.Equal(a.T(), 0, p.TotalRows)
	// error select count
	rows := sqlmock.NewRows([]string{"id", "name", "nominal", "unit"})
	for _, hc := range mock.MockAppHCDatas {
		rows.AddRow(hc.Id, hc.Name, hc.Nominal, hc.Unit)
	}
	expectedQuery = `SELECT id, name, nominal, unit FROM application_handling_cost LIMIT $2 OFFSET $1;`
	a.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	a.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM application_handling_cost`)).WillReturnError(errors.New("failed"))
	uc, p, err = a.repo.Pagging(mock.MockPageReq)
	assert.Error(a.T(), err)
	assert.Nil(a.T(), uc)
	assert.Equal(a.T(), 0, p.TotalRows)
}
