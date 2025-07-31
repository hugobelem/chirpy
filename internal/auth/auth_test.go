package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestCheckPasswordHash(t *testing.T) {
	// First, we need to create some hashed passwords for testing
	password1 := "Password123"
	password2 := "Password456"
	hash1, _ := HashPassword(password1)
	hash2, _ := HashPassword(password2)

	tests := []struct {
		name string
		password string
		hash string
		wantErr bool
	}{
		{
			name:     "Correct password",
			password: password1,
			hash:     hash1,
			wantErr:  false,
		},
		{
			name:     "Incorrect password",
			password: "wrongPassword",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Password doesn't match different hash",
			password: password1,
			hash:     hash2,
			wantErr:  true,
		},
		{
			name:     "Empty password",
			password: "",
			hash:     hash1,
			wantErr:  true,
		},
		{
			name:     "Invalid hash",
			password: password1,
			hash:     "invalidhash",
			wantErr:  true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := CheckPasswordHash(test.password, test.hash)
			if (err != nil) != test.wantErr {
				t.Errorf(
					"CheckPasswordHash() error = %v, wantErr %v",
					err,
					test.wantErr,
				)
			}
		})
	}
}

func TestValidadeJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "secret")

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
	
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			gotUserID, err := ValidateJWT(test.tokenString, test.tokenSecret)
			if (err != nil) != test.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if gotUserID != test.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", gotUserID, test.wantUserID)
			}
		})
	}
}