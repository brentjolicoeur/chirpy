package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	method := jwt.SigningMethodHS256
	claims := jwt.RegisteredClaims{
		Issuer:    "chirpy-access",
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
		Subject:   userID.String(),
	}
	newToken, err := jwt.NewWithClaims(method, claims).SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("Error creating token: %w", err)
	}
	return newToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(tokenSecret), nil
		})
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error validating token: %w", err)
	}
	id, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("Error parsing id: %w", err)
	}

	return id, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	tokenString := headers.Get("Authorization")
	if tokenString == "" {
		return "", errors.New("Missing Authorization header")
	}
	cleanedToken, ok := strings.CutPrefix(tokenString, "Bearer ")
	if !ok {
		return "", errors.New("Error cleaning token string")
	}
	trimmedToken := strings.TrimSpace(cleanedToken)
	if trimmedToken == "" {
		return trimmedToken, errors.New("Bearer string token is empty")
	}
	return trimmedToken, nil
}

func MakeRefreshToken() string {
	key := make([]byte, 32)
	rand.Read(key)
	return hex.EncodeToString(key)
}
