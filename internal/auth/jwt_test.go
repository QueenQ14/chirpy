package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestJWTValidToken(t *testing.T) {
	userID := uuid.New()
	secret := "secret-key"
	expiryTime := 30 * time.Minute

	//generate token using MakeJWT
	token, err := MakeJWT(userID, secret, expiryTime)
	if err != nil {
		t.Fatalf("error generating token: %v", err)
	}

	//test if token is non-empty
	if token == "" {
		t.Fatalf("generated token is empty")
	}

	//calling ValidateJWT
	retrievedID, err := ValidateJWT(token, secret)
	if err != nil {
		t.Fatalf("error in validating the token: %v", err)
	}

	if retrievedID != userID {
		t.Fatalf("expected userID: %v, got: %v", userID, retrievedID)
	}
}

func TestExpiredJWTToken(t *testing.T) {
	//Making an expired JWT Token
	userID := uuid.New()
	secret := "secret-key"
	expiryTime := -1 * time.Hour
	token, _ := MakeJWT(userID, secret, expiryTime)

	// ValidateJWT should return an error since the token is expired
	retrievedID, err := ValidateJWT(token, secret)
	if retrievedID != uuid.Nil {
		t.Errorf("Retrieved ID should be nil UUID, got: %v", retrievedID)
	}
	if err == nil {
		t.Fatalf("expected token to be expired, token still alive")
	}
}

func TestDifferentSecrets(t *testing.T) {
	userID := uuid.New()
	secret := "secret-key"
	secretTwo := "not-so-secret-key"
	// Creating token with 'secret'
	token, _ := MakeJWT(userID, secret, time.Hour)

	// ValidateJWT with 'secretTwo' instead of 'secret'
	retrievedID, err := ValidateJWT(token, secretTwo)
	if retrievedID != uuid.Nil {
		t.Errorf("Retrieved ID should be nil UUID, got: %v", retrievedID)
	}
	if err == nil {
		t.Fatalf("expected validation to fail, but passed")
	}
}

func TestGarbageStringAsToken(t *testing.T) {
	token := "this.is.not.a.jwt"
	secret := "secret-key"
	retrievedID, err := ValidateJWT(token, secret)
	if retrievedID != uuid.Nil {
		t.Errorf("Retrieved ID should be nil UUID, got: %v", retrievedID)
	}
	if err == nil {
		t.Fatalf("expected validation to fail, but passed")
	}
}
