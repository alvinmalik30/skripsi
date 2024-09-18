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

type TopUpRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    TopUp
}

func (t *TopUpRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(t.T(), err)
	t.mockDB = db
	t.mockSQL = mock
	t.repo = NewTopUpRepository(t.mockDB)
}
func TestTopUpRepoTestSuite(t *testing.T) {
	suite.Run(t, new(TopUpRepoTestSuite))
}
func (t *TopUpRepoTestSuite) TestPagging_Success() {
	rows := sqlmock.NewRows([]string{"id", "maturity_time", "top_up_amount", "accepted_time", "accepted_status", "status_information", "transfer_confirmation_recipt", "recipt_file"})
	for _, row := range mock.MockListTopUp {
		rows.AddRow(row.Id, row.MaturityTime, row.TopUpAmount, row.AcceptedTime, row.Accepted, row.Status, row.TransferConfirmRecipe, row.File)
	}
	expectedSQL := `SELECT id, maturity_time, top_up_amount, accepted_time, accepted_status, status_information, transfer_confirmation_recipt, recipt_file FROM top_up LIMIT $2 OFFSET $1;`
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM top_up`)).WillReturnRows(rowCount)

	uc, p, err := t.repo.Pagging(mock.MockPageReq)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), 1, len(uc))
	assert.Equal(t.T(), 1, p.TotalRows)
}
func (t *TopUpRepoTestSuite) TestPagging_Fail() {
	// error select paging
	expectedSQL := `SELECT id, maturity_time, top_up_amount, accepted_time, accepted_status, status_information, transfer_confirmation_recipt, recipt_file FROM top_up LIMIT $2 OFFSET $1;`
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("failed"))
	uc, p, err := t.repo.Pagging(dto.PageRequest{})
	assert.Error(t.T(), err)
	assert.Nil(t.T(), uc)
	assert.Equal(t.T(), 0, p.TotalRows)
	// error select count
	rows := sqlmock.NewRows([]string{"id", "maturity_time", "top_up_amount", "accepted_time", "accepted_status", "status_information", "transfer_confirmation_recipt", "recipt_file"})
	for _, row := range mock.MockListTopUp {
		rows.AddRow(row.Id, row.MaturityTime, row.TopUpAmount, row.AcceptedTime, row.Accepted, row.Status, row.TransferConfirmRecipe, row.File)
	}
	expectedSQL = `SELECT id, maturity_time, top_up_amount, accepted_time, accepted_status, status_information, transfer_confirmation_recipt, recipt_file FROM top_up LIMIT $2 OFFSET $1;`
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM top_up`)).WillReturnError(errors.New("failed"))
	uc, p, err = t.repo.Pagging(mock.MockPageReq)
	assert.Error(t.T(), err)
	assert.Nil(t.T(), uc)
	assert.Equal(t.T(), 0, p.TotalRows)
}
func (t *TopUpRepoTestSuite) TestFindUploadedFile_Success() {
	rows := sqlmock.NewRows([]string{"id", "maturity_time", "top_up_amount", "accepted_time", "accepted_status", "status_information", "transfer_confirmation_recipt", "recipt_file"})
	for _, row := range mock.MockListTopUp {
		rows.AddRow(row.Id, row.MaturityTime, row.TopUpAmount, row.AcceptedTime, row.Accepted, row.Status, row.TransferConfirmRecipe, row.File)
	}
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT id, maturity_time, top_up_amount, accepted_time, accepted_status, status_information, transfer_confirmation_recipt, recipt_file FROM top_up WHERE transfer_confirmation_recipt = true AND accepted_status = false;`)).WillReturnRows(rows)
	tu, err := t.repo.FindUploadedFile()
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), 1, len(tu))
}
func (t *TopUpRepoTestSuite) TestFindUploadedFile_Failed() {
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT id, maturity_time, top_up_amount, accepted_time, accepted_status, status_information, transfer_confirmation_recipt, recipt_file FROM top_up WHERE transfer_confirmation_recipt = true AND accepted_status = false;`)).WillReturnError(errors.New("error"))
	tu, err := t.repo.FindUploadedFile()
	assert.Error(t.T(), err)
	assert.Nil(t.T(), tu)
}
func (t *TopUpRepoTestSuite) TestConfimUpload_Success() {
	t.mockSQL.ExpectBegin()
	t.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE top_up SET accepted_status = $2, accepted_time = $3, status_information = $4, transfer_confirmation_recipt = $5, recipt_file = $6, maturity_time = $7 WHERE id = $1;`)).WithArgs(mock.MockTopUp.Id, mock.MockTopUp.Accepted, mock.MockTopUp.AcceptedTime, mock.MockTopUp.Status, mock.MockTopUp.TransferConfirmRecipe, mock.MockTopUp.File, mock.MockTopUp.MaturityTime).WillReturnResult(sqlmock.NewResult(1, 1))
	t.mockSQL.ExpectExec(regexp.QuoteMeta(`UPDATE saldo SET total_saving = total_saving + $1 WHERE user_credential_id = $2`)).WithArgs(mock.MockTopUp.TopUpAmount, mock.MockTopUp.UserCredential.Id).WillReturnResult(sqlmock.NewResult(1, 1))
	t.mockSQL.ExpectCommit()
	err := t.repo.ConfimUpload(mock.MockTopUp)
	assert.Nil(t.T(), err)
	assert.NoError(t.T(), err)
}
func (t *TopUpRepoTestSuite) TestConfimUpload_Failed() {
	// Begin failed
	t.mockSQL.ExpectBegin().WillReturnError(errors.New("begin failed"))
	err := t.repo.ConfimUpload(mock.MockTopUp)
	assert.Error(t.T(), err)
	// Update topup failed
	t.mockSQL.ExpectBegin()
	t.mockSQL.ExpectExec(`UPDATE top_up SET accepted_status = $2, accepted_time = $3, status_information = $4, transfer_confirmation_recipt = $5, recipt_file = $6, maturity_time = $7 WHERE id = $1;`).WillReturnError(errors.New("update failed"))
	err = t.repo.ConfimUpload(mock.MockTopUp)
	assert.Error(t.T(), err)
	// Update saldo failed
	t.mockSQL.ExpectBegin()
	t.mockSQL.ExpectExec(`UPDATE top_up SET accepted_status = $2, accepted_time = $3, status_information = $4, transfer_confirmation_recipt = $5, recipt_file = $6, maturity_time = $7 WHERE id = $1;`).WillReturnResult(sqlmock.NewResult(1, 1))
	t.mockSQL.ExpectExec(`UPDATE saldo SET total_saving = total_saving + $1 WHERE user_credential_id = $2`).WillReturnError(errors.New("insert failed"))
	err = t.repo.ConfimUpload(mock.MockTopUp)
	assert.Error(t.T(), err)
}
func (t *TopUpRepoTestSuite) TestNotConfimUpload_Success() {
	expectedSQL := `UPDATE top_up SET accepted_status = $2, status_information = $3, transfer_confirmation_recipt = $4, recipt_file = $5, maturity_time = $6 WHERE id = $1;`
	t.mockSQL.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockTopUp.Id, mock.MockTopUp.Accepted, mock.MockTopUp.Status, mock.MockTopUp.TransferConfirmRecipe, mock.MockTopUp.File, mock.MockTopUp.MaturityTime).WillReturnResult(sqlmock.NewResult(1, 1))
	err := t.repo.NotConfimUpload(mock.MockTopUp)
	assert.Nil(t.T(), err)
	assert.NoError(t.T(), err)
}
func (t *TopUpRepoTestSuite) TestNotConfimUpload_Failed() {
	expectedSQL := `UPDATE top_up SET accepted_status = $2, status_information = $3, transfer_confirmation_recipt = $4, recipt_file = $5, maturity_time = $6 WHERE id = $1;`
	t.mockSQL.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockTopUp.Id, mock.MockTopUp.Accepted, mock.MockTopUp.Status, mock.MockTopUp.TransferConfirmRecipe, mock.MockTopUp.File, mock.MockTopUp.MaturityTime).WillReturnError(errors.New("error"))
	err := t.repo.NotConfimUpload(mock.MockTopUp)
	assert.Error(t.T(), err)
	assert.NotNil(t.T(), err)
}
func (t *TopUpRepoTestSuite) TestUpload_Success() {
	expectedSQL := `UPDATE top_up SET status_information = $2, transfer_confirmation_recipt = $3, recipt_file = $4 WHERE id = $1;`
	t.mockSQL.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockTopUp.Id, mock.MockTopUp.Status, mock.MockTopUp.TransferConfirmRecipe, mock.MockTopUp.File).WillReturnResult(sqlmock.NewResult(1, 1))
	err := t.repo.Upload(mock.MockTopUp)
	assert.Nil(t.T(), err)
	assert.NoError(t.T(), err)
}
func (t *TopUpRepoTestSuite) TestUpload_Failed() {
	expectedSQL := `UPDATE top_up SET status_information = $2, transfer_confirmation_recipt = $3, recipt_file = $4 WHERE id = $1;`
	t.mockSQL.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockTopUp.Id, mock.MockTopUp.Status, mock.MockTopUp.TransferConfirmRecipe, mock.MockTopUp.File).WillReturnError(errors.New("error"))
	err := t.repo.Upload(mock.MockTopUp)
	assert.Error(t.T(), err)
	assert.NotNil(t.T(), err)
}
func (t *TopUpRepoTestSuite) TestFindById_Success() {
	rows := sqlmock.NewRows([]string{"user_credential_id", "user_credential_username", "user_credential_email", "user_credential_role", "user_credential_is_active", "user_credential_virtual_account_number", "biodata_full_name", "biodata_nik", "biodata_phone_number", "biodata_occupation", "biodata_place_of_birth", "biodata_date_of_birth", "biodata_postal_code", "top_up_id", "top_up_amount", "top_up_maturity_time", "top_up_accepted_time", "top_up_accepted_status", "top_up_status_information", "top_up_transfer_confirmation_recipt", "top_up_recipt_file"})
	rows.AddRow(mock.MockTopUpId.UserCredential.Id, mock.MockTopUpId.UserCredential.Username, mock.MockTopUpId.UserCredential.Email, mock.MockTopUpId.UserCredential.Role, mock.MockTopUpId.UserCredential.IsActive, mock.MockTopUpId.UserCredential.VaNumber, mock.MockTopUpId.UserBio.NamaLengkap, mock.MockTopUpId.UserBio.Nik, mock.MockTopUpId.UserBio.NomorTelepon, mock.MockTopUpId.UserBio.Pekerjaan, mock.MockTopUpId.UserBio.TempatLahir, mock.MockTopUpId.UserBio.TanggalLahir, mock.MockTopUpId.UserBio.KodePos, mock.MockTopUpId.TopUp.Id, mock.MockTopUpId.TopUp.TopUpAmount, mock.MockTopUpId.TopUp.MaturityTime, mock.MockTopUpId.TopUp.AcceptedTime, mock.MockTopUpId.TopUp.Accepted, mock.MockTopUpId.TopUp.Status, mock.MockTopUpId.TopUp.TransferConfirmRecipe, mock.MockTopUpId.TopUp.File)
	expectedSQL := `SELECT u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code, t.id, t.top_up_amount, t.maturity_time, t.accepted_time, t.accepted_status, t.status_information, t.transfer_confirmation_recipt, t.recipt_file FROM top_up t JOIN biodata b ON b.user_credential_id = t.user_credential_id JOIN user_credential u ON u.id = b.user_credential_id WHERE t.id = $1`
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockTopUp.Id).WillReturnRows(rows)
	uc, err := t.repo.FindById(mock.MockTopUpId.TopUp.Id)
	assert.Nil(t.T(), err)
	assert.NoError(t.T(), err)
	assert.Equal(t.T(), mock.MockTopUpId.TopUp.Id, uc.TopUp.Id)
}
func (t *TopUpRepoTestSuite) TestFindById_Failed() {
	expectedSQL := `SELECT u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code, t.id, t.top_up_amount, t.maturity_time, t.accepted_time, t.accepted_status, t.status_information, t.transfer_confirmation_recipt, t.recipt_file FROM top_up t JOIN biodata b ON b.user_credential_id = t.user_credential_id JOIN user_credential u ON u.id = b.user_credential_id WHERE t.id = $1`
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockTopUp.Id).WillReturnError(errors.New("error"))
	uc, err := t.repo.FindById(mock.MockTopUpId.TopUp.Id)
	assert.Error(t.T(), err)
	assert.NotNil(t.T(), err)
	assert.Equal(t.T(), dto.TopUpById{}, uc)
}
func (t *TopUpRepoTestSuite) TestSave_Success() {
	t.mockSQL.ExpectExec(`INSERT INTO top_up`).WillReturnResult(sqlmock.NewResult(1, 1))
	err := t.repo.Save(mock.MockTopUp)
	assert.Nil(t.T(), err)
	assert.NoError(t.T(), err)
}
func (t *TopUpRepoTestSuite) TestSave_Failed() {
	t.mockSQL.ExpectExec(`INSERT INTO top_up`).WillReturnError(errors.New("error"))
	err := t.repo.Save(mock.MockTopUp)
	assert.Error(t.T(), err)
	assert.NotNil(t.T(), err)
}
func (t *TopUpRepoTestSuite) TestFindByIdUser_Success() {
	row := sqlmock.NewRows([]string{"user_credential_id", "user_credential_username", "user_credential_email", "user_credential_role", "user_credential_is_active", "user_credential_virtual_account_number", "biodata_full_name", "biodata_nik", "biodata_phone_number", "biodata_occupation", "biodata_place_of_birth", "biodata_date_of_birth", "biodata_postal_code"})
	row.AddRow(mock.MockUserCred.Id, mock.MockUserCred.Username, mock.MockUserCred.Email, mock.MockUserCred.Role, mock.MockUserCred.IsActive, mock.MockUserCred.VANumber, mock.MockBiodata.NamaLengkap, mock.MockBiodata.Nik, mock.MockBiodata.NomorTelepon, mock.MockBiodata.Pekerjaan, mock.MockBiodata.TempatLahir, mock.MockBiodata.TanggalLahir, mock.MockBiodata.KodePos)
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code FROM biodata b JOIN user_credential u ON u.id = b.user_credential_id WHERE b.user_credential_id = $1;`)).WithArgs(mock.MockBiodata.UserCredential.Id).WillReturnRows(row)

	rows := sqlmock.NewRows([]string{"id", "maturity_time", "top_up_amount", "accepted_time", "accepted_status", "status_information", "transfer_confirmation_recipt", "recipt_file"})
	for _, row := range mock.MockListTopUp {
		rows.AddRow(row.Id, row.MaturityTime, row.TopUpAmount, row.AcceptedTime, row.Accepted, row.Status, row.TransferConfirmRecipe, row.File)
	}
	expectedSQL := `SELECT id, maturity_time, top_up_amount, accepted_time, accepted_status, status_information, transfer_confirmation_recipt, recipt_file FROM top_up WHERE user_credential_id = $1`
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockTopUp.UserCredential.Id).WillReturnRows(rows)

	tubu, err := t.repo.FindByIdUser(mock.MockTopUpByUser.UserCredential.Id)
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), mock.MockTopUpByUser.UserCredential.Id, tubu.UserCredential.Id)
}
func (t *TopUpRepoTestSuite) TestFindByIdUser_Failed() {
	// error select biodata by user id
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code FROM biodata b JOIN user_credential u ON u.id = b.user_credential_id WHERE b.user_credential_id = $1;`)).WillReturnError(errors.New("error"))
	tubu, err := t.repo.FindByIdUser(mock.MockTopUpByUser.UserCredential.Id)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), dto.TopUpByUser{}, tubu)

	// error select topup by user id
	row := sqlmock.NewRows([]string{"user_credential_id", "user_credential_username", "user_credential_email", "user_credential_role", "user_credential_is_active", "user_credential_virtual_account_number", "biodata_full_name", "biodata_nik", "biodata_phone_number", "biodata_occupation", "biodata_place_of_birth", "biodata_date_of_birth", "biodata_postal_code"})
	row.AddRow(mock.MockUserCred.Id, mock.MockUserCred.Username, mock.MockUserCred.Email, mock.MockUserCred.Role, mock.MockUserCred.IsActive, mock.MockUserCred.VANumber, mock.MockBiodata.NamaLengkap, mock.MockBiodata.Nik, mock.MockBiodata.NomorTelepon, mock.MockBiodata.Pekerjaan, mock.MockBiodata.TempatLahir, mock.MockBiodata.TanggalLahir, mock.MockBiodata.KodePos)
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code FROM biodata b JOIN user_credential u ON u.id = b.user_credential_id WHERE b.user_credential_id = $1;`)).WithArgs(mock.MockBiodata.UserCredential.Id).WillReturnRows(row)

	rows := sqlmock.NewRows([]string{"id", "maturity_time", "top_up_amount", "accepted_time", "accepted_status", "status_information", "transfer_confirmation_recipt", "recipt_file"})
	for _, row := range mock.MockListTopUp {
		rows.AddRow(row.Id, row.MaturityTime, row.TopUpAmount, row.AcceptedTime, row.Accepted, row.Status, row.TransferConfirmRecipe, row.File)
	}
	expectedSQL := `SELECT id, maturity_time, top_up_amount, accepted_time, accepted_status, status_information, transfer_confirmation_recipt, recipt_file FROM top_up WHERE user_credential_id = $1`
	t.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockTopUp.UserCredential.Id).WillReturnError(errors.New("error"))
	tubu, err = t.repo.FindByIdUser(mock.MockTopUpByUser.UserCredential.Id)
	assert.Error(t.T(), err)
	assert.Equal(t.T(), dto.TopUpByUser{}, tubu)
}
