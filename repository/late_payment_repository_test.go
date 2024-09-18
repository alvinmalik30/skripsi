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

type LatePaymentRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    LatePaymentFee
}

func (l *LatePaymentRepoTestSuite) SetupTest() {
	d, s, err := sqlmock.New()
	assert.NoError(l.T(), err)
	l.mockDB = d
	l.mockSQL = s
	l.repo = NewLatePaymentFeeRepository(l.mockDB)
}

func TestLatePaymentRepoTestSuite(t *testing.T) {
	suite.Run(t, new(LatePaymentRepoTestSuite))
}
func (l *LatePaymentRepoTestSuite) TestUpdate_Success() {
	expectedQuery := `UPDATE late_payment_fee SET name = $2, nominal = $3, unit =$4 WHERE id = $1;`
	l.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnResult(sqlmock.NewResult(1, 1))
	err := l.repo.Update(mock.MockLatePF)
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LatePaymentRepoTestSuite) TestUpdate_Failed() {
	expectedQuery := `UPDATE late_payment_fee SET name = $2, nominal = $3, unit =$4 WHERE id = $1;`
	l.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error update"))
	err := l.repo.Update(mock.MockLatePF)
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LatePaymentRepoTestSuite) TestFindById_Success() {
	r := sqlmock.NewRows([]string{"id", "name", "nominal", "unit"})
	r.AddRow(mock.MockLatePF.Id, mock.MockLatePF.Name, mock.MockLatePF.Nominal, mock.MockLatePF.Unit)
	expectedQuery := `SELECT id, name, nominal, unit FROM late_payment_fee WHERE id =$1`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(r)
	ahc, err := l.repo.FindById(mock.MockLatePF.Id)
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), mock.MockLatePF.Id, ahc.Id)
}
func (l *LatePaymentRepoTestSuite) TestFindById_Failed() {
	expectedQuery := `SELECT id, name, nominal, unit FROM late_payment_fee WHERE id =$1`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error"))
	ahc, err := l.repo.FindById(mock.MockLatePF.Id)
	assert.Error(l.T(), err)
	assert.Equal(l.T(), model.LatePaymentFee{}, ahc)
}
func (l *LatePaymentRepoTestSuite) TestDeleteById_Success() {
	expectedQuery := `DELETE FROM late_payment_fee WHERE id = $1`
	l.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnResult(sqlmock.NewResult(1, 1))
	err := l.repo.DeleteById(mock.MockLatePF.Id)
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LatePaymentRepoTestSuite) TestDeleteById_Failed() {
	expectedQuery := `DELETE FROM late_payment_fee WHERE id = $1`
	l.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error delete"))
	err := l.repo.DeleteById(mock.MockLatePF.Id)
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LatePaymentRepoTestSuite) TestCreateNew_Success() {
	expectedQuery := `INSERT INTO late_payment_fee`
	l.mockSQL.ExpectExec(expectedQuery).WillReturnResult(sqlmock.NewResult(1, 1))
	err := l.repo.CreateNew(mock.MockLatePF)
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LatePaymentRepoTestSuite) TestCreateNew_Failed() {
	expectedQuery := `INSERT INTO late_payment_fee`
	l.mockSQL.ExpectExec(expectedQuery).WillReturnError(errors.New("error crated"))
	err := l.repo.CreateNew(mock.MockLatePF)
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LatePaymentRepoTestSuite) TestPagging_Success() {
	rows := sqlmock.NewRows([]string{"id", "name", "nominal", "unit"})
	for _, pf := range mock.MockLatePFDatas {
		rows.AddRow(pf.Id, pf.Name, pf.Nominal, pf.Unit)
	}
	expectedQuery := `SELECT id, name, nominal, unit FROM late_payment_fee LIMIT $2 OFFSET $1;`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM late_payment_fee`)).WillReturnRows(rowCount)

	uc, p, err := l.repo.Pagging(mock.MockPageReq)
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), 1, len(uc))
	assert.Equal(l.T(), 1, p.TotalRows)
}
func (l *LatePaymentRepoTestSuite) TestPagging_Fail() {
	// error select paging
	expectedQuery := `SELECT id, name, nominal, unit FROM late_payment_fee LIMIT $2 OFFSET $1;`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnError(errors.New("failed"))
	uc, p, err := l.repo.Pagging(dto.PageRequest{})
	assert.Error(l.T(), err)
	assert.Nil(l.T(), uc)
	assert.Equal(l.T(), 0, p.TotalRows)
	// error select count
	rows := sqlmock.NewRows([]string{"id", "name", "nominal", "unit"})
	for _, pf := range mock.MockLatePFDatas {
		rows.AddRow(pf.Id, pf.Name, pf.Nominal, pf.Unit)
	}
	expectedQuery = `SELECT id, name, nominal, unit FROM late_payment_fee LIMIT $2 OFFSET $1;`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM late_payment_fee`)).WillReturnError(errors.New("failed"))
	uc, p, err = l.repo.Pagging(mock.MockPageReq)
	assert.Error(l.T(), err)
	assert.Nil(l.T(), uc)
	assert.Equal(l.T(), 0, p.TotalRows)
}
