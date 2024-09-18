package usecase

import (
	"fmt"
	"os"
	"polen/model"
	"polen/model/dto"
	"polen/repository"
	"polen/utils/common"
	"time"
)

type TopUpUseCase interface {
	CreateNew(payload dto.TopUpUser) (dto.TopUpUser, error)
	FindByIdUser(id string) (dto.TopUpByUser, error)
	FindById(id string) (dto.TopUpById, error)
	UploadFile(payload dto.TopUpUser) (int, error)
	FindUploadedFile() ([]dto.TopUp, error)
	ConfimUploadFile(payload dto.TopUpUser) (int, error)
	Pagging(payload dto.PageRequest) ([]dto.TopUp, dto.Paging, error)
}

type topUpUseCase struct {
	repo   repository.TopUp
	userUC UserUseCase
}

// Pagging implements TopUpUseCase.
func (t *topUpUseCase) Pagging(payload dto.PageRequest) ([]dto.TopUp, dto.Paging, error) {
	// limit Size, Offset (page - 1) * size
	if payload.Page < 0 {
		payload.Page = 1
	}
	return t.repo.Pagging(payload)
}

// FindUploadedFile implements TopUpUseCase.
func (t *topUpUseCase) FindUploadedFile() ([]dto.TopUp, error) {
	return t.repo.FindUploadedFile()
}

// ConfimUploadFile implements TopUpUseCase.
func (t *topUpUseCase) ConfimUploadFile(payload dto.TopUpUser) (int, error) {
	if payload.Id == "" {
		return 400, fmt.Errorf("id is required")
	}
	data, err := t.FindById(payload.Id)
	if err != nil {
		if data.UserCredential.Id == "" {
			return 404, fmt.Errorf("data to update isnt found")
		}
		return 500, err
	}
	if payload.Status == "" {
		return 400, fmt.Errorf("information is required")
	}
	if !payload.Accepted /*false*/ {
		fmt.Println("benar")
		err = t.repo.NotConfimUpload(model.TopUp{
			Id:                    payload.Id,
			Accepted:              payload.Accepted,
			AcceptedTime:          data.TopUp.AcceptedTime,
			MaturityTime:          time.Now().Add(2 * 24 * time.Hour),
			Status:                payload.Status,
			TransferConfirmRecipe: false,
			File:                  "",
		})
		if err != nil {
			return 500, err
		}
		// delete file
		_ = os.Remove(payload.File)
		return 200, nil
	} else {
		fmt.Println("salah")
		err = t.repo.ConfimUpload(model.TopUp{
			Id:                    payload.Id,
			UserCredential:        model.UserCredential{Id: data.UserCredential.Id},
			TopUpAmount:           data.TopUp.TopUpAmount,
			Accepted:              payload.Accepted,
			AcceptedTime:          time.Now(),
			MaturityTime:          data.TopUp.MaturityTime,
			Status:                payload.Status,
			TransferConfirmRecipe: data.TopUp.TransferConfirmRecipe,
			File:                  data.TopUp.File,
		})
		if err != nil {
			return 500, err
		}
		return 200, nil
	}
}

// FindById implements TopUpUseCase.
func (t *topUpUseCase) FindById(id string) (dto.TopUpById, error) {
	return t.repo.FindById(id)
}

// FindByIdUser implements TopUpUseCase.
func (t *topUpUseCase) FindByIdUser(id string) (dto.TopUpByUser, error) {
	return t.repo.FindByIdUser(id)
}

func (t *topUpUseCase) UploadFile(payload dto.TopUpUser) (int, error) {
	data := model.TopUp{
		Id:                    payload.Id,
		File:                  payload.File,
		Status:                "waiting for acceptance",
		TransferConfirmRecipe: true,
	}
	// check is data available
	dataTopUp, err := t.FindById(data.Id)
	if err != nil {
		return 500, err
	}

	if dataTopUp.UserCredential.Id == "" {
		return 404, fmt.Errorf("data top up not found")
	}

	if dataTopUp.UserCredential.Id != payload.UserCredential.Id {
		return 403, fmt.Errorf("you are not allowed to do this")
	}

	err = t.repo.Upload(data)
	if err != nil {
		return 500, err
	}
	return 200, err
}

// CreateNew implements TopUpUseCase.
func (t *topUpUseCase) CreateNew(payload dto.TopUpUser) (dto.TopUpUser, error) {
	// cek apakah jumlah top up lebih dari 0
	if payload.TopUpAmount <= 0 {
		return dto.TopUpUser{}, fmt.Errorf("top up must be greater than zero")
	}

	uc, err := t.userUC.FindById(payload.UserCredential.Id)
	if err != nil {
		return dto.TopUpUser{}, err
	}

	modelPayload := model.TopUp{
		Id:                    common.GenerateID(),
		UserCredential:        uc,
		TopUpAmount:           payload.TopUpAmount,
		MaturityTime:          time.Now().Add(2 * 24 * time.Hour), // Menambahkan 2 hari ke waktu sekarang
		Accepted:              false,
		Status:                "waiting upoload of transfer recipe",
		TransferConfirmRecipe: false,
		File:                  "",
	}

	err = t.repo.Save(modelPayload)
	if err != nil {
		return dto.TopUpUser{}, fmt.Errorf("failed to save new topup: %v", err)
	}

	payload.UserCredential.Email = modelPayload.UserCredential.Email
	payload.UserCredential.Password = "*******"
	payload.UserCredential.Username = modelPayload.UserCredential.Username
	payload.UserCredential.VaNumber = modelPayload.UserCredential.VANumber
	payload.UserCredential.IsActive = modelPayload.UserCredential.IsActive
	payload.UserCredential.Role = modelPayload.UserCredential.Role
	payload.MaturityTime = modelPayload.MaturityTime
	payload.AcceptedTime = modelPayload.AcceptedTime
	payload.Accepted = modelPayload.Accepted
	payload.Status = modelPayload.Status
	payload.TransferConfirmRecipe = modelPayload.TransferConfirmRecipe

	return payload, nil
}

func NewTopUpUseCase(repo repository.TopUp, userUC UserUseCase) TopUpUseCase {
	return &topUpUseCase{
		repo:   repo,
		userUC: userUC,
	}
}
