package repository

import (
	"database/sql"
	"math"
	"polen/model"
	"polen/model/dto"
	"time"
)

type UserRepository interface {
	Save(payload model.UserCredential, bioId string) error
	Saldo(payload model.UserCredential, idsaldo string, bioId string) error
	FindByUsername(username string) (model.UserCredential, error)
	Pagging(payload dto.PageRequest) ([]model.UserCredential, dto.Paging, error)
	FindById(id string) (model.UserCredential, error)
}

type userRepository struct {
	db *sql.DB
}

// pagging implements UserRepository.
func (u *userRepository) Pagging(payload dto.PageRequest) ([]model.UserCredential, dto.Paging, error) {
	query := `SELECT id, username, email, role, virtual_account_number, is_active FROM user_credential LIMIT $2 OFFSET $1`
	rows, err := u.db.Query(query, (payload.Page-1)*payload.Page, payload.Size)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()
	var data []model.UserCredential
	for rows.Next() {
		var datum model.UserCredential
		err := rows.Scan(
			&datum.Id,
			&datum.Username,
			&datum.Email,
			&datum.Role,
			&datum.VANumber,
			&datum.IsActive,
		)
		datum.Password = "********"
		if err != nil {
			return nil, dto.Paging{}, err
		}
		data = append(data, datum)
	}

	var count int
	row := u.db.QueryRow("SELECT COUNT(id) FROM user_credential")
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

// saldo implements UserRepository.
func (u *userRepository) Saldo(payload model.UserCredential, idsaldo string, bioId string) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO user_credential VALUES ($1, $2, $3, $4, $5, $6, $7)", payload.Id, payload.Username, payload.Email, payload.Password, payload.Role, payload.VANumber, true)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec("INSERT INTO biodata VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		bioId,
		payload.Id,
		"",
		"",
		"",
		"",
		"",
		time.DateOnly,
		"",
		false,
		false,
		"biodata is not updated",
	)
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO saldo VALUES ($1, $2, $3)", idsaldo, payload.Id, 0)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

// FindById implements UserRepository.
func (u *userRepository) FindById(id string) (model.UserCredential, error) {
	row := u.db.QueryRow(`
		SELECT 
			id, 
			username, 
			email,
			role, 
			virtual_account_number,
			is_active
		FROM 
			user_credential 
		WHERE 
			id =$1`, id)
	var userCredential model.UserCredential
	userCredential.Password = "**********"
	err := row.Scan(
		&userCredential.Id,
		&userCredential.Username,
		&userCredential.Email,
		&userCredential.Role,
		&userCredential.VANumber,
		&userCredential.IsActive,
	)
	if err != nil {
		return model.UserCredential{}, err
	}
	return userCredential, nil
}

// FindByUsername implements UserRepository.
func (u *userRepository) FindByUsername(username string) (model.UserCredential, error) {
	row := u.db.QueryRow("SELECT id, username, role, password FROM user_credential WHERE username = $1", username)
	var userCredential model.UserCredential
	err := row.Scan(&userCredential.Id, &userCredential.Username, &userCredential.Role, &userCredential.Password)
	if err != nil {
		return model.UserCredential{}, err
	}
	return userCredential, nil
}

// Save implements UserRepository.
func (u *userRepository) Save(payload model.UserCredential, bioId string) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO user_credential VALUES ($1, $2, $3, $4, $5, $6, $7)", payload.Id, payload.Username, payload.Email, payload.Password, payload.Role, payload.VANumber, true)
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO biodata VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)",
		bioId,
		payload.Id,
		"",
		"",
		"",
		"",
		"",
		time.DateOnly,
		"",
		false,
		false,
		"biodata is not updated",
	)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}
