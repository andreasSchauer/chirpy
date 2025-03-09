package auth

import (
	"fmt"
    "golang.org/x/crypto/bcrypt"
)


func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("couldn't generate hash from password: %v", err)
	}

	return string(hash), nil
}


func CheckPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}