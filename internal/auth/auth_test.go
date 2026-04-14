package auth

import (
	"testing"
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
