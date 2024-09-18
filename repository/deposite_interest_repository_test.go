package repository

import (
	"database/sql"
	"errors"
	"polen/mock"
	"polen/model/dto"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DepositeInterestRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    DepositeInterest
}

func (d *DepositeInterestRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(d.T(), err)
	d.mockDB = db
	d.mockSQL = mock
	d.repo = NewDepositeInterestRepository(d.mockDB)
}
func TestDepositeInterestRepoTestSuite(t *testing.T) {
	suite.Run(t, new(DepositeInterestRepoTestSuite))
}
func (d *DepositeInterestRepoTestSuite) TestDeleteById_Success() {
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`DELETE FROM deposit_interest WHERE id = $1`)).WithArgs(mock.MockDepositeInterestReq.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	err := d.repo.DeleteById(mock.MockDepositeInterestReq.Id)
	assert.Nil(d.T(), err)
}
func (d *DepositeInterestRepoTestSuite) TestDeleteById_Failed() {
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`DELETE FROM deposit_interest WHERE id = $1`)).WithArgs(mock.MockDepositeInterestReq.Id).WillReturnError(errors.New("failed delete"))
	err := d.repo.DeleteById(mock.MockDepositeInterestReq.Id)
	assert.Error(d.T(), err)
}
func (d *DepositeInterestRepoTestSuite) TestFindById_Success() {
	rows := sqlmock.NewRows([]string{"id", "interest_rate", "tax_rate", "duration_mounth"})
	rows.AddRow(mock.MockDepositeInterest.Id, mock.MockDepositeInterest.InterestRate, mock.MockDepositeInterest.TaxRate, mock.MockDepositeInterest.DurationMounth)
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT id, interest_rate, tax_rate, duration_mounth FROM deposit_interest WHERE id =$1`)).WithArgs(mock.MockDepositeInterest.Id).WillReturnRows(rows)
	di, err := d.repo.FindById(mock.MockDepositeInterest.Id)
	assert.Nil(d.T(), err)
	assert.Equal(d.T(), mock.MockDepositeInterest.Id, di.Id)
}
func (d *DepositeInterestRepoTestSuite) TestFindById_Failed() {
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT id, interest_rate, tax_rate, duration_mounth FROM deposit_interest WHERE id =$1`)).WithArgs(mock.MockDepositeInterest.Id).WillReturnError(errors.New("failed find "))
	di, err := d.repo.FindById(mock.MockDepositeInterestReq.Id)
	assert.Error(d.T(), err)
	assert.Equal(d.T(), dto.DepositeInterestRequest{}, di)
}
func (d *DepositeInterestRepoTestSuite) TestUpdate_Success() {
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE deposit_interest SET interest_rate = $2, tax_rate = $3, duration_mounth = $4 WHERE id = $1;`)).WithArgs(mock.MockDepositeInterest.Id, mock.MockDepositeInterest.InterestRate, mock.MockDepositeInterest.TaxRate, mock.MockDepositeInterest.DurationMounth).WillReturnResult(sqlmock.NewResult(1, 1))
	err := d.repo.Update(mock.MockDepositeInterestReq)
	assert.Nil(d.T(), err)
}
func (d *DepositeInterestRepoTestSuite) TestUpdate_Failed() {
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE deposit_interest SET interest_rate = $2, tax_rate = $3, duration_mounth = $4 WHERE id = $1;`)).WithArgs(mock.MockDepositeInterest.Id, mock.MockDepositeInterest.InterestRate, mock.MockDepositeInterest.TaxRate, mock.MockDepositeInterest.DurationMounth).WillReturnError(errors.New("failed update"))
	err := d.repo.Update(mock.MockDepositeInterestReq)
	assert.Error(d.T(), err)
}
func (d *DepositeInterestRepoTestSuite) TestSave_Success() {
	d.mockSQL.ExpectExec(`INSERT INTO deposit_interest`).WithArgs(mock.MockDepositeInterest.Id, mock.MockDepositeInterest.CreateDate, mock.MockDepositeInterest.InterestRate, mock.MockDepositeInterest.TaxRate, mock.MockDepositeInterest.DurationMounth).WillReturnResult(sqlmock.NewResult(1, 1))
	err := d.repo.Save(mock.MockDepositeInterest)
	assert.Nil(d.T(), err)
	assert.NoError(d.T(), err)
}
func (d *DepositeInterestRepoTestSuite) TestSave_Failed() {
	d.mockSQL.ExpectExec(`INSERT INTO deposit_interest`).WithArgs(mock.MockDepositeInterest.Id, mock.MockDepositeInterest.CreateDate, mock.MockDepositeInterest.InterestRate, mock.MockDepositeInterest.TaxRate, mock.MockDepositeInterest.DurationMounth).WillReturnError(errors.New("failed save"))
	err := d.repo.Save(mock.MockDepositeInterest)
	assert.Error(d.T(), err)
	assert.NotNil(d.T(), err)
}
func (d *DepositeInterestRepoTestSuite) TestPagging_Success() {
	rows := sqlmock.NewRows([]string{"id", "interest_rate", "tax_rate", "duration_mounth"})
	for _, row := range mock.MockDeposites {
		rows.AddRow(row.Id, row.InterestRate, row.TaxRate, row.DurationMounth)
	}
	expectedSQL := `SELECT id, interest_rate, tax_rate, duration_mounth FROM deposit_interest LIMIT $2 OFFSET $1`
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM deposit_interest`)).WillReturnRows(rowCount)

	uc, p, err := d.repo.Pagging(mock.MockPageReq)
	assert.Nil(d.T(), err)
	assert.Equal(d.T(), 1, len(uc))
	assert.Equal(d.T(), 1, p.TotalRows)
}
func (d *DepositeInterestRepoTestSuite) TestPagging_Fail() {
	// error select paging
	expectedSQL := `SELECT id, interest_rate, tax_rate, duration_mounth FROM deposit_interest LIMIT $2 OFFSET $1`
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("failed"))
	uc, p, err := d.repo.Pagging(dto.PageRequest{})
	assert.Error(d.T(), err)
	assert.Nil(d.T(), uc)
	assert.Equal(d.T(), 0, p.TotalRows)
	// error select count
	rows := sqlmock.NewRows([]string{"id", "interest_rate", "tax_rate", "duration_mounth"})
	for _, row := range mock.MockDeposites {
		rows.AddRow(row.Id, row.InterestRate, row.TaxRate, row.DurationMounth)
	}
	expectedSQL = `SELECT id, interest_rate, tax_rate, duration_mounth FROM deposit_interest LIMIT $2 OFFSET $1`
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM deposit_interest`)).WillReturnError(errors.New("failed"))
	uc, p, err = d.repo.Pagging(mock.MockPageReq)
	assert.Error(d.T(), err)
	assert.Nil(d.T(), uc)
	assert.Equal(d.T(), 0, p.TotalRows)
}
