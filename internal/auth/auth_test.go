package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCheckpasswordHash(t *testing.T) {
	cases := []struct {
		hash     string
		expected string
	}{
		{
			hash:     "$argon2id$v=19$m=65536,t=1,p=12$Wqy+XYFBdO8TyI5+H3MaFw$0niFMML0m0celkB9gZXDviCarF/aLy/ygE67YSk52h0",
			expected: "$argon2id$v=19$m=65536,t=1,p=12$Wqy+XYFBdO8TyI5+H3MaFw$0niFMML0m0celkB9gZXDviCarF/aLy/ygE67YSk52h0",
		},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			CheckPasswordHash(tt.hash, tt.expected)
		})
	}
}

func TestHashPassword(t *testing.T) {
	cases := []struct {
		password string
		expected string
	}{
		{password: "123", expected: "$argon2id$v=19$m=65536,t=1,p=12$Wqy+XYFBdO8TyI5+H3MaFw$0niFMML0m0celkB9gZXDviCarF/aLy/ygE67YSk52h0"},
	}

	for _, tt := range cases {
		t.Run("", func(t *testing.T) {
			hashedPassword, _ := HashPassword(tt.password)
			CheckPasswordHash(hashedPassword, tt.expected)
		})
	}
}

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret", time.Hour)

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "Valid token",
			tokenString: validToken,
			tokenSecret: "secret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid token",
			tokenString: "invalid.token.string",
			tokenSecret: "secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
		{
			name:        "Wrong secret",
			tokenString: validToken,
			tokenSecret: "wrong_secret",
			wantUserID:  uuid.Nil,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, tt.wantUserID)
			}
		})
	}
}

func TestGetBearerToken(t *testing.T) {
	cases := []struct {
		name        string
		auth_header map[string][]string
		token       string
		wantErr     bool
	}{
		{
			name:        "regular authorization",
			auth_header: map[string][]string{"Authorization": []string{"Bearer 123456"}},
			token:       "123456",
			wantErr:     false,
		},
		{
			name:        "authorization not present",
			auth_header: map[string][]string{},
			token:       "",
			wantErr:     true,
		},
		{
			name:        "no bearer",
			auth_header: map[string][]string{"Authorization": []string{"123456"}},
			token:       "",
			wantErr:     true,
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.auth_header)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBearerToken error = %v, wantErr = %v", err, tt.wantErr)
			}

			if token != tt.token {
				t.Errorf("GetBearerToken token = %v, want = %v", token, tt.token)
			}
		})
	}
}
