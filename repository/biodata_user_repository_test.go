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

type BiodataUserRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    BiodataUser
}

func (b *BiodataUserRepoTestSuite) SetupTest() {
	db, mock, err := sqlmock.New()
	assert.NoError(b.T(), err)
	b.mockDB = db
	b.mockSQL = mock
	b.repo = NewBiodataUserRepository(b.mockDB)
}
func TestBiodataUserRepoTestSuite(t *testing.T) {
	suite.Run(t, new(BiodataUserRepoTestSuite))
}
func (b *BiodataUserRepoTestSuite) TestAdminUpdate_Success() {
	expectedSQL := `UPDATE biodata SET is_eglible = $2, status_update = $3, additional_information = $4 WHERE id = $1`
	b.mockSQL.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockBiodata.Id, mock.MockBiodata.IsAglible, mock.MockBiodata.StatusUpdate, mock.MockBiodata.Information).WillReturnResult(sqlmock.NewResult(1, 1))
	err := b.repo.AdminUpdate(mock.MockBiodata)
	assert.Nil(b.T(), err)
}
func (b *BiodataUserRepoTestSuite) TestAdminUpdate_Failed() {
	expectedSQL := `UPDATE biodata SET is_eglible = $2, status_update = $3, additional_information = $4 WHERE id = $1`
	b.mockSQL.ExpectExec(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("failed update"))
	err := b.repo.AdminUpdate(mock.MockBiodata)
	assert.Error(b.T(), err)
}
func (b *BiodataUserRepoTestSuite) TestUserUpdate_Success() {
	expectedSQL := `UPDATE biodata SET user_credential_id = $2, full_name = $3, nik = $4, phone_number = $5, occupation = $6, place_of_birth = $7, date_of_birth = $8, postal_code = $9, is_eglible = $10, status_update = $11, additional_information = $12 WHERE id = $1`
	b.mockSQL.ExpectExec(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockBiodata.Id, mock.MockBiodata.UserCredential.Id, mock.MockBiodata.NamaLengkap, mock.MockBiodata.Nik, mock.MockBiodata.NomorTelepon, mock.MockBiodata.Pekerjaan, mock.MockBiodata.TempatLahir, mock.MockBiodata.TanggalLahir, mock.MockBiodata.KodePos, mock.MockBiodata.IsAglible, mock.MockBiodata.StatusUpdate, mock.MockBiodata.Information).WillReturnResult(sqlmock.NewResult(1, 1))
	err := b.repo.UserUpdate(mock.MockBiodata)
	assert.Nil(b.T(), err)
}
func (b *BiodataUserRepoTestSuite) TestUserUpdate_Failed() {
	expectedSQL := `UPDATE biodata SET user_credential_id = $2, full_name = $3, nik = $4, phone_number = $5, occupation = $6, place_of_birth = $7, date_of_birth = $8, postal_code = $9, is_eglible = $10, status_update = $11, additional_information = $12 WHERE id = $1`
	b.mockSQL.ExpectExec(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("failed update"))
	err := b.repo.UserUpdate(mock.MockBiodata)
	assert.Error(b.T(), err)
}
func (b *BiodataUserRepoTestSuite) TestFindByUcId_Success() {
	rows := sqlmock.NewRows([]string{"id", "user_credential_id", "user_credential_username", "user_credential_email", "user_credential_role", "user_credential_is_active", "user_credential_virtual_account_number", "full_name", "nik", "phone_number", "occupation", "place_of_birth", "date_of_birth", "postal_code", "is_eglible", "status_update", "additional_information"})
	rows.AddRow(mock.MockBiodata.Id, mock.MockBiodata.UserCredential.Id, mock.MockBiodata.UserCredential.Username, mock.MockBiodata.UserCredential.Email, mock.MockBiodata.UserCredential.Role, mock.MockBiodata.UserCredential.IsActive, mock.MockBiodata.UserCredential.VANumber, mock.MockBiodata.NamaLengkap, mock.MockBiodata.Nik, mock.MockBiodata.NomorTelepon, mock.MockBiodata.Pekerjaan, mock.MockBiodata.TempatLahir, mock.MockBiodata.TanggalLahir, mock.MockBiodata.KodePos, mock.MockBiodata.IsAglible, mock.MockBiodata.StatusUpdate, mock.MockBiodata.Information)
	expectedSQL := `SELECT b.id, u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code, b.is_eglible, b.status_update, b.additional_information FROM biodata b JOIN user_credential u ON u.id = b.user_credential_id WHERE b.user_credential_id = $1;`
	b.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(mock.MockBiodata.UserCredential.Id).WillReturnRows(rows)
	br, err := b.repo.FindByUcId(mock.MockBiodata.UserCredential.Id)
	assert.Nil(b.T(), err)
	assert.Equal(b.T(), mock.MockBiodata.UserCredential.Id, br.UserCredential.Id)
	assert.Equal(b.T(), mock.MockBiodata.UserCredential.Username, br.UserCredential.Username)
}
func (b *BiodataUserRepoTestSuite) TestFindByUcId_Failed() {
	expectedSQL := `SELECT b.id, u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code, b.is_eglible, b.status_update, b.additional_information FROM biodata b JOIN user_credential u ON u.id = b.user_credential_id WHERE b.user_credential_id = $1;`
	b.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("failed find"))
	br, err := b.repo.FindByUcId(mock.MockBiodata.UserCredential.Id)
	assert.Error(b.T(), err)
	assert.Equal(b.T(), dto.BiodataResponse{}, br)
}
func (b *BiodataUserRepoTestSuite) TestFindUserUpdated_Success() {
	rows := sqlmock.NewRows([]string{"id", "user_credential_id", "user_credential_username", "user_credential_email", "user_credential_role", "user_credential_is_active", "user_credential_virtual_account_number", "full_name", "nik", "phone_number", "occupation", "place_of_birth", "date_of_birth", "postal_code", "is_eglible", "status_update", "additional_information"})
	for _, row := range mock.MockBiodataResponses {
		rows.AddRow(row.Id, row.UserCredential.Id, row.UserCredential.Username, row.UserCredential.Email, row.UserCredential.Role, row.UserCredential.IsActive, row.UserCredential.VaNumber, row.NamaLengkap, row.Nik, row.NomorTelepon, row.Pekerjaan, row.TempatLahir, row.TanggalLahir, row.KodePos, row.IsAglible, row.StatusUpdate, row.Information)
	}
	expectedSQL := `SELECT b.id, u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code, b.is_eglible, b.status_update, b.additional_information FROM biodata b JOIN user_credential u ON u.id = b.user_credential_id WHERE b.status_update = true AND b.is_eglible = false`
	b.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnRows(rows)
	br, err := b.repo.FindUserUpdated()
	assert.Nil(b.T(), err)
	assert.Equal(b.T(), 1, len(br))
}
func (b *BiodataUserRepoTestSuite) TestFindUserUpdated_Failed() {
	expectedSQL := `SELECT b.id, u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code, b.is_eglible, b.status_update, b.additional_information FROM biodata b JOIN user_credential u ON u.id = b.user_credential_id WHERE b.status_update = true AND b.is_eglible = false`
	b.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("failed find"))
	br, err := b.repo.FindUserUpdated()
	assert.Error(b.T(), err)
	assert.Nil(b.T(), br)
}
func (b *BiodataUserRepoTestSuite) TestPagging_Success() {
	rows := sqlmock.NewRows([]string{"id", "user_credential_id", "user_credential_username", "user_credential_email", "user_credential_role", "user_credential_is_active", "user_credential_virtual_account_number", "full_name", "nik", "phone_number", "occupation", "place_of_birth", "date_of_birth", "postal_code", "is_eglible", "status_update", "additional_information"})
	for _, row := range mock.MockBiodataResponses {
		rows.AddRow(row.Id, row.UserCredential.Id, row.UserCredential.Username, row.UserCredential.Email, row.UserCredential.Role, row.UserCredential.IsActive, row.UserCredential.VaNumber, row.NamaLengkap, row.Nik, row.NomorTelepon, row.Pekerjaan, row.TempatLahir, row.TanggalLahir, row.KodePos, row.IsAglible, row.StatusUpdate, row.Information)
	}
	expectedSQL := `SELECT b.id, u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code, b.is_eglible, b.status_update, b.additional_information FROM biodata b JOIN user_credential u ON u.id = b.user_credential_id LIMIT $2 OFFSET $1;`
	b.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	b.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM biodata`)).WillReturnRows(rowCount)

	br, p, err := b.repo.Pagging(mock.MockPageReq)
	assert.Nil(b.T(), err)
	assert.Equal(b.T(), 1, len(br))
	assert.Equal(b.T(), 1, p.TotalRows)
}
func (b *BiodataUserRepoTestSuite) TestPagging_Failed() {
	// error select paging
	expectedSQL := `SELECT b.id, u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code, b.is_eglible, b.status_update, b.additional_information FROM biodata b JOIN user_credential u ON u.id = b.user_credential_id LIMIT $2 OFFSET $1;`
	b.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WillReturnError(errors.New("failed"))
	br, p, err := b.repo.Pagging(dto.PageRequest{})
	assert.Error(b.T(), err)
	assert.Nil(b.T(), br)
	assert.Equal(b.T(), 0, p.TotalRows)
	// error select count
	rows := sqlmock.NewRows([]string{"id", "user_credential_id", "user_credential_username", "user_credential_email", "user_credential_role", "user_credential_is_active", "user_credential_virtual_account_number", "full_name", "nik", "phone_number", "occupation", "place_of_birth", "date_of_birth", "postal_code", "is_eglible", "status_update", "additional_information"})
	for _, row := range mock.MockBiodataResponses {
		rows.AddRow(row.Id, row.UserCredential.Id, row.UserCredential.Username, row.UserCredential.Email, row.UserCredential.Role, row.UserCredential.IsActive, row.UserCredential.VaNumber, row.NamaLengkap, row.Nik, row.NomorTelepon, row.Pekerjaan, row.TempatLahir, row.TanggalLahir, row.KodePos, row.IsAglible, row.StatusUpdate, row.Information)
	}
	expectedSQL = `SELECT b.id, u.id, u.username, u.email, u.role, u.is_active, u.virtual_account_number, b.full_name, b.nik, b.phone_number, b.occupation, b.place_of_birth, b.date_of_birth, b.postal_code, b.is_eglible, b.status_update, b.additional_information FROM biodata b JOIN user_credential u ON u.id = b.user_credential_id LIMIT $2 OFFSET $1;`
	b.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedSQL)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	b.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM biodata`)).WillReturnError(errors.New("failed"))

	br, p, err = b.repo.Pagging(mock.MockPageReq)
	assert.Error(b.T(), err)
	assert.Nil(b.T(), br)
	assert.Equal(b.T(), 0, p.TotalRows)
}
