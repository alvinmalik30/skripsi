package repository

import (
	"database/sql"
	"math"
	"polen/model"
	"polen/model/dto"
	"time"
)

type DepositeRepository interface {
	CreateDeposite(payload dto.DepositeDto) error
	FindByUcId(id string) (dto.DepositeByUserResponse, error)
	FindById(id string) (dto.DepositeByIdResponse, error)
	Pagging(payload dto.PageRequest) ([]dto.Deposite, dto.Paging, error)
	Update() error
}

type depositeRepository struct {
	db   *sql.DB
	user BiodataUser
}

// Update implements DepositeRepository.
func (d *depositeRepository) Update() error {
	rows, err := d.db.Query(`
		SELECT
			id,
			user_credential_id,
			tax,
			status,
			total_return
		FROM
			deposit
		WHERE
			status = true
		AND
			maturity_date < now();
	`)
	if err != nil {
		return err
	}

	defer rows.Close()
	var data []model.Deposite
	for rows.Next() {
		var datum model.Deposite
		err := rows.Scan(
			&datum.Id,
			&datum.UserCredentialId,
			&datum.Tax,
			&datum.Status,
			&datum.TotalReturn,
		)
		if err != nil {
			return err
		}
		data = append(data, datum)
	}

	// update data
	for _, v := range data {
		tx, _ := d.db.Begin()
		_, err = tx.Exec(
			`
			UPDATE 
				deposit
			SET 
				status = false
			WHERE 
				id = $1;
			`,
			v.Id,
		)
		if err != nil {
			tx.Rollback()
		}
		_, err = tx.Exec(
			// tambah saldo user
			`
			UPDATE saldo
			SET total_saving = total_saving + $1
			WHERE user_credential_id = $2;
			`,
			v.TotalReturn,
			v.UserCredentialId,
		)
		if err != nil {
			tx.Rollback()
		}
		// kurangi sado perusahaan
		_, err = tx.Exec(
			`
			UPDATE saldo
			SET total_saving = total_saving - $1 - $2
			WHERE user_credential_id = $3;
			`,
			v.TotalReturn,
			v.Tax,
			"456",
		)
		if err != nil {
			tx.Rollback()
		}
		// tambah saldo pajak
		_, err = tx.Exec(
			`
			UPDATE saldo
			SET total_saving = total_saving + $1
			WHERE user_credential_id = $2;
			`,
			v.Tax,
			"tax",
		)
		if err != nil {
			tx.Rollback()
		}
		tx.Commit()
	}
	return nil
}

// Pagging implements DepositeRepository.
func (d *depositeRepository) Pagging(payload dto.PageRequest) ([]dto.Deposite, dto.Paging, error) {
	query := `
	SELECT
		id,
		deposit_amount,
		interest_rate,
		tax_rate,
		duration,
		created_date,
		maturity_date,
		status,
		gross_profit,
		tax,
		net_profit,
		total_return
	FROM
		deposit
	LIMIT 
		$2 
	OFFSET 
		$1;
	`
	rows, err := d.db.Query(query, (payload.Page-1)*payload.Page, payload.Size)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()
	var data []dto.Deposite
	for rows.Next() {
		var datum dto.Deposite
		err := rows.Scan(
			&datum.Id,
			&datum.DepositeAmount,
			&datum.InterestRate,
			&datum.TaxRate,
			&datum.DurationMounth,
			&datum.CreateDate,
			&datum.MaturityDate,
			&datum.Status,
			&datum.GrossProfit,
			&datum.Tax,
			&datum.NetProfit,
			&datum.TotalReturn,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		data = append(data, datum)
	}

	var count int
	row := d.db.QueryRow("SELECT COUNT(id) FROM deposit_interest")
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

// FindById implements DepositeRepository.
func (d *depositeRepository) FindById(id string) (dto.DepositeByIdResponse, error) {
	// get data
	var result dto.DepositeByIdResponse
	var ucid string
	row := d.db.QueryRow(`
		SELECT
			id,
			user_credential_id,
			deposit_amount,
			interest_rate,
			tax_rate,
			duration,
			created_date,
			maturity_date,
			status,
			gross_profit,
			tax,
			net_profit,
			total_return
		FROM
			deposit
		WHERE
			id = $1;
	`, id)
	err := row.Scan(
		&result.Deposite.Id,
		&ucid,
		&result.Deposite.DepositeAmount,
		&result.Deposite.InterestRate,
		&result.Deposite.TaxRate,
		&result.Deposite.DurationMounth,
		&result.Deposite.CreateDate,
		&result.Deposite.MaturityDate,
		&result.Deposite.Status,
		&result.Deposite.GrossProfit,
		&result.Deposite.Tax,
		&result.Deposite.NetProfit,
		&result.Deposite.TotalReturn,
	)
	if err != nil {
		return dto.DepositeByIdResponse{}, err
	}
	// get biodata user
	bio, err := d.user.FindByUcId(ucid)
	if err != nil {
		return dto.DepositeByIdResponse{}, err
	}
	result.BioUser = bio
	return result, err
}

// FindByUcId implements DepositeRepository.
func (d *depositeRepository) FindByUcId(id string) (dto.DepositeByUserResponse, error) {
	// get biodata user
	bio, err := d.user.FindByUcId(id)
	if err != nil {
		return dto.DepositeByUserResponse{}, err
	}
	var result dto.DepositeByUserResponse
	result.BioUser = bio
	// get data
	rows, err := d.db.Query(`
		SELECT
			id,
			deposit_amount,
			interest_rate,
			tax_rate,
			duration,
			created_date,
			maturity_date,
			status,
			gross_profit,
			tax,
			net_profit,
			total_return
		FROM
			deposit
		WHERE
			user_credential_id = $1;
	`, id)
	if err != nil {
		return dto.DepositeByUserResponse{}, err
	}
	var data []dto.Deposite
	for rows.Next() {
		var datum dto.Deposite
		err := rows.Scan(
			&datum.Id,
			&datum.DepositeAmount,
			&datum.InterestRate,
			&datum.TaxRate,
			&datum.DurationMounth,
			&datum.CreateDate,
			&datum.MaturityDate,
			&datum.Status,
			&datum.GrossProfit,
			&datum.Tax,
			&datum.NetProfit,
			&datum.TotalReturn,
		)
		if err != nil {
			return dto.DepositeByUserResponse{}, err
		}
		data = append(data, datum)
	}
	result.Deposite = data
	return result, nil
}

// CreateDeposite implements DepositeRepository.
func (d *depositeRepository) CreateDeposite(payload dto.DepositeDto) error {
	tx, err := d.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO deposit VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)",
		payload.Id,
		payload.UserCredential.Id,
		payload.DepositeAmount,
		payload.InterestRate.InterestRate,
		payload.InterestRate.TaxRate,
		payload.InterestRate.DurationMounth,
		time.Now(),
		payload.MaturityDate,
		payload.Status,
		payload.GrossProfit,
		payload.Tax,
		payload.NetProfit,
		payload.TotalReturn,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		`
		UPDATE saldo
		SET total_saving = total_saving - $1
		WHERE user_credential_id = $2;
		`,
		payload.DepositeAmount,
		payload.UserCredential.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		`
		UPDATE saldo
		SET total_saving = total_saving + $1
		WHERE user_credential_id = 456;
		`,
		payload.DepositeAmount,
		payload.UserCredential.Id,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return err
}

func NewDepositeRepository(db *sql.DB, user BiodataUser) DepositeRepository {
	return &depositeRepository{
		db:   db,
		user: user,
	}
}
