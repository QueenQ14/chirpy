package auth

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
)

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("auth header missing")
	}
	if !strings.HasPrefix(auth, "Bearer ") {
		return "", fmt.Errorf("invalid auth header")
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	if token == "" {
		return "", fmt.Errorf("token missing")
	}
	return token, nil
}

func MakeRefreshToken() string {
	data := make([]byte, 32)
	if _, err := rand.Read(data); err != nil {
		return ""
	}
	return hex.EncodeToString(data)
}

func GetAPIKey(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	if auth == "" {
		return "", fmt.Errorf("auth header missing")
	}
	if !strings.HasPrefix(auth, "ApiKey ") {
		return "", fmt.Errorf("invalid auth header")
	}
	key := strings.TrimPrefix(auth, "ApiKey ")
	if key == "" {
		return "", fmt.Errorf("ApiKey missing")
	}
	return key, nil
}
