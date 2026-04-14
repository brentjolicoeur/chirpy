package auth

import (
	"net/http"
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

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		headers     http.Header
		expected    string
		expectError bool
	}{
		{
			name: "valid bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer my-secret-token"},
			},
			expected:    "my-secret-token",
			expectError: false,
		},
		{
			name: "valid bearer token with extra whitespace",
			headers: http.Header{
				"Authorization": []string{"Bearer   token-with-spaces   "},
			},
			expected:    "token-with-spaces",
			expectError: false,
		},
		{
			name:        "missing authorization header",
			headers:     http.Header{},
			expected:    "",
			expectError: true,
		},
		{
			name: "wrong scheme (Basic instead of Bearer)",
			headers: http.Header{
				"Authorization": []string{"Basic dXNlcjpwYXNz"},
			},
			expected:    "",
			expectError: true,
		},
		{
			name: "empty authorization header value",
			headers: http.Header{
				"Authorization": []string{""},
			},
			expected:    "",
			expectError: true,
		},
		{
			name: "bearer prefix only, no token",
			headers: http.Header{
				"Authorization": []string{"Bearer "},
			},
			expected:    "",
			expectError: true,
		},
		{
			name: "JWT-style bearer token",
			headers: http.Header{
				"Authorization": []string{"Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.payload.signature"},
			},
			expected:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.payload.signature",
			expectError: false,
		},
		{
			name:        "nil headers map",
			headers:     nil,
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetBearerToken(tt.headers)

			if tt.expectError && err == nil {
				t.Errorf("expected an error but got none")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if got != tt.expected {
				t.Errorf("got %q, want %q", got, tt.expected)
			}
		})
	}
}
