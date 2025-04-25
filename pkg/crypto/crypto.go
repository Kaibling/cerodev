package crypto

import (
	"crypto/rand"
	"encoding/hex"

	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string, passwordCost int) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), passwordCost)

	return string(bytes), err
}

func CheckPasswordHash(password, hash string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)); err != nil {
		return false, err
	}

	return true, nil
}

func GenerateToken(length int) (string, error) {
	key := make([]byte, length)

	_, err := rand.Read(key)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(key), nil
}
