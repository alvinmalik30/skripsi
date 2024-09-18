package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	// ambil pass dan jadikan hash
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed hash password: %v", err.Error())
	}
	return string(passwordHash), nil
}

func VerifyPassword(hashPassword string, password string) error {
	// ambil hash pass dari database dan samakan apakah sama dengan pass yang di input
	return bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(password))
}
