package repository

import (
	"database/sql"
	"math"
	"polen/model"
	"polen/model/dto"
)

type TopUp interface {
	Save(payload model.TopUp) error
	FindByIdUser(id string) (dto.TopUpByUser, error)
	FindById(id string) (dto.TopUpById, error)
	Upload(payload model.TopUp) error
	ConfimUpload(payload model.TopUp) error
	NotConfimUpload(payload model.TopUp) error
	FindUploadedFile() ([]dto.TopUp, error)
	Pagging(payload dto.PageRequest) ([]dto.TopUp, dto.Paging, error)
}

type topUpRepository struct {
	db *sql.DB
}

// Pagging implements TopUp.
func (t *topUpRepository) Pagging(payload dto.PageRequest) ([]dto.TopUp, dto.Paging, error) {
	query := `
	SELECT 
		id,
		maturity_time, 
		top_up_amount,
		accepted_time,
		accepted_status,
		status_information,
		transfer_confirmation_recipt,
		recipt_file
	FROM 
		top_up
	LIMIT 
		$2 
	OFFSET 
		$1;
	`
	rows, err := t.db.Query(query, (payload.Page-1)*payload.Page, payload.Size)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	defer rows.Close()
	var topUps []dto.TopUp
	for rows.Next() {
		var topUp dto.TopUp
		err := rows.Scan(
			&topUp.Id,
			&topUp.MaturityTime,
			&topUp.TopUpAmount,
			&topUp.AcceptedTime,
			&topUp.Accepted,
			&topUp.Status,
			&topUp.TransferConfirmRecipe,
			&topUp.File,
		)
		if err != nil {
			return nil, dto.Paging{}, err
		}
		topUps = append(topUps, topUp)
	}

	var count int
	row := t.db.QueryRow("SELECT COUNT(id) FROM top_up")
	if err := row.Scan(&count); err != nil {
		return nil, dto.Paging{}, err
	}

	paging := dto.Paging{
		Page:       payload.Page,
		Size:       payload.Size,
		TotalRows:  count,
		TotalPages: int(math.Ceil(float64(count) / float64(payload.Size))),
	}

	return topUps, paging, nil
}

// FindUploadedFile implements TopUp.
func (t *topUpRepository) FindUploadedFile() ([]dto.TopUp, error) {
	topUpResult, err := t.db.Query(`
		SELECT 
			id,
			maturity_time, 
			top_up_amount,
			accepted_time,
			accepted_status,
			status_information,
			transfer_confirmation_recipt,
			recipt_file
		FROM 
			top_up
		WHERE 
			transfer_confirmation_recipt = true
			AND
			accepted_status = false;`,
	)
	if err != nil {
		return nil, err
	}
	defer topUpResult.Close()
	var topUps []dto.TopUp
	for topUpResult.Next() {
		var topUp dto.TopUp
		err = topUpResult.Scan(
			&topUp.Id,
			&topUp.MaturityTime,
			&topUp.TopUpAmount,
			&topUp.AcceptedTime,
			&topUp.Accepted,
			&topUp.Status,
			&topUp.TransferConfirmRecipe,
			&topUp.File,
		)
		if err != nil {
			return nil, err
		}
		topUps = append(topUps, topUp)
	}
	return topUps, nil
}

