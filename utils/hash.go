package utils

import (
    "golang.org/x/crypto/bcrypt"
)

// HashPassword transforma uma senha em um hash
func HashPassword(password string) (string, error) {
    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", err
    }
    return string(hashedPassword), nil
}

// CheckPasswordHash compara uma senha com seu hash
func IsPasswordEqualHash(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}