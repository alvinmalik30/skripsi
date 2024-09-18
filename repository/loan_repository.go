package repository

import (
	"database/sql"
	"fmt"
	"polen/model"
	"polen/model/dto"
	"strconv"
	"time"
)

type LoanRepository interface {
	Create(loanReq model.Loan, installenment []model.InstallenmentLoan) error
	FindById(id string) (dto.InstallenmentLoanByIdResponse, error)
	FindByLoanId(id string) (dto.Loan, error)
	Upload(payload dto.LoanInstallenmentResponse) error
	FindUploadedFile() ([]dto.LoanInstallenmentResponse, error)
	Accepted(payload dto.InstallenmentLoanByIdResponse) error
	UpdateLateFee() error
}

type loanRepository struct {
	db   *sql.DB
	user BiodataUser
}

// FindByLoanId implements LoanRepository.
func (l *loanRepository) FindByLoanId(id string) (dto.Loan, error) {
	rows, err := l.db.Query(`
	SELECT 
		id,
		is_payed,
		payment_installenment_cost,
		payment_deadline,
		total_amount_of_dept,
		late_payment_fee_nominal,
		late_payment_fee_unit,
		late_payment_fee_total,
		late_payment_fee_day,
		payment_date,
		status,
		transfer_confirmation_recipt,
		recipt_file
	FROM
		installenment_loan
	WHERE 
		loan_id = $1;
		`, id)
	if err != nil {
		return dto.Loan{}, err
	}
	defer rows.Close()
	var data []dto.LoanInstallenmentResponse
	for rows.Next() {
		var datum dto.LoanInstallenmentResponse
		var nominal float64
		var unit string
		rows.Scan(
			&datum.Id,
			&datum.IsPayed,
			&datum.PaymentInstallment,
			&datum.PaymentDeadLine,
			&datum.TotalAmountOfDepth,
			&nominal,
			&unit,
			&datum.LatePayment.LatePaymentFeesTotal,
			&datum.LatePayment.LatePaymentDays,
			&datum.PaymentDate,
			&datum.Status,
			&datum.TransferConfirmRecipe,
			&datum.File,
		)
		datum.LatePayment.LatePaymentFees = strconv.Itoa(int(nominal)) + " " + unit
		data = append(data, datum)
	}
	row := l.db.QueryRow(`
	SELECT 
		id,
		user_credential_id,
		loan_amount,
		loan_interest_rate,
		loan_interest_nominal,
		total_amount_of_dept,
		application_handling_cost_nominal,
		application_handling_cost_unit,
		loan_date_created
	FROM
		loan
	WHERE 
		id = $1;
		`, id)
	var loan dto.Loan
	var ahcn float64
	var ahcu string
	err = row.Scan(
		&loan.Id,
		&loan.UserCredentialId,
		&loan.LoanAmount,
		&loan.LoanInterestRate,
		&loan.LoanInterestNominal,
		&loan.TotalAmountOfDepth,
		&ahcn,
		&ahcu,
		&loan.LoanDateCreate,
	)
	loan.AppHandlingCost = strconv.Itoa(int(ahcn)) + " " + ahcu
	loan.Installment = data
	if err != nil {
		return dto.Loan{}, err
	}
	return loan, nil
}

