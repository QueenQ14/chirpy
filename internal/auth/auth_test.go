package auth

import "testing"

func TestPasswordHashing(t *testing.T) {
	password := "mySecretPassword123"

	// 1. Test that hashing succeeds
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("unexpected error hashing password: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash string")
	}

	// 2. Test that correct password validates to true
	match, err := CheckPasswordHash(password, hash)
	if err != nil {
		t.Fatalf("unexpected error checking password: %v", err)
	}
	if !match {
		t.Errorf("expected password to match hash, but it did not")
	}

	// 3. Test that wrong password validates to false
	match, err = CheckPasswordHash("wrongPassword", hash)
	if err != nil {
		t.Fatalf("unexpected error checking wrong password: %v", err)
	}
	if match {
		t.Errorf("expected wrong password to fail verification, but it matched")
	}
}

func TestPasswordSaltRandomness(t *testing.T) {
	password := "samePassword"

	// Argon2 generates a random salt each time, so two hashes of the same password must differ
	hash1, err := HashPassword(password)
	if err != nil {
		t.Fatalf("error hashing password 1: %v", err)
	}

	hash2, err := HashPassword(password)
	if err != nil {
		t.Fatalf("error hashing password 2: %v", err)
	}

	if hash1 == hash2 {
		t.Errorf("expected different hashes due to random salt, but got identical hashes")
	}
}
