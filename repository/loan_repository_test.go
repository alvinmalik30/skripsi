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

type LoanRepoTestSuite struct {
	suite.Suite
	mockDB   *sql.DB
	mockSQL  sqlmock.Sqlmock
	repo     LoanRepository
	userRepo BiodataUser
}

func (l *LoanRepoTestSuite) SetupTest() {
	d, s, err := sqlmock.New()
	assert.NoError(l.T(), err)
	l.mockDB = d
	l.mockSQL = s
	l.userRepo = NewBiodataUserRepository(l.mockDB)
	l.repo = NewLoanRepository(l.mockDB, l.userRepo)
}

func TestLoanRepoTestSuite(t *testing.T) {
	suite.Run(t, new(LoanRepoTestSuite))
}
func (l *LoanRepoTestSuite) TestCreate_Success() {
	l.mockSQL.ExpectBegin()
	l.mockSQL.ExpectExec(`INSERT INTO loan`).WillReturnResult(sqlmock.NewResult(1, 1))
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving - $1 WHERE user_credential_id = $2;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	l.mockSQL.ExpectExec(`INSERT INTO installenment_loan`).WillReturnResult(sqlmock.NewResult(1, 1))
	l.mockSQL.ExpectCommit()
	err := l.repo.Create(mock.MockLoan, mock.MockInstallLoanDatas)
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LoanRepoTestSuite) TestCreate_Failed() {
	// Begin failed
	l.mockSQL.ExpectBegin().WillReturnError(errors.New("begin failed"))
	err := l.repo.Create(mock.MockLoan, mock.MockInstallLoanDatas)
	assert.Error(l.T(), err)
	// Insert loan failed
	l.mockSQL.ExpectBegin()
	l.mockSQL.ExpectExec(`INSERT INTO loan`).WillReturnError(errors.New("insert failed"))
	err = l.repo.Create(mock.MockLoan, mock.MockInstallLoanDatas)
	assert.Error(l.T(), err)
	// Update saldo failed
	l.mockSQL.ExpectBegin()
	l.mockSQL.ExpectExec(`INSERT INTO loan`).WillReturnResult(sqlmock.NewResult(1, 1))
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving - $1 WHERE user_credential_id = $2;`)).WillReturnError(errors.New("update failed"))
	err = l.repo.Create(mock.MockLoan, mock.MockInstallLoanDatas)
	assert.Error(l.T(), err)
	// Insert installement loan
	l.mockSQL.ExpectBegin()
	l.mockSQL.ExpectExec(`INSERT INTO loan`).WillReturnResult(sqlmock.NewResult(1, 1))
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving - $1 WHERE user_credential_id = $2;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	l.mockSQL.ExpectExec(`INSERT INTO installenment_loan`).WillReturnError(errors.New("insert failed"))
	err = l.repo.Create(mock.MockLoan, mock.MockInstallLoanDatas)
	assert.Error(l.T(), err)
}
func (l *LoanRepoTestSuite) TestFindUploadedFile_Success() {
	rows := sqlmock.NewRows([]string{"id", "loan_id", "is_payed", "payment_installenment_cost", "payment_deadline", "total_amount_of_dept", "late_payment_fee_nominal", "late_payment_fee_unit", "late_payment_fee_day", "payment_date", "status", "transfer_confirmation_recipt", "recipt_file"})
	for _, hc := range mock.MockLoanInstallRespons {
		rows.AddRow(hc.Id, "1", hc.IsPayed, hc.PaymentInstallment, hc.PaymentDeadLine, hc.TotalAmountOfDepth, hc.LatePayment.LatePaymentFees, hc.LatePayment.LatePaymentDays, hc.LatePayment.LatePaymentFeesTotal, hc.PaymentDate, hc.Status, hc.TransferConfirmRecipe, hc.File)
	}
	expectedQuery := `SELECT id, loan_id, is_payed, payment_installenment_cost, payment_deadline, total_amount_of_dept, late_payment_fee_nominal, late_payment_fee_unit, late_payment_fee_day, payment_date, status, transfer_confirmation_recipt, recipt_file FROM installenment_loan WHERE is_payed = false AND transfer_confirmation_recipt = true;`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(rows)
	lir, err := l.repo.FindUploadedFile()
	assert.Nil(l.T(), err)
	assert.Equal(l.T(), 1, len(lir))
}
func (l *LoanRepoTestSuite) TestFindUploadedFile_Failed() {
	expectedQuery := `SELECT id, loan_id, is_payed, payment_installenment_cost, payment_deadline, total_amount_of_dept, late_payment_fee_nominal, late_payment_fee_unit, late_payment_fee_day, payment_date, status, transfer_confirmation_recipt, recipt_file FROM installenment_loan WHERE is_payed = false AND transfer_confirmation_recipt = true;`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error"))
	lir, err := l.repo.FindUploadedFile()
	assert.Error(l.T(), err)
	assert.Equal(l.T(), 0, len(lir))
}
func (l *LoanRepoTestSuite) TestUpload_Success() {
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE installenment_loan SET payment_date = $2, status = 'witing for accepment', transfer_confirmation_recipt = true, recipt_file = $3 WHERE id = $1;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	err := l.repo.Upload(mock.MockLoanInstallRespons[0])
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LoanRepoTestSuite) TestUpload_Failed() {
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE installenment_loan SET payment_date = $2, status = 'witing for accepment', transfer_confirmation_recipt = true, recipt_file = $3 WHERE id = $1;`)).WillReturnError(errors.New("error"))
	err := l.repo.Upload(mock.MockLoanInstallRespons[0])
	assert.Error(l.T(), err)
	assert.NotNil(l.T(), err)
}
func (l *LoanRepoTestSuite) TestAccepted_Success() {
	l.mockSQL.ExpectBegin()
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE installenment_loan SET is_payed = $2, status = $3 WHERE id = $1;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving + $1 WHERE user_credential_id = '456';`)).WillReturnResult(sqlmock.NewResult(1, 1))
	l.mockSQL.ExpectCommit()
	err := l.repo.Accepted(mock.MockInstallLoanByIdResp)
	assert.Nil(l.T(), err)
	assert.NoError(l.T(), err)
}
func (l *LoanRepoTestSuite) TestAccepted_Failed() {
	// Begin failed
	l.mockSQL.ExpectBegin().WillReturnError(errors.New("begin failed"))
	err := l.repo.Accepted(mock.MockInstallLoanByIdResp)
	assert.Error(l.T(), err)
	// Update installenment_loan failed
	l.mockSQL.ExpectBegin()
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE installenment_loan SET is_payed = $2, status = $3 WHERE id = $1;`)).WillReturnError(errors.New("error"))
	err = l.repo.Accepted(mock.MockInstallLoanByIdResp)
	assert.Error(l.T(), err)
	// Update saldo failed
	l.mockSQL.ExpectBegin()
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE installenment_loan SET is_payed = $2, status = $3 WHERE id = $1;`)).WillReturnResult(sqlmock.NewResult(1, 1))
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving + $1 WHERE user_credential_id = '456';`)).WillReturnError(errors.New("update failed"))
	err = l.repo.Accepted(mock.MockInstallLoanByIdResp)
	assert.Error(l.T(), err)
}
func (l *LoanRepoTestSuite) TestUpdateLateFee_Success() {
	rows := sqlmock.NewRows([]string{"id", "total_amount_of_dept", "payment_deadline", "late_payment_fee_nominal", "late_payment_fee_unit"})
	for _, hc := range mock.MockInstallLoanDatas {
		rows.AddRow(hc.Id, hc.TotalAmountOfDepth, hc.PaymentDeadLine, hc.LatePaymentFeesNominal, hc.LatePaymentFeesUnit)
	}
	expectedQuery := `SELECT id, total_amount_of_dept, payment_deadline, late_payment_fee_nominal, late_payment_fee_unit FROM installenment_loan WHERE is_payed = false AND transfer_confirmation_recipt = true AND payment_deadline < now();`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(rows)
	l.mockSQL.ExpectBegin()
	l.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE installenment_loan SET late_payment_fee_day = $2, late_payment_fee_total = $3, total_amount_of_dept = total_amount_of_dept + $4 WHERE id = $1`)).WillReturnResult(sqlmock.NewResult(1, 1))
	l.mockSQL.ExpectCommit()
	err := l.repo.UpdateLateFee()
	assert.Nil(l.T(), err)
}
func (l *LoanRepoTestSuite) TestUpdateLateFee_Failed() {
	// error select installenment_loan
	expectedQuery := `SELECT id, total_amount_of_dept, payment_deadline, late_payment_fee_nominal, late_payment_fee_unit FROM installenment_loan WHERE is_payed = false AND transfer_confirmation_recipt = true AND payment_deadline < now();`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error"))
	err := l.repo.UpdateLateFee()
	assert.Error(l.T(), err)
	// error begin
	rows := sqlmock.NewRows([]string{"id", "total_amount_of_dept", "payment_deadline", "late_payment_fee_nominal", "late_payment_fee_unit"})
	for _, hc := range mock.MockInstallLoanDatas {
		rows.AddRow(hc.Id, hc.TotalAmountOfDepth, hc.PaymentDeadLine, hc.LatePaymentFeesNominal, hc.LatePaymentFeesUnit)
	}
	expectedQuery = `SELECT id, total_amount_of_dept, payment_deadline, late_payment_fee_nominal, late_payment_fee_unit FROM installenment_loan WHERE is_payed = false AND transfer_confirmation_recipt = true AND payment_deadline < now();`
	l.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(rows)
	l.mockSQL.ExpectBegin().WillReturnError(errors.New("error"))
	err = l.repo.UpdateLateFee()
	assert.Error(l.T(), err)
}
