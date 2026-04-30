package utils

import (
	"testing"
)

func TestPasswordHashAndCheck(t *testing.T) {
	tests := []struct {
		name     string
		password string
	}{
		{
			name:     "standard password",
			password: "mysecretpassword",
		},
		{
			name:     "short password",
			password: "123",
		},
		{
			name:     "password with symbols",
			password: "password!@#$%^&*()",
		},
		{
			name:     "empty password",
			password: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test PasswordHash
			hash, err := PasswordHash(tt.password)
			if err != nil {
				t.Fatalf("PasswordHash() failed for %q: %v", tt.name, err)
			}

			if hash == "" {
				t.Errorf("PasswordHash() returned empty hash for %q", tt.name)
			}

			// Test CheckPasswordHash with correct password
			if !CheckPasswordHash(tt.password, hash) {
				t.Errorf("CheckPasswordHash() failed for %q with correct password", tt.name)
			}

			// Test CheckPasswordHash with incorrect password
			wrongPassword := tt.password + "_wrong"
			if CheckPasswordHash(wrongPassword, hash) {
				t.Errorf("CheckPasswordHash() should have failed for %q with incorrect password", tt.name)
			}
		})
	}
}

func TestCheckPasswordHashInvalidHash(t *testing.T) {
	tests := []struct {
		name     string
		password string
		hash     string
		want     bool
	}{
		{
			name:     "invalid hash format",
			password: "password",
			hash:     "not-a-bcrypt-hash",
			want:     false,
		},
		{
			name:     "empty hash",
			password: "password",
			hash:     "",
			want:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CheckPasswordHash(tt.password, tt.hash); got != tt.want {
				t.Errorf("CheckPasswordHash() got = %v, want %v for %q", got, tt.want, tt.name)
			}
		})
	}
}
