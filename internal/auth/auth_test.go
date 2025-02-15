package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeAndValidateToken(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"

	tokenString, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("unexpected error while creating JWT: %v", err)
	}
	validatedUserID, err := ValidateJWT(tokenString, secret)
	if err != nil {
		t.Fatalf("unexpected error while validating JWT: %v", err)
	}

	if validatedUserID != userID {
		t.Errorf("unexpected userID %v, got %v", userID, validatedUserID)
	}
}

func TestExpiredTokens(t *testing.T) {
	userID := uuid.New()
	secret := "test_secret"

	tokenString, err := MakeJWT(userID, secret)
	if err != nil {
		t.Fatalf("unexpected error while creating JWT: %v", err)
	}
	time.Sleep(time.Second)
	_, err = ValidateJWT(tokenString, secret)
	if err == nil {
		t.Fatalf("expected %v error, got: %v", jwt.ErrTokenInvalidClaims, err)
	}
}
