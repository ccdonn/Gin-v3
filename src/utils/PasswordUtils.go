package utils

import "golang.org/x/crypto/bcrypt"

// EncodePassword : encode password
func EncodePassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// VerifyPassword : verify password
func VerifyPassword(inputPassword, storedPasswordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(storedPasswordHash), []byte(inputPassword))
	return err == nil
}
