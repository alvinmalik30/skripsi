package mock

import (
	"polen/model"
	"polen/model/dto"
	"time"
)

var MockUserCred = model.UserCredential{
	Id:       "1",
	Username: "alvin ",
	Email:    "alvinmalik@gmail.com",
	Password: "123",
	VANumber: "Efvfdvfdhucsucuh",
	Role:     "borrower",
	IsActive: true,
}
var MockAuthLoginReq = dto.AuthLoginRequest{
	Username: "alvinmalik",
	Password: "password",
}
var MockAuthReq = dto.AuthRequest{
	Username: "alvinmalik",
	Email:    "alvin @gmail.com",
	Password: "password",
	Role:     "borrower",
}
var MockAuthResponse = dto.AuthResponse{
	Username: MockAuthReq.Username,
	Token:    "",
}
var MockBiodata = model.BiodataUser{
	Id:             "1",
	UserCredential: model.UserCredential{Id: MockUserCred.Id},
	NamaLengkap:    "alvin malik",
	Nik:            "32010",
	NomorTelepon:   "081287743960",
	Pekerjaan:      "IT",
	TempatLahir:    "jakarta",
	TanggalLahir:   time.Date(2000, time.December, 12, 0, 0, 0, 0, time.UTC),
	KodePos:        "1610",
	IsAglible:      false,
	StatusUpdate:   false,
	Information:    "biodata is not updated",
}
var MockSaldoData = dto.Saldo{
	Id:    "1",
	UcId:  "1",
	Total: 100000,
}
var MockSaldoDatas = []dto.Saldo{
	{
		Id:    "1",
		UcId:  "1",
		Total: 100000,
	},
}
var MockSaldo = model.Saldo{
	Id:    "1",
	UcId:  "1",
	Total: 0,
}
var MockPageReq = dto.PageRequest{
	Page: 1,
	Size: 5,
}
var MockPaging = dto.Paging{
	Page:       1,
	Size:       5,
	TotalRows:  1,
	TotalPages: 1,
}
var MockDepositeDto = []dto.Deposite{
	{
		Id:             "1",
		DepositeAmount: 100,
		InterestRate:   1,
		TaxRate:        1,
		DurationMounth: 12,
		CreateDate:     time.Now(),
		MaturityDate:   time.Now(),
		Status:         "pending",
		GrossProfit:    12000,
		Tax:            10,
		NetProfit:      2000,
		TotalReturn:    10000,
	},
}
var MockTopUpByUser = dto.TopUpByUser{
	UserCredential: dto.GetAuthResponse{
		Id:       "1",
		Username: "alvin ",
		Email:    "alvinmalik@gmail.com",
		Password: "123",
		VaNumber: "Efvfdvfdhucsucuh",
		Role:     "borrower",
		IsActive: true,
	},
	UserBio: dto.BiodataRequest{
		NamaLengkap:  "alvin malik",
		Nik:          "32010",
		NomorTelepon: "081287743960",
		Pekerjaan:    "IT",
		TempatLahir:  "Jakarta",
		TanggalLahir: "2000-12-12",
		KodePos:      "1610",
	},
	TopUp: []dto.TopUp{
		{
			Id:                    "1",
			VaNumber:              "Efvfdvfdhucsucuh",
			TopUpAmount:           100,
			MaturityTime:          time.Now(),
			AcceptedTime:          time.Now(),
			Accepted:              false,
			Status:                "pending",
			TransferConfirmRecipe: false,
			File:                  "",
		},
	},
}
var MockTopUpId = dto.TopUpById{
	UserCredential: dto.GetAuthResponse{
		Id:       "1",
		Username: "alvin ",
		Email:    "alvinmalik@gmail.com",
		Password: "123",
		VaNumber: "Efvfdvfdhucsucuh",
		Role:     "borrower",
		IsActive: true,
	},
	UserBio: dto.BiodataRequest{
		NamaLengkap:  "alvin malik",
		Nik:          "32010",
		NomorTelepon: "081287743960",
		Pekerjaan:    "IT",
		TempatLahir:  "Jakarta",
		TanggalLahir: "2000-12-12",
		KodePos:      "1610",
	},
	TopUp: dto.TopUp{
		Id:                    "1",
		VaNumber:              "Efvfdvfdhucsucuh",
		TopUpAmount:           100,
		MaturityTime:          time.Now(),
		AcceptedTime:          time.Now(),
		Accepted:              false,
		Status:                "pending",
		TransferConfirmRecipe: false,
		File:                  "",
	},
}
var MockTopUp = model.TopUp{
	Id: "1",
	UserCredential: model.UserCredential{
		Id:       "1",
		Username: "alvin ",
		Email:    "alvinmalik@gmail.com",
		Password: "123",
		VANumber: "Efvfdvfdhucsucuh",
		Role:     "borrower",
		IsActive: true,
	},
	TopUpAmount:           100,
	MaturityTime:          time.Now(),
	AcceptedTime:          time.Now(),
	Accepted:              false,
	Status:                "pending",
	TransferConfirmRecipe: false,
	File:                  "",
}
var MockListTopUp = []dto.TopUp{
	{
		Id:                    "1",
		VaNumber:              "vfdbvhfdbdhf",
		TopUpAmount:           1000,
		MaturityTime:          time.Now(),
		AcceptedTime:          time.Now(),
		Accepted:              false,
		Status:                "not accepted",
		TransferConfirmRecipe: false,
		File:                  "",
	},
}
var MockDeposites = []dto.DepositeInterestRequest{
	{
		Id:             "1",
		InterestRate:   1,
		TaxRate:        1,
		DurationMounth: 12,
	},
}
var MockUserCreds = []model.UserCredential{
	{
		Id:       "1",
		Username: "alvin ",
		Email:    "alvinmalik@gmail.com",
		Password: "123",
		VANumber: "Efvfdvfdhucsucuh",
		Role:     "admin",
		IsActive: true,
	},
}
var MockUpdateBioReq = dto.UpdateBioRequest{
	UserCredentialId: "1",
	IsAglible:        false,
	Information:      "pending",
}
var MockBiodataResponse = dto.BiodataResponse{
	Id:          "1",
	NamaLengkap: "alvin malik",
	UserCredential: dto.GetAuthResponse{
		Id:       "1",
		Username: "alvinmalik",
		Email:    "alvin @gmail.com",
		Password: "123",
		Role:     "peminjam",
		VaNumber: "bfdffbfvfhvf",
		IsActive: false,
	},
	Nik:          "32010",
	NomorTelepon: "081287743960",
	Pekerjaan:    "IT",
	TempatLahir:  "Jakarta",
	TanggalLahir: "2000-12-12",
	KodePos:      "1610",
	IsAglible:    false,
	StatusUpdate: false,
	Information:  "Additional",
}
var MockBiodataResponses = []dto.BiodataResponse{
	{
		Id:          "1",
		NamaLengkap: "alvin malik",
		UserCredential: dto.GetAuthResponse{
			Id:       "1",
			Username: "alvinmalik",
			Email:    "alvin @gmail.com",
			Password: "123",
			Role:     "peminjam",
			VaNumber: "bfdffbfvfhvf",
			IsActive: false,
		},
		Nik:          "32010",
		NomorTelepon: "081287743960",
		Pekerjaan:    "IT",
		TempatLahir:  "Jakarta",
		TanggalLahir: "2000-12-12",
		KodePos:      "1610",
		IsAglible:    false,
		StatusUpdate: false,
		Information:  "Additional",
	},
}
var MockDepositeInterest = model.DepositeInterest{
	Id:             "1",
	CreateDate:     time.Now(),
	InterestRate:   1,
	TaxRate:        1,
	DurationMounth: 12,
}
var MockDepositeInterestReq = dto.DepositeInterestRequest{
	Id:             "1",
	InterestRate:   1,
	TaxRate:        1,
	DurationMounth: 12,
}
var MockDepositeByIdResponse = dto.DepositeByIdResponse{
	BioUser: dto.BiodataResponse{
		Id:          "1",
		NamaLengkap: "alvin malik",
		UserCredential: dto.GetAuthResponse{
			Id:       "1",
			Username: "alvinmalik",
			Email:    "alvin @gmail.com",
			Password: "123",
			Role:     "peminjam",
			VaNumber: "bfdffbfvfhvf",
			IsActive: false,
		},
		Nik:          "32010",
		NomorTelepon: "081287743960",
		Pekerjaan:    "IT",
		TempatLahir:  "Jakarta",
		TanggalLahir: "2000-12-12",
		KodePos:      "1610",
		IsAglible:    false,
		StatusUpdate: false,
		Information:  "Additional",
	},
	Deposite: dto.Deposite{
		Id:             "1",
		DepositeAmount: 100,
		InterestRate:   1,
		TaxRate:        1,
		DurationMounth: 12,
		CreateDate:     time.Now(),
		MaturityDate:   time.Now(),
		Status:         "pending",
		GrossProfit:    12000,
		Tax:            10,
		NetProfit:      2000,
		TotalReturn:    10000,
	},
}
var MockAppHCDatas = []model.AppHandlingCost{
	{
		Id:      "1",
		Name:    "alvin malik",
		Nominal: 10000,
		Unit:    "Rp",
	},
}
var MockAppHC = model.AppHandlingCost{
	Id:      "1",
	Name:    "alvin malik",
	Nominal: 10000,
	Unit:    "rupiah",
}
var MockLatePFDatas = []model.LatePaymentFee{
	{
		Id:      "1",
		Name:    "alvin malik",
		Nominal: 10000,
		Unit:    "Rp",
	},
}
var MockLatePF = model.LatePaymentFee{
	Id:      "1",
	Name:    "alvin malik",
	Nominal: 10000,
	Unit:    "rupiah",
}
var MockLoanInterest = model.LoanInterest{
	Id:               "1",
	DurationMonths:   12,
	LoanInterestRate: 10,
}
var MockLoanInterestDatas = []model.LoanInterest{
	{
		Id:               "1",
		DurationMonths:   12,
		LoanInterestRate: 10,
	},
}
var MockLoan = model.Loan{
	Id:                     "1",
	UserCredentialId:       "1",
	LoanDateCreate:         time.Now(),
	LoanAmount:             1,
	LoanDuration:           12,
	LoanInterestRate:       10,
	LoanInterestNominal:    10000,
	AppHandlingCostNominal: 100,
	AppHandlingCostUnit:    "Rp",
	TotalAmountOfDepth:     100,
	IsPayed:                false,
	Status:                 "pending",
}
var MockInstallLoan = model.InstallenmentLoan{

	Id:                     "1",
	LoanId:                 "1",
	IsPayed:                false,
	PaymentInstallment:     10000,
	PaymentDeadLine:        time.Now(),
	TotalAmountOfDepth:     100,
	LatePaymentFeesNominal: 1000,
	LatePaymentFeesUnit:    "Rp",
	LatePaymentDays:        120,
	LatePaymentFeesTotal:   100000,
	PaymentDate:            time.Now(),
	Status:                 "pending",
	TransferConfirmRecipe:  false,
	File:                   "",
}
var MockInstallLoanDatas = []model.InstallenmentLoan{
	{
		Id:                     "1",
		LoanId:                 "1",
		IsPayed:                false,
		PaymentInstallment:     10000,
		PaymentDeadLine:        time.Now(),
		TotalAmountOfDepth:     100,
		LatePaymentFeesNominal: 1000,
		LatePaymentFeesUnit:    "Rp",
		LatePaymentDays:        120,
		LatePaymentFeesTotal:   100000,
		PaymentDate:            time.Now(),
		Status:                 "pending",
		TransferConfirmRecipe:  false,
		File:                   "",
	},
}
var MockLoanInstallRespons = []dto.LoanInstallenmentResponse{
	{
		Id:                 "1",
		IsPayed:            false,
		PaymentInstallment: 100,
		PaymentDeadLine:    time.Now(),
		TotalAmountOfDepth: 100,
		LatePayment: dto.LatePayment{
			LatePaymentFees:      "1000",
			LatePaymentDays:      30,
			LatePaymentFeesTotal: 300000,
		},
		PaymentDate:           time.Now(),
		Status:                "pending",
		TransferConfirmRecipe: false,
		File:                  "",
	},
}
var MockInstallLoanByIdResp = dto.InstallenmentLoanByIdResponse{
	UserDeatail: dto.BiodataResponse{
		Id:          "1",
		NamaLengkap: "alvin malik",
		UserCredential: dto.GetAuthResponse{
			Id:       "1",
			Username: "alvinmalik",
			Email:    "alvin @gmail.com",
			Password: "123",
			Role:     "peminjam",
			VaNumber: "bfdffbfvfhvf",
			IsActive: false,
		},
		Nik:          "32010",
		NomorTelepon: "081287743960",
		Pekerjaan:    "IT",
		TempatLahir:  "Jakarta",
		TanggalLahir: "2000-12-12",
		KodePos:      "1610",
		IsAglible:    false,
		StatusUpdate: false,
		Information:  "Additional",
	},
	LoanId: "1",
	LoanInst: dto.LoanInstallenmentResponse{
		Id:                 "1",
		IsPayed:            false,
		PaymentInstallment: 100,
		PaymentDeadLine:    time.Now(),
		TotalAmountOfDepth: 100,
		LatePayment: dto.LatePayment{
			LatePaymentFees:      "1000",
			LatePaymentDays:      30,
			LatePaymentFeesTotal: 300000,
		},
		PaymentDate:           time.Now(),
		Status:                "pending",
		TransferConfirmRecipe: false,
		File:                  "",
	},
}
var MockDeposite = dto.DepositeDto{
	Id: "1",
	UserCredential: dto.GetAuthResponse{
		Id:       "1",
		Username: "alvinmalik",
		Email:    "alvin @gmail.com",
		Password: "123",
		Role:     "peminjam",
		VaNumber: "bfdffbfvfhvf",
		IsActive: false,
	},
	InterestRate: dto.DepositeInterestRequest{
		Id:             "1",
		InterestRate:   5,
		TaxRate:        10,
		DurationMounth: 12,
	},
	DepositeAmount: 100000,
	MaturityDate:   time.Now(),
	Status:         false,
	GrossProfit:    100,
	Tax:            10,
	NetProfit:      10,
	TotalReturn:    110,
}
var MockDepositesDTO = []dto.DepositeDto{
	{
		Id: "1",
		UserCredential: dto.GetAuthResponse{
			Id:       "1",
			Username: "alvinmalik",
			Email:    "alvin @gmail.com",
			Password: "123",
			Role:     "peminjam",
			VaNumber: "bfdffbfvfhvf",
			IsActive: false,
		},
		InterestRate: dto.DepositeInterestRequest{
			Id:             "1",
			InterestRate:   5,
			TaxRate:        10,
			DurationMounth: 12,
		},
		DepositeAmount: 100000,
		MaturityDate:   time.Now(),
		Status:         false,
		GrossProfit:    100,
		Tax:            10,
		NetProfit:      10,
		TotalReturn:    110,
	},
}
var MockDepositeByUserResponse = dto.DepositeByUserResponse{
	BioUser: dto.BiodataResponse{
		Id:          "1",
		NamaLengkap: "alvin malik",
		UserCredential: dto.GetAuthResponse{
			Id:       "1",
			Username: "alvinmalik",
			Email:    "alvin @gmail.com",
			Password: "123",
			Role:     "peminjam",
			VaNumber: "bfdffbfvfhvf",
			IsActive: false,
		},
		Nik:          "32010",
		NomorTelepon: "081287743960",
		Pekerjaan:    "IT",
		TempatLahir:  "Jakarta",
		TanggalLahir: "2000-12-12",
		KodePos:      "1610",
		IsAglible:    false,
		StatusUpdate: false,
		Information:  "Additional",
	},
	Deposite: MockDepositeDto,
}
var MockDTOLoan = dto.Loan{
	Id:                  "1",
	UserCredentialId:    "1",
	LoanDateCreate:      time.Now(),
	LoanAmount:          0,
	LoanDuration:        0,
	LoanInterestRate:    0,
	LoanInterestNominal: 0,
	AppHandlingCost:     "",
	TotalAmountOfDepth:  0,
	IsPayed:             false,
	Status:              "pending",
	Installment: []dto.LoanInstallenmentResponse{
		{
			Id:                 "1",
			IsPayed:            false,
			PaymentInstallment: 0,
			PaymentDeadLine:    time.Now(),
			TotalAmountOfDepth: 0,
			LatePayment: dto.LatePayment{
				LatePaymentFees:      "rupiah",
				LatePaymentDays:      0,
				LatePaymentFeesTotal: 0,
			},
			PaymentDate:           time.Now(),
			Status:                "pending",
			TransferConfirmRecipe: false,
			File:                  "",
		},
	},
}
