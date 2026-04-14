package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	hash1, _ := HashPassword("mySecretPassword")
	hash2, _ := HashPassword("NotmySecretPassword")

	tests := []struct {
		name          string
		password      string
		hash          string
		wantErr       bool
		matchPassword bool
	}{
		{
			name:          "Correct password matches its own hash",
			password:      "mySecretPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: true,
		},
		{
			name:          "Wrong password doesn't match",
			password:      "NotmySecretPassword",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Password doesn't match a different user's hash",
			password:      "mySecretPassword",
			hash:          hash2,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Empty password doesn't match",
			password:      "",
			hash:          hash1,
			wantErr:       false,
			matchPassword: false,
		},
		{
			name:          "Invalid hash returns error",
			password:      "NotmySecretPassword",
			hash:          "Garbage string",
			wantErr:       true,
			matchPassword: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			match, err := CheckPasswordHash(tt.password, tt.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error checking hash when wanted none: %v", err)
			}
			if !tt.wantErr && match != tt.matchPassword {
				t.Errorf("got %v, want %v", match, tt.matchPassword)
			}
		})
	}
}

func TestMakeAndValidateJWT(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	tokenString, err := MakeJWT(userID, secret, time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned unexpected error: %v", err)
	}

	gotID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("ValidateJWT returned unexpected error: %v", err)
	}

	if gotID != userID {
		t.Errorf("got user ID %v, want %v", gotID, userID)
	}
}

func TestValidateJWT_ExpiredToken(t *testing.T) {
	userID := uuid.New()
	secret := "test-secret"

	// Create a token that expires immediately (negative duration puts it in the past)
	tokenString, err := MakeJWT(userID, secret, -time.Second)
	if err != nil {
		t.Fatalf("MakeJWT returned unexpected error: %v", err)
	}

	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Error("expected error for expired token, got nil")
	}
}

func TestValidateJWT_WrongSecret(t *testing.T) {
	userID := uuid.New()

	tokenString, err := MakeJWT(userID, "correct-secret", time.Hour)
	if err != nil {
		t.Fatalf("MakeJWT returned unexpected error: %v", err)
	}

	_, err = ValidateJWT(tokenString, "wrong-secret")
	if err == nil {
		t.Error("expected error for wrong secret, got nil")
	}
}

func TestValidateJWT_InvalidToken(t *testing.T) {
	_, err := ValidateJWT("this.is.not.a.valid.jwt", "some-secret")
	if err == nil {
		t.Error("expected error for invalid token string, got nil")
	}
}
