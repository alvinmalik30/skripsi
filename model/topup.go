package model

import "time"

type TopUp struct {
	Id                    string
	UserCredential        UserCredential
	TopUpAmount           int
	MaturityTime          time.Time
	AcceptedTime          time.Time
	Accepted              bool
	Status                string
	TransferConfirmRecipe bool
	File                  string
}
