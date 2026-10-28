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