// ConfimUpload implements TopUp.
func (t *topUpRepository) ConfimUpload(payload model.TopUp) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		UPDATE 
			top_up 
		SET 
			accepted_status = $2,
			accepted_time = $3,
			status_information = $4, 
			transfer_confirmation_recipt = $5,
			recipt_file = $6,
			maturity_time = $7
		WHERE 
			id = $1;`,
		payload.Id,
		payload.Accepted,
		payload.AcceptedTime,
		payload.Status,
		payload.TransferConfirmRecipe,
		payload.File,
		payload.MaturityTime,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	_, err = tx.Exec(
		`
		UPDATE saldo
		SET total_saving = total_saving + $1
		WHERE user_credential_id = $2;
		`,
		payload.TopUpAmount,
		payload.UserCredential.Id,
	)
	if err != nil {
		return err
	}
	tx.Commit()
	return err
}

func (t *topUpRepository) NotConfimUpload(payload model.TopUp) error {
	_, err := t.db.Exec(`
		UPDATE 
			top_up 
		SET 
			accepted_status = $2,
			status_information = $3, 
			transfer_confirmation_recipt = $4,
			recipt_file = $5,
			maturity_time = $6
		WHERE 
			id = $1;`,
		payload.Id,
		payload.Accepted,
		payload.Status,
		payload.TransferConfirmRecipe,
		payload.File,
		payload.MaturityTime,
	)
	return err
}

// Upload implements TopUp.
func (t *topUpRepository) Upload(payload model.TopUp) error {
	_, err := t.db.Exec(`
		UPDATE 
			top_up 
		SET 
			status_information = $2, 
			transfer_confirmation_recipt = $3,
			recipt_file = $4
		WHERE 
			id = $1;`,
		payload.Id,
		payload.Status,
		payload.TransferConfirmRecipe,
		payload.File,
	)
	return err
}

// FindById implements TopUp.
func (t *topUpRepository) FindById(id string) (dto.TopUpById, error) {
	rows := t.db.QueryRow(`
    SELECT 
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
        t.id, 
        t.top_up_amount, 
        t.maturity_time, 
        t.accepted_time,
        t.accepted_status,
        t.status_information,
        t.transfer_confirmation_recipt,
				t.recipt_file
    FROM 
        top_up t
    JOIN 
        biodata b ON b.user_credential_id = t.user_credential_id
    JOIN
        user_credential u ON u.id = b.user_credential_id
    WHERE 
        t.id = $1`,
		id,
	)

	var result dto.TopUpById
	result.UserCredential.Password = "********"
	err := rows.Scan(
		&result.UserCredential.Id,
		&result.UserCredential.Username,
		&result.UserCredential.Email,
		&result.UserCredential.Role,
		&result.UserCredential.IsActive,
		&result.UserCredential.VaNumber,
		&result.UserBio.NamaLengkap,
		&result.UserBio.Nik,
		&result.UserBio.NomorTelepon,
		&result.UserBio.Pekerjaan,
		&result.UserBio.TempatLahir,
		&result.UserBio.TanggalLahir,
		&result.UserBio.KodePos,
		&result.TopUp.Id,
		&result.TopUp.TopUpAmount,
		&result.TopUp.MaturityTime,
		&result.TopUp.AcceptedTime,
		&result.TopUp.Accepted,
		&result.TopUp.Status,
		&result.TopUp.TransferConfirmRecipe,
		&result.TopUp.File,
	)
	result.TopUp.VaNumber = result.UserCredential.VaNumber
	if err != nil {
		return dto.TopUpById{}, err
	}
	return result, nil
}

// FindById implements TopUp.
func (t *topUpRepository) FindByIdUser(id string) (dto.TopUpByUser, error) {
	userResult := t.db.QueryRow(`
	SELECT 
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
		b.postal_code
	FROM 
		biodata b 
	JOIN 
		user_credential u 
	ON 
		u.id = b.user_credential_id
	WHERE
	b.user_credential_id = $1;
	`, id)
	var result dto.TopUpByUser
	result.UserCredential.Password = "********"
	err := userResult.Scan(
		&result.UserCredential.Id,
		&result.UserCredential.Username,
		&result.UserCredential.Email,
		&result.UserCredential.Role,
		&result.UserCredential.IsActive,
		&result.UserCredential.VaNumber,
		&result.UserBio.NamaLengkap,
		&result.UserBio.Nik,
		&result.UserBio.NomorTelepon,
		&result.UserBio.Pekerjaan,
		&result.UserBio.TempatLahir,
		&result.UserBio.TanggalLahir,
		&result.UserBio.KodePos,
	)
	if err != nil {
		return dto.TopUpByUser{}, err
	}
	topUpResult, err := t.db.Query(`
		SELECT 
			id,
			maturity_time, 
			top_up_amount,
			accepted_time,
			accepted_status,
			status_information,
			transfer_confirmation_recipt,
			recipt_file
		FROM 
			top_up
		WHERE 
			user_credential_id = $1`,
		id,
	)
	if err != nil {
		return dto.TopUpByUser{}, err
	}
	defer topUpResult.Close()
	var topUps []dto.TopUp
	for topUpResult.Next() {
		var topUp dto.TopUp
		err = topUpResult.Scan(
			&topUp.Id,
			&topUp.MaturityTime,
			&topUp.TopUpAmount,
			&topUp.AcceptedTime,
			&topUp.Accepted,
			&topUp.Status,
			&topUp.TransferConfirmRecipe,
			&topUp.File,
		)
		topUp.VaNumber = result.UserCredential.VaNumber
		if err != nil {
			return dto.TopUpByUser{}, err
		}
		topUps = append(topUps, topUp)
	}
	result.TopUp = topUps
	return result, err
}

// Save implements TopUp.
func (t *topUpRepository) Save(payload model.TopUp) error {
	_, err := t.db.Exec("INSERT INTO top_up VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)",
		payload.Id,
		payload.UserCredential.Id,
		payload.TopUpAmount,
		payload.MaturityTime,
		payload.AcceptedTime,
		payload.Accepted,
		payload.Status,
		payload.TransferConfirmRecipe,
		payload.File,
	)
	return err
}

func NewTopUpRepository(db *sql.DB) TopUp {
	return &topUpRepository{
		db: db,
	}
}
