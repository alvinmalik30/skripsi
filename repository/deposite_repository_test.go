package repository

import (
	"database/sql"
	"errors"
	"polen/mock"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type DepositeRepoTestSuite struct {
	suite.Suite
	mockSQL sqlmock.Sqlmock
	mockDB  *sql.DB
	user    BiodataUser
	repo    DepositeRepository
}

func (d *DepositeRepoTestSuite) SetupTest() {
	db, sql, err := sqlmock.New()
	assert.NoError(d.T(), err)
	d.mockDB = db
	d.mockSQL = sql
	d.repo = NewDepositeRepository(d.mockDB, d.user)
}

func TestDepositeRepoTestSuite(t *testing.T) {
	suite.Run(t, new(DepositeRepoTestSuite))
}
func (d *DepositeRepoTestSuite) TestPagging_Success() {
	rows := sqlmock.NewRows([]string{"id", "deposit_amount", "interest_rate", "tax_rate", "duration", "created_date", "maturity_date", "status", "gross_profit", "tax", "net_profit", "total_return"})
	for _, deposite := range mock.MockDepositeDto {
		rows.AddRow(deposite.Id, deposite.DepositeAmount, deposite.InterestRate, deposite.TaxRate, deposite.DurationMounth, deposite.CreateDate, deposite.MaturityDate, deposite.Status, deposite.GrossProfit, deposite.Tax, deposite.NetProfit, deposite.TotalReturn)
	}
	expectedSQL := `SELECT id, deposit_amount, interest_rate, tax_rate, duration, created_date, maturity_date, status, gross_profit, tax, net_profit, total_return FROM deposit LIMIT $2 OFFSET $1;`
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM deposit_interest`)).WillReturnRows(rowCount)

	uc, p, err := d.repo.Pagging(mock.MockPageReq)
	assert.Nil(d.T(), err)
	assert.Equal(d.T(), 1, len(uc))
	assert.Equal(d.T(), 1, p.TotalRows)
}
func (d *DepositeRepoTestSuite) TestPagging_Failed() {
	// error select paging
	expectedSQL := `SELECT id, deposit_amount, interest_rate, tax_rate, duration, created_date, maturity_date, status, gross_profit, tax, net_profit, total_return FROM deposit LIMIT $2 OFFSET $1;`
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnError(errors.New("error"))
	uc, p, err := d.repo.Pagging(mock.MockPageReq)
	assert.Error(d.T(), err)
	assert.Nil(d.T(), uc)
	assert.Equal(d.T(), 0, p.TotalRows)
	// error select count
	rows := sqlmock.NewRows([]string{"id", "deposit_amount", "interest_rate", "tax_rate", "duration", "created_date", "maturity_date", "status", "gross_profit", "tax", "net_profit", "total_return"})
	for _, deposite := range mock.MockDepositeDto {
		rows.AddRow(deposite.Id, deposite.DepositeAmount, deposite.InterestRate, deposite.TaxRate, deposite.DurationMounth, deposite.CreateDate, deposite.MaturityDate, deposite.Status, deposite.GrossProfit, deposite.Tax, deposite.NetProfit, deposite.TotalReturn)
	}
	expectedSQL = `SELECT id, deposit_amount, interest_rate, tax_rate, duration, created_date, maturity_date, status, gross_profit, tax, net_profit, total_return FROM deposit LIMIT $2 OFFSET $1;`
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)
	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM deposit_interest`)).WillReturnError(errors.New("error"))
	uc, p, err = d.repo.Pagging(mock.MockPageReq)
	assert.Error(d.T(), err)
	assert.Nil(d.T(), uc)
	assert.Equal(d.T(), 0, p.TotalRows)
}
func (d *DepositeRepoTestSuite) TestCreateDeposite_Success() {
	d.mockSQL.ExpectBegin()
	d.mockSQL.ExpectExec(`INSERT INTO deposit`).WillReturnResult(sqlmock.NewResult(1, 1))
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving - $1 WHERE user_credential_id = $2;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving + $1 WHERE user_credential_id = 456;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	d.mockSQL.ExpectCommit()
	err := d.repo.CreateDeposite(mock.MockDeposite)
	assert.Nil(d.T(), err)
	assert.NoError(d.T(), err)
}
func (d *DepositeRepoTestSuite) TestCreateDeposite_Failed() {
	// Begin failed
	d.mockSQL.ExpectBegin().WillReturnError(errors.New("begin failed"))
	err := d.repo.CreateDeposite(mock.MockDeposite)
	assert.Error(d.T(), err)
	// Insert deposit failed
	d.mockSQL.ExpectBegin()
	d.mockSQL.ExpectExec(`INSERT INTO deposit`).WillReturnError(errors.New("insert failed"))
	err = d.repo.CreateDeposite(mock.MockDeposite)
	assert.Error(d.T(), err)
	// Update saldo failed
	d.mockSQL.ExpectBegin()
	d.mockSQL.ExpectExec(`INSERT INTO deposit`).WillReturnResult(sqlmock.NewResult(1, 1))
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving - $1 WHERE user_credential_id = $2;`)).WillReturnError(errors.New("update failed"))
	err = d.repo.CreateDeposite(mock.MockDeposite)
	assert.Error(d.T(), err)
	// Update saldo failed
	d.mockSQL.ExpectBegin()
	d.mockSQL.ExpectExec(`INSERT INTO deposit`).WillReturnResult(sqlmock.NewResult(1, 1))
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving - $1 WHERE user_credential_id = $2;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving + $1 WHERE user_credential_id = 456;`)).WillReturnError(errors.New("updated failed"))
	err = d.repo.CreateDeposite(mock.MockDeposite)
	assert.Error(d.T(), err)
}
func (d *DepositeRepoTestSuite) TestUpdate_Success() {
	rows := sqlmock.NewRows([]string{"id", "user_credential_id", "tax", "status", "total_return"})
	for _, hc := range mock.MockDepositesDTO {
		rows.AddRow(hc.Id, hc.UserCredential.Id, hc.Tax, hc.Status, hc.TotalReturn)
	}
	expectedQuery := `SELECT id, user_credential_id, tax, status, total_return FROM deposit WHERE status = true AND maturity_date < now();`
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(rows)
	d.mockSQL.ExpectBegin()
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE deposit SET status = false WHERE id = $1;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving + $1 WHERE user_credential_id = $2;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving - $1 - $2 WHERE user_credential_id = $3;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	d.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving + $1 WHERE user_credential_id = $2;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	d.mockSQL.ExpectCommit()
	err := d.repo.Update()
	assert.Nil(d.T(), err)
}

func (d *DepositeRepoTestSuite) TestUpdate_Failed() {
	// error select installenment_loan
	expectedQuery := `SELECT id, user_credential_id, tax, status, total_return FROM deposit WHERE status = true AND maturity_date < now();`
	d.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error"))
	err := d.repo.Update()
	assert.Error(d.T(), err)
}
