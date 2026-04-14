package auth

import (
	"fmt"
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
	return uuid.Nil, nil
}
