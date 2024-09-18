package usecase

import (
	"fmt"
	"polen/model"
	"polen/model/dto"
	"polen/repository"
	"polen/utils/common"
	"polen/utils/security"
	"regexp"

	"github.com/gin-gonic/gin"
)

type UserUseCase interface {
	FindByUsername(username string, ctx *gin.Context) (model.UserCredential, error)
	Register(payload dto.AuthRequest) error
	Paging(payload dto.PageRequest, ctx *gin.Context) ([]model.UserCredential, dto.Paging, error)
	FindById(id string) (model.UserCredential, error)
}

type userUseCase struct {
	repo repository.UserRepository
	ctx  *gin.Context
}

// Paging implements UserUseCase.
func (u *userUseCase) Paging(payload dto.PageRequest, ctx *gin.Context) ([]model.UserCredential, dto.Paging, error) {
	// limit Size, Offset (page - 1) * size
	if payload.Page < 0 {
		payload.Page = 1
	}
	role, err := common.GetRole(ctx)
	if err != nil {
		return nil, dto.Paging{}, err
	}
	if role != "admin" {
		return nil, dto.Paging{}, fmt.Errorf("you are not allowed")
	}
	return u.repo.Pagging(payload)
}

// FindById implements UserUseCase.
func (u *userUseCase) FindById(id string) (model.UserCredential, error) {
	return u.repo.FindById(id)
}

// FindByUsername implements UserUseCase.
func (u *userUseCase) FindByUsername(username string, ctx *gin.Context) (model.UserCredential, error) {
	role, err := common.GetRole(ctx)
	if err != nil {
		return model.UserCredential{}, err
	}
	name, err := common.GetName(ctx)
	if err != nil {
		return model.UserCredential{}, err
	}
	if role != "admin" {
		if name != username {
			return model.UserCredential{}, fmt.Errorf("you are not allowed")
		}
	}
	return u.repo.FindByUsername(username)
}

// Register implements UserUseCase.
func (u *userUseCase) Register(payload dto.AuthRequest) error {
	if payload.Username == "" {
		return fmt.Errorf("username required")
	}
	if payload.Password == "" {
		return fmt.Errorf("password required")
	}
	if payload.Email == "" {
		return fmt.Errorf("email required")
	}
	if payload.Role == "" {
		return fmt.Errorf("role is required")
	}
	if payload.Role != "peminjam" && payload.Role != "pemodal" && payload.Role != "admin" {
		return fmt.Errorf("role you has choose isnt available")
	}
	// Pola regex untuk memeriksa format email
	emailPattern := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

	// Mencocokkan alamat email dengan pola regex
	isValid := isValidEmail(emailPattern, payload.Email)
	if !isValid {
		return fmt.Errorf("is not valid email")
	}

	hashPassword, err := security.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	userCredential := model.UserCredential{
		Id:       common.GenerateID(),
		Username: payload.Username,
		Email:    payload.Email,
		Password: hashPassword,
		VANumber: common.GenerateID(),
		Role:     payload.Role,
	}

	biodataId := common.GenerateID()

	if userCredential.Role == "pemodal" {
		saldoId := common.GenerateID()
		err = u.repo.Saldo(userCredential, saldoId, biodataId)
	} else {
		err = u.repo.Save(userCredential, biodataId)
	}
	if err != nil {
		return fmt.Errorf("failed save user: %v", err.Error())
	}
	return nil
}

func isValidEmail(pattern, email string) bool {
	match, _ := regexp.MatchString(pattern, email)
	return match
}

func NewUserUseCase(repo repository.UserRepository, ctx *gin.Context) UserUseCase {
	return &userUseCase{
		repo: repo,
		ctx:  ctx,
	}
}
