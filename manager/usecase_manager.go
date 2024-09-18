package manager

import (
	"polen/usecase"

	"github.com/gin-gonic/gin"
)

type UseCaseManager interface {
	AuthUseCase() usecase.AuthUseCase
	UserUseCase() usecase.UserUseCase
	BiodataUserUseCase() usecase.BiodataUserUseCase
	TopUpUsecase() usecase.TopUpUseCase
	DepositerInterestUseCase() usecase.DepositeInterestUseCase
	LoanInterestUseCase() usecase.LoanInterestUseCase
	SaldoUsecase() usecase.SaldoUsecase
	DepositeUsecase() usecase.DepositeUseCase
	AppHandlingCostUseCase() usecase.AppHandlingCostUsecase
	LatePaymentFee() usecase.LatePaymentFeeUsecase
	LoanUsecase() usecase.LoanUseCase
}

type useCaseManager struct {
	repoManager RepoManager
	ctx         *gin.Context
}

// LatePaymentFee implements UseCaseManager.
func (u *useCaseManager) LatePaymentFee() usecase.LatePaymentFeeUsecase {
	return usecase.NewLatePaymentFeeUseCase(u.repoManager.LatePaymentFee())
}

// LoanUsecase implements UseCaseManager.
func (u *useCaseManager) LoanUsecase() usecase.LoanUseCase {
	return usecase.NewLoanUseCase(u.repoManager.LoanRepo(), u.LoanInterestUseCase(), u.AppHandlingCostUseCase(), u.LatePaymentFee())
}

// DepositrUsecase implements UseCaseManager.
func (u *useCaseManager) DepositeUsecase() usecase.DepositeUseCase {
	return usecase.NewDepositeUseCase(u.repoManager.DepositeRepo(), u.DepositerInterestUseCase(), u.SaldoUsecase())
}

// AppHandlingCostUseCase implements UseCaseManager.
func (u *useCaseManager) AppHandlingCostUseCase() usecase.AppHandlingCostUsecase {
	return usecase.NewAppHandlingCostUseCase(u.repoManager.AppHandlingCostRepo())
}

// LoanInterestUseCase implements UseCaseManager.
func (u *useCaseManager) LoanInterestUseCase() usecase.LoanInterestUseCase {
	return usecase.NewLoanInterestUseCase(u.repoManager.LoanInterestRepo())

}

// SaldoUsecase implements UseCaseManager.
func (u *useCaseManager) SaldoUsecase() usecase.SaldoUsecase {
	return usecase.NewSaldoUsecase(u.repoManager.SaldoRepo())
}

// DepositerInterestUseCase implements UseCaseManager.
func (u *useCaseManager) DepositerInterestUseCase() usecase.DepositeInterestUseCase {
	return usecase.NewDepositeInterestUseCase(u.repoManager.DepositeInterestRepo())
}

// TopUpUsecase implements UseCaseManager.
func (u *useCaseManager) TopUpUsecase() usecase.TopUpUseCase {
	return usecase.NewTopUpUseCase(u.repoManager.TopUpRepo(), u.UserUseCase())
}

// BiodataUserUseCase implements UseCaseManager.
func (u *useCaseManager) BiodataUserUseCase() usecase.BiodataUserUseCase {
	return usecase.NewBiodataUserUseCase(u.repoManager.BiodataRepo(), u.UserUseCase(), u.ctx)
}

// AuthUseCase implements UseCaseManager.
func (u *useCaseManager) AuthUseCase() usecase.AuthUseCase {
	return usecase.NewAuthUseCase(u.repoManager.UserRepo())
}

// UserUseCase implements UseCaseManager.
func (u *useCaseManager) UserUseCase() usecase.UserUseCase {
	return usecase.NewUserUseCase(u.repoManager.UserRepo(), u.ctx)
}

func NewUsecaseManager(repoManager RepoManager) UseCaseManager {
	return &useCaseManager{
		repoManager: repoManager,
		ctx:         &gin.Context{},
	}
}
