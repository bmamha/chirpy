package auth

import (
	"testing"

	"github.com/google/uuid"
)

func TestValidateJWT(t *testing.T) {
	userID := uuid.New()
	validToken, _ := MakeJWT(userID, "sec")

	tests := []struct {
		name        string
		tokenString string
		tokenSecret string
		wantUserID  uuid.UUID
		wantErr     bool
	}{
		{
			name:        "validToken",
			tokenString: validToken,
			tokenSecret: "tokenSecret",
			wantUserID:  userID,
			wantErr:     false,
		},
		{
			name:        "Invalid Token",
			tokenString: "invalid token string",
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
			validatedUserID, err := ValidateJWT(tt.tokenString, tt.tokenSecret)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateJWT error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if validatedUserID != tt.wantUserID {
				t.Errorf("ValidateJWT() gotUserID = %v, want %v", validatedUserID, tt.wantUserID)
			}
		})
	}
}
