package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var secret string = "lzoAD+JoSY7HJrKmtuYDfB67sypYiWeStwDQhxVKzkfJ3++4Jkzoh63/kay5pG1DCs/je8s97ov+VfTMkzbcGA=="

func TestPasswordHashed(t *testing.T) {
    pass := "password"

    passHash, err := HashPassword(pass)
    if err != nil {
        t.Fatalf(err.Error())
    }
    
    err = bcrypt.CompareHashAndPassword([]byte(passHash), []byte(pass))
    if err != nil {
        t.Fatalf("Produced hash doesn't equal expected!")
    }
}

func TestHashVerification(t *testing.T) {
    pass := "password"
    hash := "$2a$10$d03WIijVOR1cVww3Lq.9fekb1qsAtBTMnIT10z1iVDg1lkLA6L4f."

    err := CheckPasswordHash(pass, hash)
    if err != nil {
        t.Fatalf("Password validation against hash failed %v", err)
    }
}

func TestInvalidHashVerification(t *testing.T) {
    pass := "hello-world"
    invalidHash := "$2y$10$M8NsO9Km6XhztdnyXVf1VuGC6obfKqxfX5ep8SFb5CY0VEAdDsDza" //hello-worldz

    err := CheckPasswordHash(pass, invalidHash)
    if err == nil {
        t.Fatalf("Password validated against invalid hash")
    }
}

func TestCreateValidToken(t *testing.T) {
    expiresIn, err := time.ParseDuration("10s")
    if err != nil {
        t.Fatalf("Error parsing duration\n")
    }
    userID, err := uuid.Parse("35c8d40b-f8d7-48d2-9211-33e1ce5aa272")
    if err != nil {
        t.Fatalf("Error parsing uuid\n")
    }

    token, err := MakeJWT(userID, secret, expiresIn)
    if err != nil {
        t.Fatalf("Error creating token\n")
    }

    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })
    if err != nil || !parsedToken.Valid {
        t.Fatalf("Token is not valid\n")
    }
}

func TestCreateExpiredToken(t *testing.T) {
    expiresIn := time.Duration(-10)
    userID, err := uuid.Parse("35c8d40b-f8d7-48d2-9211-33e1ce5aa272")
    if err != nil {
        t.Fatalf("Error parsing uuid\n")
    }

    token, err := MakeJWT(userID, secret, expiresIn)
    if err != nil {
        t.Fatalf("Error creating token\n")
    }

    parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
        return []byte(secret), nil
    })

    if parsedToken.Valid {
        t.Fatalf("Token is not valid\n")
    }
}

func TestJWTValidates(t *testing.T) {
    //Signed with var secret string
    expiresIn, err := time.ParseDuration("10s")
    if err != nil {
        t.Fatalf("Error parsing duration\n")
    }

    expectedUserId := "35c8d40b-f8d7-48d2-9211-33e1ce5aa272"
    claims := &jwt.RegisteredClaims{
        Issuer: issuer,
        IssuedAt: jwt.NewNumericDate(time.Now()),
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
        Subject: expectedUserId,
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

    tokenString, err := token.SignedString([]byte(secret))
    if err != nil {
        t.Fatalf(err.Error())
    }

    actualUserId, err := ValidateJWT(tokenString, secret)
    if err != nil {
        t.Fatalf("Actual userID is not equal expected expected: %v actual: %v\n", expectedUserId, actualUserId)
    }
}