// UpdateLateFee implements LoanRepository.
func (l *loanRepository) UpdateLateFee() error {
	rows, err := l.db.Query(`
	SELECT 
		id,
		total_amount_of_dept,
		payment_deadline,
		late_payment_fee_nominal,
		late_payment_fee_unit
	FROM
		installenment_loan
	WHERE 
		is_payed = false
	AND
		transfer_confirmation_recipt = true
		AND
			payment_deadline < now();
		`)
	if err != nil {
		return err
	}
	defer rows.Close()
	var data []dto.LoanInstallenmentResponse
	for rows.Next() {
		var datum dto.LoanInstallenmentResponse
		var nominal float64
		var unit string
		err := rows.Scan(
			&datum.Id,
			&datum.TotalAmountOfDepth,
			&datum.PaymentDeadLine,
			&nominal,
			&unit,
		)
		if err != nil {
			return err
		}
		// Hitung selisih antara tanggal saat ini dan tanggal target
		selisih := time.Since(datum.PaymentDeadLine)
		// Konversi selisih ke dalam hari
		hariTerlewat := int(selisih.Hours() / 24)
		datum.LatePayment.LatePaymentDays = hariTerlewat
		if unit == "rupiah" {
			datum.LatePayment.LatePaymentFeesTotal = int(nominal * float64(hariTerlewat))
		} else if unit == "percent" {
			datum.LatePayment.LatePaymentFeesTotal = int((nominal * float64(datum.TotalAmountOfDepth)) * float64(hariTerlewat))
		}
		data = append(data, datum)
		fmt.Println(datum)
	}
	fmt.Println(data)
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	for _, v := range data {
		tx.Exec(`
		UPDATE 
			installenment_loan 
		SET 
			late_payment_fee_day = $2, 
			late_payment_fee_total = $3, 
			total_amount_of_dept = total_amount_of_dept + $4
		WHERE 
			id = $1
		`,
			v.Id,
			v.LatePayment.LatePaymentDays,
			v.LatePayment.LatePaymentFeesTotal,
			v.LatePayment.LatePaymentFeesTotal,
		)
		fmt.Println(v)
	}
	tx.Commit()
	return nil
}

