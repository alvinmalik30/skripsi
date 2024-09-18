package usecase

import (
	"fmt"
	"polen/model"
	"polen/model/dto"
	"polen/repository"
	"polen/utils/common"
	"time"

	"github.com/gin-gonic/gin"
)

type BiodataUserUseCase interface {
	FindByUserCredential(ctx *gin.Context) (dto.BiodataResponse, error)
	FindByUcId(id string) (dto.BiodataResponse, error)
	FindUserUpdated() ([]dto.BiodataResponse, error)
	UserUpdate(payload dto.BiodataRequest, ctx *gin.Context) (int, error)
	AdminUpdate(payload dto.UpdateBioRequest, ctx *gin.Context) (int, error)
	Paging(payload dto.PageRequest) ([]dto.BiodataResponse, dto.Paging, error)
}

type biodataUserUseCase struct {
	repo   repository.BiodataUser
	userUC UserUseCase
	ctx    *gin.Context
}

// Paging implements BiodataUserUseCase.
func (bio *biodataUserUseCase) Paging(payload dto.PageRequest) ([]dto.BiodataResponse, dto.Paging, error) {
	// limit Size, Offset (page - 1) * size
	if payload.Page < 0 {
		payload.Page = 1
	}
	return bio.repo.Pagging(payload)
}

// FindByUcId implements BiodataUserUseCase.
func (bio *biodataUserUseCase) FindByUcId(id string) (dto.BiodataResponse, error) {
	return bio.repo.FindByUcId(id)
}

// AdminUpdate implements BiodataUserUseCase.
func (bio *biodataUserUseCase) AdminUpdate(payload dto.UpdateBioRequest, ctx *gin.Context) (int, error) {
	if payload.Information == "" {
		return 400, fmt.Errorf("information detail is required")
	}
	biodata := model.BiodataUser{
		IsAglible:   payload.IsAglible,
		Information: payload.Information,
	}
	result, err := bio.FindByUcId(payload.UserCredentialId)
	if err != nil {
		return 500, err
	}
	biodata.Id = result.Id
	if biodata.IsAglible {
		biodata.StatusUpdate = true
	}
	err = bio.repo.AdminUpdate(biodata)
	if err != nil {
		return 500, err
	}
	return 200, nil
}

// UserUpdate implements BiodataUserUseCase.
func (bio *biodataUserUseCase) UserUpdate(payload dto.BiodataRequest, ctx *gin.Context) (int, error) {
	if payload.NamaLengkap == "" {
		return 400, fmt.Errorf("name is required")
	}
	if payload.Nik == "" {
		return 400, fmt.Errorf("nik is required")
	}
	if payload.Nik == "" {
		return 400, fmt.Errorf("nik is required")
	}
	if payload.NomorTelepon == "" {
		return 400, fmt.Errorf("phone is required")
	}
	if payload.Pekerjaan == "" {
		return 400, fmt.Errorf("job is required")
	}
	if payload.TanggalLahir == "" {
		return 400, fmt.Errorf("birth date is required")
	}
	if payload.TempatLahir == "" {
		return 400, fmt.Errorf("place of birth is required")
	}
	if payload.KodePos == "" {
		return 400, fmt.Errorf("postal code is required")
	}
	// find data by user credential
	result, err := bio.FindByUserCredential(ctx)
	if err != nil {
		return 500, err
	}
	// made payload
	layout := "2006-01-02"
	parsedTime, err := time.Parse(layout, payload.TanggalLahir)
	if err != nil {
		return 500, err
	}
	ucId, err := common.GetId(ctx)
	if err != nil {
		return 500, err
	}
	biodata := model.BiodataUser{
		Id:             result.Id,
		UserCredential: model.UserCredential{Id: ucId},
		NamaLengkap:    payload.NamaLengkap,
		Nik:            payload.Nik,
		NomorTelepon:   payload.NomorTelepon,
		Pekerjaan:      payload.Pekerjaan,
		TempatLahir:    payload.TempatLahir,
		TanggalLahir:   parsedTime,
		KodePos:        payload.KodePos,
		StatusUpdate:   true,
		Information:    "waiting for acceptence",
	}
	if result.IsAglible {
		biodata.IsAglible = false
	}
	err = bio.repo.UserUpdate(biodata)
	if err != nil {
		return 500, err
	}
	return 200, nil
}

// FindUserUpdated implements BiodataUserUseCase.
func (bio *biodataUserUseCase) FindUserUpdated() ([]dto.BiodataResponse, error) {
	return bio.repo.FindUserUpdated()
}

// FindByUserCredential implements BiodataUserUseCase.
func (bio *biodataUserUseCase) FindByUserCredential(ctx *gin.Context) (dto.BiodataResponse, error) {
	id, err := common.GetId(ctx)
	if err != nil {
		return dto.BiodataResponse{}, err
	}
	return bio.repo.FindByUcId(id)
}

func NewBiodataUserUseCase(repo repository.BiodataUser, userUC UserUseCase, ctx *gin.Context) BiodataUserUseCase {
	return &biodataUserUseCase{
		repo:   repo,
		userUC: userUC,
		ctx:    ctx,
	}
}
