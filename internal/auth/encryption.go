package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var defaultCost int = 10
var issuer string = "chirpy"

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

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) { 
    claims := &jwt.RegisteredClaims{
        Issuer: "chirpy",
        IssuedAt: jwt.NewNumericDate(time.Now()),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
        Subject: userID.String() ,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString([]byte(tokenSecret))
    if err != nil {
        return "", err
    }

    return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
    claims := &jwt.RegisteredClaims{}
    token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
        return []byte(tokenSecret), nil
    })
    if err != nil {
        return uuid.Nil, err
    }
    
    sub, err := token.Claims.GetSubject()
    if err != nil {
        return uuid.Nil, err
    }

    userId, err := uuid.Parse(sub)
    if err != nil {
        return uuid.Nil, err
    }
    return userId, err
}
