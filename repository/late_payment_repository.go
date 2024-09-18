package repository

import (
	"database/sql"
	"math"
	"polen/model"
	"polen/model/dto"
)

type LatePaymentFee interface {
	CreateNew(payload model.LatePaymentFee) error
	FindById(id string) (model.LatePaymentFee, error)
	Pagging(payload dto.PageRequest) ([]model.LatePaymentFee, dto.Paging, error)
	Update(payload model.LatePaymentFee) error
	DeleteById(id string) error
}
type latePaymentFeeRepository struct {
	db *sql.DB
}

// CreateNew implements LatePaymentFee.
func (a *latePaymentFeeRepository) CreateNew(payload model.LatePaymentFee) error {
	_, err := a.db.Exec("INSERT INTO late_payment_fee VALUES ($1, $2, $3, $4)", payload.Id, payload.Name, payload.Nominal, payload.Unit)
	if err != nil {
		return err
	}
	return nil
}

// DeleteById implements LatePaymentFee.
func (a *latePaymentFeeRepository) DeleteById(id string) error {
	_, err := a.db.Exec("DELETE FROM late_payment_fee WHERE id = $1", id)
	return err
}

// FindById implements LatePaymentFee.
func (a *latePaymentFeeRepository) FindById(id string) (model.LatePaymentFee, error) {
	row := a.db.QueryRow(`SELECT 
		id, 
		name,
		nominal,
		unit
	FROM 
		late_payment_fee
	WHERE id =$1`, id)
	LatePaymentFee := model.LatePaymentFee{}
	err := row.Scan(
		&LatePaymentFee.Id,
		&LatePaymentFee.Name,
		&LatePaymentFee.Nominal,
		&LatePaymentFee.Unit,
	)
	if err != nil {
		return model.LatePaymentFee{}, err
	}
	return LatePaymentFee, nil
}

// Pagging implements LatePaymentFee.
func (a *latePaymentFeeRepository) Pagging(payload dto.PageRequest) ([]model.LatePaymentFee, dto.Paging, error) {
	if payload.Page < 0 {
		payload.Page = 1
	}
	query := `SELECT 
		id, 
		name,
		nominal,
		unit
	FROM 
		late_payment_fee
	LIMIT 
		$2 
	OFFSET 
		$1;
	`
	rows, err := a.db.Query(query, (payload.Page-1)*payload.Page, payload.Size)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()
	var data []model.LatePaymentFee
	for rows.Next() {
		var LatePaymentFee model.LatePaymentFee
		err := rows.Scan(
			&LatePaymentFee.Id,
			&LatePaymentFee.Name,
			&LatePaymentFee.Nominal,
			&LatePaymentFee.Unit,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		data = append(data, LatePaymentFee)
	}

	var count int
	row := a.db.QueryRow("SELECT COUNT(id) FROM late_payment_fee")
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

// Update implements LatePaymentFee.
func (a *latePaymentFeeRepository) Update(payload model.LatePaymentFee) error {
	_, err := a.db.Exec(`
		UPDATE 
			late_payment_fee
		SET 
			name = $2, 
			nominal = $3,
			unit =$4
		WHERE 
			id = $1;`,
		payload.Id,
		payload.Name,
		payload.Nominal,
		payload.Unit,
	)
	if err != nil {
		return err
	}
	return nil
}

func NewLatePaymentFeeRepository(db *sql.DB) LatePaymentFee {
	return &latePaymentFeeRepository{
		db: db,
	}
}
