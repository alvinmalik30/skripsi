package usecase

import (
	"fmt"
	"polen/model"
	"polen/model/dto"
	"polen/repository"
)

type LatePaymentFeeUsecase interface {
	CreateNew(payload model.LatePaymentFee) (int, error)
	Pagging(payload dto.PageRequest) ([]model.LatePaymentFee, dto.Paging, error)
	FindById(id string) (model.LatePaymentFee, error)
	Update(payload model.LatePaymentFee) error
	DeleteById(id string) error
}

type latePaymentFeeUseCase struct {
	repo repository.LatePaymentFee
}

// CreateNew implements LatePaymentFeeUsecase.
func (a *latePaymentFeeUseCase) CreateNew(payload model.LatePaymentFee) (int, error) {
	if payload.Id == "" {
		return 400, fmt.Errorf("id is required")
	}
	if payload.Name == "" {
		return 400, fmt.Errorf("name is required")
	}
	if payload.Nominal == 0 {
		return 400, fmt.Errorf("nominal is required")
	}

	if payload.Unit == "" {
		return 400, fmt.Errorf("unit is required")
	}

	if payload.Unit != "rupiah" && payload.Unit != "percent" {
		return 400, fmt.Errorf("unit must be rupiah or percent")
	}

	if err := a.repo.CreateNew(payload); err != nil {
		return 500, err
	}
	return 201, nil
}

// DeleteById implements LatePaymentFeeUsecase.
func (a *latePaymentFeeUseCase) DeleteById(id string) error {
	app, err := a.repo.FindById(id)
	if err != nil {
		return err
	}

	err = a.repo.DeleteById(app.Id)
	if err != nil {
		return fmt.Errorf("failed to delete app handling cost: %v", err)
	}

	return nil
}

// FindById implements LatePaymentFeeUsecase.
func (a *latePaymentFeeUseCase) FindById(id string) (model.LatePaymentFee, error) {
	return a.repo.FindById(id)
}

// Pagging implements LatePaymentFeeUsecase.
func (a *latePaymentFeeUseCase) Pagging(payload dto.PageRequest) ([]model.LatePaymentFee, dto.Paging, error) {
	return a.repo.Pagging(payload)
}

// Update implements LatePaymentFeeUsecase.
func (a *latePaymentFeeUseCase) Update(payload model.LatePaymentFee) error {
	if payload.Id == "" {
		return fmt.Errorf("id is required")
	}

	result, err := a.FindById(payload.Id)
	if err != nil {
		return err
	}

	if payload.Name == "" {
		payload.Name = result.Name
	}

	if payload.Nominal == 0 {
		payload.Nominal = result.Nominal
	}
	if payload.Unit == "" {
		payload.Unit = result.Unit
	}

	if payload.Unit != "rupiah" && payload.Unit != "percent" {
		return fmt.Errorf("unit must be rupiah or percent")
	}

	err = a.repo.Update(payload)
	if err != nil {
		return fmt.Errorf("failed to update app: %v", err)
	}

	return nil
}

func NewLatePaymentFeeUseCase(repo repository.LatePaymentFee) LatePaymentFeeUsecase {
	return &latePaymentFeeUseCase{
		repo: repo,
	}
}
