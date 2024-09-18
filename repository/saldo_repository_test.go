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

type SaldoRepoTestSuite struct {
	suite.Suite
	mockDB  *sql.DB
	mockSQL sqlmock.Sqlmock
	repo    SaldoRepository
}

func (s *SaldoRepoTestSuite) SetupTest() {
	d, s2, err := sqlmock.New()
	assert.NoError(s.T(), err)
	s.mockDB = d
	s.mockSQL = s2
	s.repo = NewSaldoRepository(s.mockDB)
}
func TestSaldoRepoTestSuite(t *testing.T) {
	suite.Run(t, new(SaldoRepoTestSuite))
}
func (s *SaldoRepoTestSuite) TestFindByIdUser_Success() {
	r := sqlmock.NewRows([]string{"id", "total_saving"})
	r.AddRow(mock.MockSaldo.Id, mock.MockSaldo.Total)
	expectedQuery := `SELECT id, total_saving FROM saldo WHERE user_credential_id = $1`
	s.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnRows(r)
	ahc, err := s.repo.FindByIdUser(mock.MockSaldo.UcId)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), mock.MockSaldo.UcId, ahc.UcId)
}
func (s *SaldoRepoTestSuite) TestFindByIdUser_Failed() {
	expectedQuery := `SELECT id, total_saving FROM saldo WHERE user_credential_id = $1`
	s.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WillReturnError(errors.New("error"))
	ahc, err := s.repo.FindByIdUser(mock.MockSaldo.UcId)
	assert.Error(s.T(), err)
	assert.Equal(s.T(), dto.Saldo{}, ahc)
}
func (s *SaldoRepoTestSuite) TestPagging_Success() {
	rows := sqlmock.NewRows([]string{"id", "total_saving"})
	for _, sd := range mock.MockSaldoDatas {
		rows.AddRow(sd.Id, sd.Total)
	}
	expectedQuery := `SELECT id, total_saving FROM saldo LIMIT $2 OFFSET $1;`
	s.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	s.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM saldo`)).WillReturnRows(rowCount)

	uc, p, err := s.repo.Pagging(mock.MockPageReq)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(uc))
	assert.Equal(s.T(), 1, p.TotalRows)
}
func (s *SaldoRepoTestSuite) TestPagging_Fail() {
	// error select paging
	expectedQuery := `SELECT id, total_saving FROM saldo LIMIT $2 OFFSET $1;`
	s.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnError(errors.New("failed"))
	uc, p, err := s.repo.Pagging(dto.PageRequest{})
	assert.Error(s.T(), err)
	assert.Nil(s.T(), uc)
	assert.Equal(s.T(), 0, p.TotalRows)
	// error select count
	rows := sqlmock.NewRows([]string{"id", "total_saving"})
	for _, sd := range mock.MockSaldoDatas {
		rows.AddRow(sd.Id, sd.Total)
	}
	expectedQuery = `SELECT id, total_saving FROM saldo LIMIT $2 OFFSET $1;`
	s.mockSQL.ExpectQuery(regexp.QuoteMeta(expectedQuery)).WithArgs(((mock.MockPageReq.Page - 1) * mock.MockPageReq.Size), mock.MockPageReq.Size).WillReturnRows(rows)

	rowCount := sqlmock.NewRows([]string{"count"})
	rowCount.AddRow(1)
	s.mockSQL.ExpectQuery(regexp.QuoteMeta(`SELECT COUNT(id) FROM saldo`)).WillReturnError(errors.New("failed"))
	uc, p, err = s.repo.Pagging(mock.MockPageReq)
	assert.Error(s.T(), err)
	assert.Nil(s.T(), uc)
	assert.Equal(s.T(), 0, p.TotalRows)
}
