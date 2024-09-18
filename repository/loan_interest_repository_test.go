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

type LoanInterestRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    LoanInterest
}

func (l *LoanInterestRepoTestSuite) SetupTest() {
	d, s, err := sqlmock.New()
	assert.NoError(l.T(), err)
	l.mockDB = d
	l.mockSQL = s
	l.repo = NewLoanInterestRepository(l.mockDB)
}

func TestLoanInterestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(LoanInterestRepoTestSuite))
}
func (l *LoanInterestRepoTestSuite) TestUpdate_Success() {
	expectedQuery := `UPDATE loan_interest SET duration_months = $2, loan_interest_rate = $3 WHERE id = $1;`
	l.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnResult(sqlmock.NewResult(1, 1))
	err := l.repo.Update(mock.MockLoanInterest)
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LoanInterestRepoTestSuite) TestUpdate_Failed() {
	expectedQuery := `UPDATE loan_interest SET duration_months = $2, loan_interest_rate = $3 WHERE id = $1;`
	l.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error update"))
	err := l.repo.Update(mock.MockLoanInterest)
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LoanInterestRepoTestSuite) TestFindById_Success() {
	r := sqlmock.NewRows([]string{"id", "duration_months", "loan_interest_rate"})
	r.AddRow(mock.MockLoanInterest.Id, mock.MockLoanInterest.DurationMonths, mock.MockLoanInterest.LoanInterestRate)
	expectedQuery := `SELECT id, duration_months, loan_interest_rate FROM loan_interest WHERE id =$1`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(r)
	ahc, err := l.repo.FindById(mock.MockLoanInterest.Id)
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), mock.MockLoanInterest.Id, ahc.Id)
}
func (l *LoanInterestRepoTestSuite) TestFindById_Failed() {
	expectedQuery := `SELECT id, duration_months, loan_interest_rate FROM loan_interest WHERE id =$1`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error"))
	ahc, err := l.repo.FindById(mock.MockLoanInterest.Id)
	assert.Error(l.T(), err)
	assert.Equal(l.T(), model.LoanInterest{}, ahc)
}
func (l *LoanInterestRepoTestSuite) TestDeleteById_Success() {
	expectedQuery := `DELETE FROM loan_interest WHERE id = $1`
	l.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnResult(sqlmock.NewResult(1, 1))
	err := l.repo.DeleteById(mock.MockLoanInterest.Id)
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LoanInterestRepoTestSuite) TestDeleteById_Failed() {
	expectedQuery := `DELETE FROM loan_interest WHERE id = $1`
	l.mockSQL.ExpectExec(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error delete"))
	err := l.repo.DeleteById(mock.MockLoanInterest.Id)
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LoanInterestRepoTestSuite) TestCreateNew_Success() {
	expectedQuery := `INSERT INTO loan_interest`
	l.mockSQL.ExpectExec(expectedQuery).WillReturnResult(sqlmock.NewResult(1, 1))
	err := l.repo.CreateNew(mock.MockLoanInterest)
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LoanInterestRepoTestSuite) TestCreateNew_Failed() {
	expectedQuery := `INSERT INTO loan_interest`
	l.mockSQL.ExpectExec(expectedQuery).WillReturnError(errors.New("error crated"))
	err := l.repo.CreateNew(mock.MockLoanInterest)
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LoanInterestRepoTestSuite) TestPagging_Success() {
	rows := sqlmock.NewRows([]string{"id", "duration_months", "loan_interest_rate"})
	for _, li := range mock.MockLoanInterestDatas {
		rows.AddRow(li.Id, li.DurationMonths, li.LoanInterestRate)
	}
	expectedQuery := `SELECT id, duration_months, loan_interest_rate FROM loan_interest LIMIT $2 OFFSET $1;`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM loan_interest`)).WillReturnRows(rowCount)

	uc, p, err := l.repo.Pagging(mock.MockPageReq)
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), 1, len(uc))
	assert.Equal(l.T(), 1, p.TotalRows)
}
func (l *LoanInterestRepoTestSuite) TestPagging_Fail() {
	// error select paging
	expectedQuery := `SELECT id, duration_months, loan_interest_rate FROM loan_interest LIMIT $2 OFFSET $1;`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnError(errors.New("failed"))
	uc, p, err := l.repo.Pagging(dto.PageRequest{})
	assert.Error(l.T(), err)
	assert.Nil(l.T(), uc)
	assert.Equal(l.T(), 0, p.TotalRows)
	// error select count
	rows := sqlmock.NewRows([]string{"id", "duration_months", "loan_interest_rate"})
	for _, li := range mock.MockLoanInterestDatas {
		rows.AddRow(li.Id, li.DurationMonths, li.LoanInterestRate)
	}
	expectedQuery = `SELECT id, duration_months, loan_interest_rate FROM loan_interest LIMIT $2 OFFSET $1;`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM loan_interest`)).WillReturnError(errors.New("failed"))
	uc, p, err = l.repo.Pagging(mock.MockPageReq)
	assert.Error(l.T(), err)
	assert.Nil(l.T(), uc)
	assert.Equal(l.T(), 0, p.TotalRows)
}
