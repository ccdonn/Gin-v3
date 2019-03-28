package utils

import "golang.org/x/crypto/bcrypt"

func EncodePassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func VerifyPassword(inputPassword, storedPasswordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(inputPassword))
	return err == nil
}
