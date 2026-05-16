package auth

import (
	"testing"
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
