package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestNonEmptyStringJWT(t *testing.T) {
	id := uuid.New()
	jwt, err := MakeJWT(id, "thisIsASecret", time.Hour) 
	if err != nil {
		t.Fatalf("failed to MakeJWT: %v", err)
	}
	if jwt == "" {
		t.Fatal("jwt return empty string")
	}
}

func TestValidJWT(t *testing.T) {
	userID := uuid.New()
	secret := "thisIsASecret"
	jwt, err := MakeJWT(userID, secret, time.Hour) 
	if err != nil {
		t.Fatalf("failed to make jwt: %v", err)
	}

	validateResult, err := ValidateJWT(jwt, secret)
	if err != nil {
		t.Fatalf("failed to validate jwt: %v", err)
	}
	if userID != validateResult {
		t.Fatal("jwt is not valid")
	}
}

// Wrong secret
func TestWrongSecret(t *testing.T) {
	userID := uuid.New()
	secret := "thisIsASecret"
	
	jwt, err := MakeJWT(userID, secret, time.Hour) 
	if err != nil {
		t.Fatalf("failed to make jwt: %v", err)
	}

	wrongSecret := "ThisIsAWrongSecret!"
	_, err = ValidateJWT(jwt, wrongSecret)
	if err == nil {
		t.Fatal("expected an error but got none")
	}
}

// Expired token
func TestExpiredJWT(t *testing.T) {
	userID := uuid.New()
	secret := "thisIsASecret"
	
	// Create a token that expired 1 second ago
	jwt, err := MakeJWT(userID, secret, -time.Second) 
	if err != nil {
		t.Fatalf("failed to make jwt: %v", err)
	}

	// Validating it should fail because it's expired
	_, err = ValidateJWT(jwt, secret)
	if err == nil {
		t.Fatal("expected an error for expired token but got none")
	}
}

// Tampered token
func TestTamperedJWT(t *testing.T) {
	userID := uuid.New()
	secret := "thisIsASecret"
	
	// Create a token
	jwt, err := MakeJWT(userID, secret, time.Second) 
	if err != nil {
		t.Fatalf("failed to make jwt: %v", err)
	}

	// Tampered jwt result
	jwt = jwt[:len(jwt)-1] + "AB"

	// Validating it should fail because it's tampered
	_, err = ValidateJWT(jwt, secret)
	if err == nil {
		t.Fatal("expected an error for tampered token but got none")
	}
}

// Empty/garbage string
func TestEmptyOrGarbageString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
	}{
		{"empty string", ""},
		{"purposely not a jwt string", "not.a.jwt"},
		{"another not jwt string", "abs.def.ghi"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			secret := "thisIsASecret"
			_, err := ValidateJWT(tt.input, secret)
			if err == nil {
				t.Fatal("expected an error for tampered token but got none")
			}
		})
	}
}