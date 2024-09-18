package dto

import "time"

type TopUpUser struct {
	Id                    string          `json:"id"`
	UserCredential        GetAuthResponse `json:"user credential"`
	TopUpAmount           int             `json:"top up amount"`
	MaturityTime          time.Time       `json:"maturity time"`
	AcceptedTime          time.Time       `json:"accepted time"`
	Accepted              bool            `json:"accepted status"`
	Status                string          `json:"information status"`
	TransferConfirmRecipe bool            `json:"confirm recipe send"`
	File                  string
}

type TopUpReq struct {
	TopUpAmount int `json:"top up amount"`
}

type TopUpupdate struct {
	Id       string `json:"id"`
	Accepted bool   `json:"accepted status"`
	Status   string `json:"information status"`
}

type TopUpByUser struct {
	UserCredential GetAuthResponse `json:"user credential"`
	UserBio        BiodataRequest  `json:"user biodata"`
	TopUp          []TopUp         `json:"top up detail"`
}

type TopUpById struct {
	UserCredential GetAuthResponse `json:"user credential"`
	UserBio        BiodataRequest  `json:"user biodata"`
	TopUp          TopUp           `json:"top up detail"`
}

type TopUp struct {
	Id                    string    `json:"id"`
	VaNumber              string    `json:"virtual account"`
	TopUpAmount           int       `json:"top up amount"`
	MaturityTime          time.Time `json:"maturity time"`
	AcceptedTime          time.Time `json:"accepted time"`
	Accepted              bool      `json:"accepted status"`
	Status                string    `json:"information status"`
	TransferConfirmRecipe bool      `json:"confirm recipe send"`
	File                  string    `json:"recipt file"`
}
