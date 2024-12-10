package auth

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

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
