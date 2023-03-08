package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerificationPasswordAndHash(password, hash string) error {
	errCompare := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if errCompare != nil {
		fmt.Println(errCompare.Error())
		return errCompare
	}

	return nil
}
