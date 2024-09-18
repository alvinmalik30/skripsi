package repository

import (
	"database/sql"
	"math"
	"polen/model/dto"
)

type SaldoRepository interface {
	FindByIdUser(id string) (dto.Saldo, error)
	Pagging(payload dto.PageRequest) ([]dto.Saldo, dto.Paging, error)
}

type saldoRepository struct {
	db *sql.DB
}

// FindByIdUser implements SaldoRepository.
func (s *saldoRepository) FindByIdUser(id string) (dto.Saldo, error) {
	row := s.db.QueryRow(`
		SELECT 
			id, 
			total_saving
		FROM 
			saldo 
		WHERE 
			user_credential_id = $1`, id)
	var data dto.Saldo
	data.UcId = id
	err := row.Scan(
		&data.Id,
		&data.Total,
	)
	if err != nil {
		return dto.Saldo{}, err
	}
	return data, nil
}

// Pagging implements SaldoRepository.
func (s *saldoRepository) Pagging(payload dto.PageRequest) ([]dto.Saldo, dto.Paging, error) {
	// limit Size, Offset (page - 1) * size
	if payload.Page < 0 {
		payload.Page = 1
	}
	query := `
	SELECT 
		id, 
		total_saving
	FROM 
		saldo 
	LIMIT 
		$2 
	OFFSET 
		$1;
	`
	rows, err := s.db.Query(query, (payload.Page-1)*payload.Page, payload.Size)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()
	var data []dto.Saldo
	for rows.Next() {
		var datum dto.Saldo
		err := rows.Scan(
			&datum.Id,
			&datum.Total,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		data = append(data, datum)
	}

	var count int
	row := s.db.QueryRow("SELECT COUNT(id) FROM saldo")
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

func NewSaldoRepository(db *sql.DB) SaldoRepository {
	return &saldoRepository{
		db: db,
	}
}
