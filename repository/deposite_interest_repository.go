package repository

import (
	"database/sql"
	"math"
	"polen/model"
	"polen/model/dto"
)

type DepositeInterest interface {
	Save(payload model.DepositeInterest) error
	Pagging(payload dto.PageRequest) ([]dto.DepositeInterestRequest, dto.Paging, error)
	Update(payload dto.DepositeInterestRequest) error
	DeleteById(id string) error
	FindById(id string) (dto.DepositeInterestRequest, error)
}

type depositeInterest struct {
	db *sql.DB
}

func (d *depositeInterest) Pagging(payload dto.PageRequest) ([]dto.DepositeInterestRequest, dto.Paging, error) {
	query := `SELECT 
		id, 
		interest_rate,
		tax_rate,
		duration_mounth
	FROM 
		deposit_interest
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
	var data []dto.DepositeInterestRequest
	for rows.Next() {
		var datum dto.DepositeInterestRequest
		err := rows.Scan(
			&datum.Id,
			&datum.InterestRate,
			&datum.TaxRate,
			&datum.DurationMounth,
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

// DeleteById implements DepositeInterest.
func (d *depositeInterest) DeleteById(id string) error {
	_, err := d.db.Exec("DELETE FROM deposit_interest WHERE id = $1", id)
	return err
}

// FindById implements DepositeInterest.
func (d *depositeInterest) FindById(id string) (dto.DepositeInterestRequest, error) {
	row := d.db.QueryRow(`SELECT 
		id, 
		interest_rate,
		tax_rate,
		duration_mounth
	FROM 
		deposit_interest 
	WHERE id =$1`, id)
	deposite := dto.DepositeInterestRequest{}
	err := row.Scan(
		&deposite.Id,
		&deposite.InterestRate,
		&deposite.TaxRate,
		&deposite.DurationMounth,
	)
	if err != nil {
		return dto.DepositeInterestRequest{}, err
	}
	return deposite, nil
}

// Update implements DepositeInterest.
func (d *depositeInterest) Update(payload dto.DepositeInterestRequest) error {
	_, err := d.db.Exec(`
		UPDATE 
			deposit_interest 
		SET 
			interest_rate = $2, 
			tax_rate = $3,
			duration_mounth = $4
		WHERE 
			id = $1;`,
		payload.Id,
		payload.InterestRate,
		payload.TaxRate,
		payload.DurationMounth,
	)
	return err
}

// Save implements DepositeIntereset.
func (d *depositeInterest) Save(payload model.DepositeInterest) error {
	_, err := d.db.Exec("INSERT INTO deposit_interest VALUES ($1, $2, $3, $4, $5)", payload.Id, payload.CreateDate, payload.InterestRate, payload.TaxRate, payload.DurationMounth)
	if err != nil {
		return err
	}
	return nil
}

func NewDepositeInterestRepository(db *sql.DB) DepositeInterest {
	return &depositeInterest{
		db: db,
	}
}
