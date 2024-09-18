package usecase

import (
	"fmt"
	"polen/model"
	"polen/model/dto"
	"polen/repository"
)

type AppHandlingCostUsecase interface {
	CreateNew(payload model.AppHandlingCost) (int, error)
	Pagging(payload dto.PageRequest) ([]model.AppHandlingCost, dto.Paging, error)
	FindById(id string) (model.AppHandlingCost, error)
	Update(payload model.AppHandlingCost) error
	DeleteById(id string) error
}

type appHandlingCostUseCase struct {
	repo repository.AppHandlingCost
}

// CreateNew implements AppHandlingCostUsecase.
func (a *appHandlingCostUseCase) CreateNew(payload model.AppHandlingCost) (int, error) {
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

// DeleteById implements AppHandlingCostUsecase.
func (a *appHandlingCostUseCase) DeleteById(id string) error {
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

// FindById implements AppHandlingCostUsecase.
func (a *appHandlingCostUseCase) FindById(id string) (model.AppHandlingCost, error) {
	return a.repo.FindById(id)
}

// Pagging implements AppHandlingCostUsecase.
func (a *appHandlingCostUseCase) Pagging(payload dto.PageRequest) ([]model.AppHandlingCost, dto.Paging, error) {
	return a.repo.Pagging(payload)
}

// Update implements AppHandlingCostUsecase.
func (a *appHandlingCostUseCase) Update(payload model.AppHandlingCost) error {
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

func NewAppHandlingCostUseCase(repo repository.AppHandlingCost) AppHandlingCostUsecase {
	return &appHandlingCostUseCase{
		repo: repo,
	}
}
