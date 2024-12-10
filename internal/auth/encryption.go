package auth

import (
    "golang.org/x/crypto/bcrypt"
)

var defaultCost int = 10

func HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), defaultCost)
    if err != nil {
        return "", err
    }

    return string(hash), nil
}

func CheckPasswordHash(password, hash string) error {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if err != nil {
        return err
    }

    return nil
}
