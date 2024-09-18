package repository

import (
	"database/sql"
	"math"
	"polen/model"
	"polen/model/dto"
)

type BiodataUser interface {
	FindUserUpdated() ([]dto.BiodataResponse, error)
	FindByUcId(id string) (dto.BiodataResponse, error)
	UserUpdate(payload model.BiodataUser) error
	AdminUpdate(payload model.BiodataUser) error
	Pagging(payload dto.PageRequest) ([]dto.BiodataResponse, dto.Paging, error)
}

type biodataUserRepository struct {
	db *sql.DB
}

// Pagging implements BiodataUser.
func (bio *biodataUserRepository) Pagging(payload dto.PageRequest) ([]dto.BiodataResponse, dto.Paging, error) {
	query := `SELECT 
		b.id, 
			u.id, 
			u.username, 
			u.email, 
			u.role, 
			u.is_active, 
			u.virtual_account_number,
		b.full_name, 
		b.nik, 
		b.phone_number, 
		b.occupation, 
		b.place_of_birth, 
		b.date_of_birth, 
		b.postal_code,
		b.is_eglible,
		b.status_update,
		b.additional_information
	FROM 
		biodata b 
	JOIN 
		user_credential u 
	ON 
		u.id = b.user_credential_id
	LIMIT 
		$2 
	OFFSET 
		$1;
	`
	rows, err := bio.db.Query(query, (payload.Page-1)*payload.Page, payload.Size)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()
	var data []dto.BiodataResponse
	for rows.Next() {
		var datum dto.BiodataResponse
		err := rows.Scan(
			&datum.Id,
			&datum.UserCredential.Id,
			&datum.UserCredential.Username,
			&datum.UserCredential.Email,
			&datum.UserCredential.Role,
			&datum.UserCredential.IsActive,
			&datum.UserCredential.VaNumber,
			&datum.NamaLengkap,
			&datum.Nik,
			&datum.NomorTelepon,
			&datum.Pekerjaan,
			&datum.TempatLahir,
			&datum.TanggalLahir,
			&datum.KodePos,
			&datum.IsAglible,
			&datum.StatusUpdate,
			&datum.Information,
		)
		datum.UserCredential.Password = "********"
		if err != nil {
			return nil, dto.Paging{}, err
		}
		data = append(data, datum)
	}

	var count int
	row := bio.db.QueryRow("SELECT COUNT(id) FROM biodata")
	if err := row.Scan(&count); err != nil {
		return nil, dto.Paging{}, err
	}

	paging := dto.Paging{
		Page:       payload.Page,
		Size:       payload.Size,
		TotalRows:  count,
		TotalPages: int(math.Ceil(float64(count) / float64(payload.Size))),
	}

	return data, paging, nil
}

// AdminUpdate implements BiodataUser.
func (bio *biodataUserRepository) AdminUpdate(payload model.BiodataUser) error {
	_, err := bio.db.Exec(`
	UPDATE biodata SET 
			is_eglible = $2,
			status_update = $3,
			additional_information = $4
		WHERE id = $1`,
		payload.Id,
		payload.IsAglible,
		payload.StatusUpdate,
		payload.Information,
	)
	return err
}

// UserUpdate implements BiodataUser.
func (bio *biodataUserRepository) UserUpdate(payload model.BiodataUser) error {
	_, err := bio.db.Exec(`
	UPDATE biodata SET 
			user_credential_id = $2, 
			full_name = $3, 
			nik = $4, 
			phone_number = $5, 
			occupation = $6, 
			place_of_birth = $7, 
			date_of_birth = $8, 
			postal_code = $9,
			is_eglible = $10,
			status_update = $11,
			additional_information = $12
		WHERE id = $1`,
		payload.Id,
		payload.UserCredential.Id,
		payload.NamaLengkap,
		payload.Nik,
		payload.NomorTelepon,
		payload.Pekerjaan,
		payload.TempatLahir,
		payload.TanggalLahir,
		payload.KodePos,
		payload.IsAglible,
		payload.StatusUpdate,
		payload.Information,
	)
	return err
}

// FindByUcId implements BiodataUser.
func (bio *biodataUserRepository) FindByUcId(id string) (dto.BiodataResponse, error) {
	row := bio.db.QueryRow(`
	SELECT 
			b.id, 
				u.id, 
				u.username, 
				u.email, 
				u.role, 
				u.is_active, 
				u.virtual_account_number,
			b.full_name, 
			b.nik, 
			b.phone_number, 
			b.occupation, 
			b.place_of_birth, 
			b.date_of_birth, 
			b.postal_code,
			b.is_eglible,
			b.status_update,
			b.additional_information
		FROM 
			biodata b 
		JOIN 
			user_credential u 
		ON 
			u.id = b.user_credential_id
		WHERE
			b.user_credential_id = $1;
`, id)
	biodata := dto.BiodataResponse{}
	biodata.UserCredential.Password = "********"
	err := row.Scan(
		&biodata.Id,
		&biodata.UserCredential.Id,
		&biodata.UserCredential.Username,
		&biodata.UserCredential.Email,
		&biodata.UserCredential.Role,
		&biodata.UserCredential.IsActive,
		&biodata.UserCredential.VaNumber,
		&biodata.NamaLengkap,
		&biodata.Nik,
		&biodata.NomorTelepon,
		&biodata.Pekerjaan,
		&biodata.TempatLahir,
		&biodata.TanggalLahir,
		&biodata.KodePos,
		&biodata.IsAglible,
		&biodata.StatusUpdate,
		&biodata.Information,
	)
	if err != nil {
		return dto.BiodataResponse{}, err
	}
	return biodata, nil
}

// FindByUcId implements BiodataUser.
func (bio *biodataUserRepository) FindUserUpdated() ([]dto.BiodataResponse, error) {
	rows, err := bio.db.Query(`
	SELECT 
			b.id, 
				u.id, 
				u.username, 
				u.email, 
				u.role, 
				u.is_active,
				u.virtual_account_number,
			b.full_name, 
			b.nik, 
			b.phone_number, 
			b.occupation, 
			b.place_of_birth, 
			b.date_of_birth, 
			b.postal_code,
			b.is_eglible,
			b.status_update,
			b.additional_information
		FROM 
			biodata b 
		JOIN 
			user_credential u 
		ON 
			u.id = b.user_credential_id
		WHERE
			b.status_update = true
			AND
			b.is_eglible = false
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var biodatas []dto.BiodataResponse
	for rows.Next() {
		biodata := dto.BiodataResponse{}
		biodata.UserCredential.Password = "********"
		err := rows.Scan(
			&biodata.Id,
			&biodata.UserCredential.Id,
			&biodata.UserCredential.Username,
			&biodata.UserCredential.Email,
			&biodata.UserCredential.Role,
			&biodata.UserCredential.IsActive,
			&biodata.UserCredential.VaNumber,
			&biodata.NamaLengkap,
			&biodata.Nik,
			&biodata.NomorTelepon,
			&biodata.Pekerjaan,
			&biodata.TempatLahir,
			&biodata.TanggalLahir,
			&biodata.KodePos,
			&biodata.IsAglible,
			&biodata.StatusUpdate,
			&biodata.Information,
		)
		if err != nil {
			return nil, err
		}
		biodatas = append(biodatas, biodata)
	}
	return biodatas, nil
}

func NewBiodataUserRepository(db *sql.DB) BiodataUser {
	return &biodataUserRepository{db: db}
}
