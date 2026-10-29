package auth

import (
	"net/http"
	"testing"
)

func TestGetBearerToken(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", "Bearer meowmeow")
	originalToken := "meowmeow"
	token, err := GetBearerToken(header)
	if err != nil {
		t.Fatalf("error fetching token from header")
	}
	if token != originalToken {
		t.Fatalf("token was not extracted correctly. Expected %v, got %v", originalToken, token)
	}
}

func TestGetAPIKey(t *testing.T) {
	header := http.Header{}
	header.Add("Authorization", "ApiKey meowmeow")
	originalToken := "meowmeow"
	token, err := GetAPIKey(header)
	if err != nil {
		t.Fatalf("error fetching token from header")
	}
	if token != originalToken {
		t.Fatalf("ApiKey was not extracted correctly. Expected %v, got %v", originalToken, token)
	}
}
