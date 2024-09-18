package usecase

import (
	"fmt"
	"polen/model"
	"polen/model/dto"
	"polen/repository"
	"polen/utils/common"
	"time"
)

type DepositeInterestUseCase interface {
	CreateNew(payload dto.DepositeInterestRequest) (int, error)
	Pagging(payload dto.PageRequest) ([]dto.DepositeInterestRequest, dto.Paging, error)
	FindById(id string) (dto.DepositeInterestRequest, error)
	Update(payload dto.DepositeInterestRequest) error
	DeleteById(id string) error
}

type depositeInterestUseCase struct {
	repo repository.DepositeInterest
}

// Pagging implements DepositeInterestUseCase.
func (d *depositeInterestUseCase) Pagging(payload dto.PageRequest) ([]dto.DepositeInterestRequest, dto.Paging, error) {
	// limit Size, Offset (page - 1) * size
	if payload.Page < 0 {
		payload.Page = 1
	}
	return d.repo.Pagging(payload)
}

// DeleteById implements DepositeInterestUseCase.
func (d *depositeInterestUseCase) DeleteById(id string) error {
	deposite, err := d.repo.FindById(id)
	if err != nil {
		return err
	}

	err = d.repo.DeleteById(deposite.Id)
	if err != nil {
		return fmt.Errorf("failed to delete deposite: %v", err)
	}

	return nil
}

// FindById implements DepositeInterestUseCase.
func (d *depositeInterestUseCase) FindById(id string) (dto.DepositeInterestRequest, error) {
	return d.repo.FindById(id)
}

// Update implements DepositeInterestUseCase.
func (d *depositeInterestUseCase) Update(payload dto.DepositeInterestRequest) error {
	if payload.Id == "" {
		return fmt.Errorf("id is required")
	}

	result, err := d.FindById(payload.Id)
	if err != nil {
		return err
	}

	if payload.InterestRate == 0 {
		payload.InterestRate = result.InterestRate
	}
	if payload.TaxRate == 0 {
		payload.TaxRate = result.TaxRate
	}
	if payload.DurationMounth == 0 {
		payload.DurationMounth = result.DurationMounth
	}

	err = d.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update deposite: %v", err)
	}

	return nil
}

// CreateNew implements DepositeInteresetUseCase.
func (d *depositeInterestUseCase) CreateNew(payload dto.DepositeInterestRequest) (int, error) {
	if payload.Id == "" {
		return 400, fmt.Errorf("id is required")
	}
	if payload.InterestRate == 0 {
		return 400, fmt.Errorf("interest rate is required")
	}
	if payload.TaxRate == 0 {
		return 400, fmt.Errorf("tax is required")
	}
	if payload.DurationMounth == 0 {
		return 400, fmt.Errorf("duration month is required")
	}
	model := model.DepositeInterest{
		Id:             common.GenerateID(),
		CreateDate:     time.Now(),
		InterestRate:   payload.InterestRate,
		TaxRate:        payload.TaxRate,
		DurationMounth: payload.DurationMounth,
	}

	if err := d.repo.Save(model); err != nil {
		return 500, err
	}
	return 201, nil
}

func NewDepositeInterestUseCase(repo repository.DepositeInterest) DepositeInterestUseCase {
	return &depositeInterestUseCase{
		repo: repo,
	}
}