// Accepted implements LoanRepository.
func (t *loanRepository) Accepted(payload dto.InstallenmentLoanByIdResponse) error {
	tx, err := t.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(`
		UPDATE 
			installenment_loan 
		SET 
			is_payed = $2, 
			status = $3
		WHERE 
			id = $1;`,
		payload.LoanInst.Id,
		payload.LoanInst.IsPayed,
		payload.LoanInst.Status,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if !payload.LoanInst.IsPayed {
		_, err = tx.Exec(
			`
			UPDATE saldo
			SET total_saving = total_saving + $1
			WHERE user_credential_id = '456';
			`,
			payload.LoanInst.TotalAmountOfDepth,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return err
}

// Upload implements LoanRepository.
func (t *loanRepository) Upload(payload dto.LoanInstallenmentResponse) error {
	_, err := t.db.Exec(`
		UPDATE 
			installenment_loan 
		SET 
			payment_date = $2, 
			status = 'witing for accepment',
			transfer_confirmation_recipt = true,
			recipt_file = $3
		WHERE 
			id = $1;`,
		payload.Id,
		time.Now(),
		payload.File,
	)
	return err
}

// FindById implements LoanRepository.
func (l *loanRepository) FindById(id string) (dto.InstallenmentLoanByIdResponse, error) {
	row := l.db.QueryRow(`
	SELECT 
		installenment_loan.id, 
		user_credential_id,
		loan_id,
		installenment_loan.is_payed,
		payment_installenment_cost,
		payment_deadline,
		installenment_loan.total_amount_of_dept,
		late_payment_fee_nominal,
		late_payment_fee_unit,
		late_payment_fee_day,
		payment_date,
		installenment_loan.status,
		transfer_confirmation_recipt,
		recipt_file
	FROM 
		installenment_loan 
	JOIN
		loan on loan.id = installenment_loan.loan_id
	WHERE 
		installenment_loan.id = $1`, id)

	var loanInstallenment dto.LoanInstallenmentResponse
	var loanId string
	var ucid string
	var latePaymentFeesNominal float64
	var latepaymentfeeUnit string
	row.Scan(
		&loanInstallenment.Id,
		&ucid,
		&loanId,
		&loanInstallenment.IsPayed,
		&loanInstallenment.PaymentInstallment,
		&loanInstallenment.PaymentDeadLine,
		&loanInstallenment.TotalAmountOfDepth,
		&latePaymentFeesNominal,
		&latepaymentfeeUnit,
		&loanInstallenment.LatePayment.LatePaymentDays,
		&loanInstallenment.PaymentDate,
		&loanInstallenment.Status,
		&loanInstallenment.TransferConfirmRecipe,
		&loanInstallenment.File,
	)
	loanInstallenment.LatePayment.LatePaymentFees = strconv.Itoa(int(latePaymentFeesNominal)) + " " + latepaymentfeeUnit
	fmt.Println(loanInstallenment)
	fmt.Println(ucid)
	// get biodata user
	bio, err := l.user.FindByUcId(ucid)
	fmt.Println(bio)
	if err != nil {
		return dto.InstallenmentLoanByIdResponse{}, err
	}

	return dto.InstallenmentLoanByIdResponse{
		UserDeatail: bio,
		LoanId:      loanId,
		LoanInst:    loanInstallenment,
	}, nil
}

// FindById implements LoanRepository.
func (l *loanRepository) FindUploadedFile() ([]dto.LoanInstallenmentResponse, error) {
	rows, err := l.db.Query(`
	SELECT 
		id,
		loan_id,
		is_payed,
		payment_installenment_cost,
		payment_deadline,
		total_amount_of_dept,
		late_payment_fee_nominal,
		late_payment_fee_unit,
		late_payment_fee_day,
		payment_date,
		status,
		transfer_confirmation_recipt,
		recipt_file
	FROM
		installenment_loan
	WHERE 
		is_payed = false
	AND
		transfer_confirmation_recipt = true;
		`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var datum []dto.LoanInstallenmentResponse
	for rows.Next() {
		var loanInstallenment dto.LoanInstallenmentResponse
		var loanId string
		var latePaymentFeesNominal float64
		var latepaymentfeeUnit string
		err := rows.Scan(
			&loanInstallenment.Id,
			&loanId,
			&loanInstallenment.IsPayed,
			&loanInstallenment.PaymentInstallment,
			&loanInstallenment.PaymentDeadLine,
			&loanInstallenment.TotalAmountOfDepth,
			&latePaymentFeesNominal,
			&latepaymentfeeUnit,
			&loanInstallenment.LatePayment.LatePaymentDays,
			&loanInstallenment.PaymentDate,
			&loanInstallenment.Status,
			&loanInstallenment.TransferConfirmRecipe,
			&loanInstallenment.File,
		)
		if err != nil {
			return nil, err
		}
		loanInstallenment.LatePayment.LatePaymentFees = strconv.Itoa(int(latePaymentFeesNominal)) + " " + latepaymentfeeUnit
		datum = append(datum, loanInstallenment)
	}

	return datum, nil
}

// create implements LoanRepository.
func (l *loanRepository) Create(loan model.Loan, installenment []model.InstallenmentLoan) error {
	tx, err := l.db.Begin()
	if err != nil {
		return err
	}
	// save loan
	_, err = tx.Exec(`
	INSERT INTO
		loan
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);
	`,
		loan.Id,
		loan.UserCredentialId,
		loan.LoanAmount,
		loan.LoanDuration,
		loan.LoanInterestRate,
		loan.LoanInterestNominal,
		loan.TotalAmountOfDepth,
		loan.AppHandlingCostNominal,
		loan.AppHandlingCostUnit,
		loan.LoanDateCreate,
		loan.IsPayed,
		loan.Status,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	// kurangi sado perusahaan
	_, err = tx.Exec(
		`
		UPDATE saldo
		SET total_saving = total_saving - $1
		WHERE user_credential_id = $2;
		`,
		loan.LoanAmount,
		"456",
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	// save installentment
	for _, v := range installenment {
		_, err = tx.Exec(`
		INSERT INTO
			installenment_loan
		VALUES
			($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);
		`,
			v.Id,
			v.LoanId,
			v.IsPayed,
			v.PaymentInstallment,
			v.PaymentDeadLine,
			v.TotalAmountOfDepth,
			v.LatePaymentFeesNominal,
			v.LatePaymentFeesUnit,
			v.LatePaymentDays,
			v.LatePaymentFeesTotal,
			v.PaymentDate,
			v.Status,
			v.TransferConfirmRecipe,
			v.File,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}

func NewLoanRepository(db *sql.DB, user BiodataUser) LoanRepository {
	return &loanRepository{
		db:   db,
		user: user,
	}
}
