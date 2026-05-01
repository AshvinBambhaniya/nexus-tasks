package jwt

import (
	"testing"
	"time"

	"github.com/AshvinBambhaniya/nexus-tasks/config"
	"github.com/stretchr/testify/assert"
)

func TestCreateToken(t *testing.T) {
	cfg := config.AppConfig{Secret: "test-secret-key-that-is-long-enough-for-hs256"}

	tests := []struct {
		name    string
		sub     string
		exp     time.Time
		wantErr bool
	}{
		{
			name:    "valid token",
			sub:     "user-123",
			exp:     time.Now().Add(time.Hour),
			wantErr: false,
		},
		{
			name:    "token with empty subject",
			sub:     "",
			exp:     time.Now().Add(time.Hour),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := CreateToken(cfg.Secret, tt.sub, tt.exp)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("CreateToken() returned empty token")
			}
		})
	}
}

func TestParseToken(t *testing.T) {
	secret := "test-secret-key-that-is-long-enough-for-hs256"
	cfg := config.AppConfig{Secret: secret}

	validSub := "user-123"
	validToken, _ := CreateToken(cfg.Secret, validSub, time.Now().Add(time.Hour))

	expiredToken, _ := CreateToken(cfg.Secret, validSub, time.Now().Add(-time.Hour))

	wrongSecretCfg := config.AppConfig{Secret: "wrong-secret-key-long-enough-for-hs256"}
	wrongSecretToken, _ := CreateToken(wrongSecretCfg.Secret, validSub, time.Now().Add(time.Hour))

	tests := []struct {
		name    string
		token   string
		wantSub string
		wantErr bool
	}{
		{
			name:    "valid token",
			token:   validToken,
			wantSub: validSub,
			wantErr: false,
		},
		{
			name:    "expired token",
			token:   expiredToken,
			wantErr: true,
		},
		{
			name:    "invalid signature (wrong secret)",
			token:   wrongSecretToken,
			wantErr: true,
		},
		{
			name:    "malformed token",
			token:   "not.a.token",
			wantErr: true,
		},
		{
			name:    "empty token",
			token:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ParseToken(cfg.Secret, tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				assert.Equal(t, tt.wantSub, claims.Subject())
			}
		})
	}
}
