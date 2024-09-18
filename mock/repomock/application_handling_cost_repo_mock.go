package repomock

import (
	"polen/model"
	"polen/model/dto"

	"github.com/stretchr/testify/mock"
)

type ApplicationHCRepoMock struct {
	mock.Mock
}

// CreateNew implements AppHandlingCost.
func (a *ApplicationHCRepoMock) CreateNew(payload model.AppHandlingCost) error {
	return a.Called(payload).Error(0)
}

// DeleteById implements AppHandlingCost.
func (a *ApplicationHCRepoMock) DeleteById(id string) error {
	return a.Called(id).Error(0)
}

// FindById implements AppHandlingCost.
func (a *ApplicationHCRepoMock) FindById(id string) (model.AppHandlingCost, error) {
	a2 := a.Called(id)
	if a2.Get(1) != nil {
		return model.AppHandlingCost{}, a2.Error(1)
	}
	return a2.Get(0).(model.AppHandlingCost), nil
}

// Pagging implements AppHandlingCost.
func (a *ApplicationHCRepoMock) Pagging(payload dto.PageRequest) ([]model.AppHandlingCost, dto.Paging, error) {
	a2 := a.Called(payload)
	if a2.Get(2) != nil {
		return nil, dto.Paging{}, a2.Error(2)
	}
	return a2.Get(0).([]model.AppHandlingCost), a2.Get(1).(dto.Paging), nil
}

// Update implements AppHandlingCost.
func (a *ApplicationHCRepoMock) Update(payload model.AppHandlingCost) error {
	return a.Called(payload).Error(0)
}
