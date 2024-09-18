package repository

import (
	"database/sql"
	"math"
	"polen/model"
	"polen/model/dto"
)

type LoanInterest interface {
	CreateNew(payload model.LoanInterest) error
	FindById(id string) (model.LoanInterest, error)
	Pagging(payload dto.PageRequest) ([]model.LoanInterest, dto.Paging, error)
	Update(payload model.LoanInterest) error
	DeleteById(id string) error
}
type loanInterestRepository struct {
	db *sql.DB
}

// Pagging implements LoanInterest.
func (l *loanInterestRepository) Pagging(payload dto.PageRequest) ([]model.LoanInterest, dto.Paging, error) {
	if payload.Page < 0 {
		payload.Page = 1
	}
	query := `SELECT 
		id, 
		duration_months,
		loan_interest_rate
	FROM 
		loan_interest
	LIMIT 
		$2 
	OFFSET 
		$1;
	`
	rows, err := l.db.Query(query, (payload.Page-1)*payload.Page, payload.Size)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()
	var data []model.LoanInterest
	for rows.Next() {
		var loan model.LoanInterest
		err := rows.Scan(
			&loan.Id,
			&loan.DurationMonths,
			&loan.LoanInterestRate,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		data = append(data, loan)
	}

	var count int
	row := l.db.QueryRow("SELECT COUNT(id) FROM loan_interest")
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

// CreateNew implements LoanInterest.
func (l *loanInterestRepository) CreateNew(payload model.LoanInterest) error {
	_, err := l.db.Exec("INSERT INTO loan_interest VALUES ($1, $2, $3)", payload.Id, payload.DurationMonths, payload.LoanInterestRate)
	if err != nil {
		return err
	}
	return nil
}

// DeleteById implements LoanInterest.
func (l *loanInterestRepository) DeleteById(id string) error {
	_, err := l.db.Exec("DELETE FROM loan_interest WHERE id = $1", id)
	return err
}

// FindById implements LoanInterest.
func (l *loanInterestRepository) FindById(id string) (model.LoanInterest, error) {
	row := l.db.QueryRow(`SELECT 
		id, 
		duration_months,
		loan_interest_rate
	FROM 
		loan_interest 
	WHERE id =$1`, id)
	LoanInterest := model.LoanInterest{}
	err := row.Scan(
		&LoanInterest.Id,
		&LoanInterest.DurationMonths,
		&LoanInterest.LoanInterestRate,
	)
	if err != nil {
		return model.LoanInterest{}, err
	}
	return LoanInterest, nil
}

// Update implements LoanInterest.
func (l *loanInterestRepository) Update(payload model.LoanInterest) error {
	_, err := l.db.Exec(`
		UPDATE 
			loan_interest 
		SET 
			duration_months = $2, 
			loan_interest_rate = $3
		WHERE 
			id = $1;`,
		payload.Id,
		payload.DurationMonths,
		payload.LoanInterestRate,
	)
	if err != nil {
		return err
	}
	return nil
}

func NewLoanInterestRepository(db *sql.DB) LoanInterest {
	return &loanInterestRepository{
		db: db,
	}
}
