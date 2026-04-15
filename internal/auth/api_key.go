package auth

import (
	"errors"
	"net/http"
	"strings"
)

func GetAPIKey(headers http.Header) (string, error) {
	apiKey := headers.Get("Authorization")
	if apiKey == "" {
		return "", errors.New("Missing Authorization header")
	}
	apiKeyNoPrefix, ok := strings.CutPrefix(apiKey, "ApiKey ")
	if !ok {
		return "", errors.New("Error cleaning token string")
	}
	trimmedApiKey := strings.TrimSpace(apiKeyNoPrefix)
	if trimmedApiKey == "" {
		return trimmedApiKey, errors.New("ApiKey is empty")
	}
	return trimmedApiKey, nil
}
