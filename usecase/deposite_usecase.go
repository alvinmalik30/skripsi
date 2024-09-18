package usecase

import (
	"fmt"
	"polen/model/dto"
	"polen/repository"
	"polen/utils/common"
	"time"
)

type DepositeUseCase interface {
	CreateDeposite(payload dto.DepositeDto) (int, error)
	FindByUcId(id string) (int, dto.DepositeByUserResponse, error)
	FindById(id string) (int, dto.DepositeByIdResponse, error)
	Pagging(payload dto.PageRequest) ([]dto.Deposite, dto.Paging, error)
	Update() error
}

type depositeUseCase struct {
	repo    repository.DepositeRepository
	intRate DepositeInterestUseCase
	saldo   SaldoUsecase
}

// Update implements DepositeUseCase.
func (d *depositeUseCase) Update() error {
	return d.repo.Update()
}

// Pagging implements DepositeUseCase.
func (d *depositeUseCase) Pagging(payload dto.PageRequest) ([]dto.Deposite, dto.Paging, error) {
	// limit Size, Offset (page - 1) * size
	if payload.Page < 0 {
		payload.Page = 1
	}
	return d.repo.Pagging(payload)
}

// FindById implements DepositeUseCase.
func (d *depositeUseCase) FindById(id string) (int, dto.DepositeByIdResponse, error) {
	if id == "" {
		return 400, dto.DepositeByIdResponse{}, fmt.Errorf("id user is required")
	}
	resut, err := d.repo.FindById(id)
	if err != nil {
		return 500, dto.DepositeByIdResponse{}, err
	}
	return 200, resut, nil
}

// FindByUcId implements DepositeUseCase.
func (d *depositeUseCase) FindByUcId(id string) (int, dto.DepositeByUserResponse, error) {
	if id == "" {
		return 400, dto.DepositeByUserResponse{}, fmt.Errorf("id user is required")
	}
	resut, err := d.repo.FindByUcId(id)
	if err != nil {
		return 500, dto.DepositeByUserResponse{}, err
	}
	return 200, resut, nil
}

// CreateDeposite implements DepositeUseCase.
func (d *depositeUseCase) CreateDeposite(payload dto.DepositeDto) (int, error) {
	// check input
	if payload.InterestRate.Id == "" {
		return 400, fmt.Errorf("interest rate is required")
	}
	if payload.DepositeAmount <= 0 {
		return 400, fmt.Errorf("deposite amount must greather than zero")
	}
	// check apakah saldo cukup
	saldo, err := d.saldo.FindByIdUser(payload.UserCredential.Id)
	if err != nil {
		return 500, err
	}
	if payload.DepositeAmount > saldo.Total {
		return 400, fmt.Errorf("saldo balance isnt enought")
	}

	// generate id
	uuid := common.GenerateID()
	payload.Id = uuid

	// get data interest rate
	payload.InterestRate, err = d.intRate.FindById(payload.InterestRate.Id)
	if err != nil {
		return 500, err
	}

	// building payload
	payload.MaturityDate = time.Now().Add(time.Duration(payload.InterestRate.DurationMounth*30) * 24 * time.Hour)
	payload.Status = true
	payload.GrossProfit = (int(float64(payload.DepositeAmount)*payload.InterestRate.InterestRate) * payload.InterestRate.DurationMounth) / 365 // (setoran pokok*interest rate*durationDay)/365
	payload.Tax = int(float64(payload.GrossProfit) * payload.InterestRate.TaxRate)                                                             // TaxRate*grossProfit
	payload.NetProfit = payload.GrossProfit - payload.Tax                                                                                      // grossProfit-tax
	payload.TotalReturn = payload.DepositeAmount + payload.NetProfit                                                                           // setoran+netProfit

	err = d.repo.CreateDeposite(payload)
	if err != nil {
		return 500, err
	}
	return 200, nil
}

func NewDepositeUseCase(repo repository.DepositeRepository, intRate DepositeInterestUseCase, saldo SaldoUsecase) DepositeUseCase {
	return &depositeUseCase{
		repo:    repo,
		intRate: intRate,
		saldo:   saldo,
	}
}
