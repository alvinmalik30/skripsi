package model

type UserCredential struct {
	Id       string
	Username string
	Email    string
	Password string
	VANumber string
	Role     string
	IsActive bool
}
