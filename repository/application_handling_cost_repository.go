package repository

import (
	"database/sql"
	"math"
	"polen/model"
	"polen/model/dto"
)

type AppHandlingCost interface {
	CreateNew(payload model.AppHandlingCost) error
	FindById(id string) (model.AppHandlingCost, error)
	Pagging(payload dto.PageRequest) ([]model.AppHandlingCost, dto.Paging, error)
	Update(payload model.AppHandlingCost) error
	DeleteById(id string) error
}
type appHandlingCostRepository struct {
	db *sql.DB
}

// CreateNew implements AppHandlingCost.
func (a *appHandlingCostRepository) CreateNew(payload model.AppHandlingCost) error {
	_, err := a.db.Exec("INSERT INTO application_handling_cost VALUES ($1, $2, $3, $4)", payload.Id, payload.Name, payload.Nominal, payload.Unit)
	if err != nil {
		return err
	}
	return nil
}

// DeleteById implements AppHandlingCost.
func (a *appHandlingCostRepository) DeleteById(id string) error {
	_, err := a.db.Exec("DELETE FROM application_handling_cost WHERE id = $1", id)
	return err
}

// FindById implements AppHandlingCost.
func (a *appHandlingCostRepository) FindById(id string) (model.AppHandlingCost, error) {
	row := a.db.QueryRow(`SELECT 
		id, 
		name,
		nominal,
		unit
	FROM 
		application_handling_cost
	WHERE id = $1`, id)
	AppHandlingCost := model.AppHandlingCost{}
	err := row.Scan(
		&AppHandlingCost.Id,
		&AppHandlingCost.Name,
		&AppHandlingCost.Nominal,
		&AppHandlingCost.Unit,
	)
	if err != nil {
		return model.AppHandlingCost{}, err
	}
	return AppHandlingCost, nil
}

// Pagging implements AppHandlingCost.
func (a *appHandlingCostRepository) Pagging(payload dto.PageRequest) ([]model.AppHandlingCost, dto.Paging, error) {
	if payload.Page < 0 {
		payload.Page = 1
	}
	query := `SELECT 
		id, 
		name,
		nominal,
		unit
	FROM 
		application_handling_cost
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
	var data []model.AppHandlingCost
	for rows.Next() {
		var AppHandlingCost model.AppHandlingCost
		err := rows.Scan(
			&AppHandlingCost.Id,
			&AppHandlingCost.Name,
			&AppHandlingCost.Nominal,
			&AppHandlingCost.Unit,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		data = append(data, AppHandlingCost)
	}

	var count int
	row := a.db.QueryRow("SELECT COUNT(id) FROM application_handling_cost")
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

// Update implements AppHandlingCost.
func (a *appHandlingCostRepository) Update(payload model.AppHandlingCost) error {
	_, err := a.db.Exec(`
		UPDATE 
			application_handling_cost
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

func NewAppHandlingCostRepository(db *sql.DB) AppHandlingCost {
	return &appHandlingCostRepository{
		db: db,
	}
}
