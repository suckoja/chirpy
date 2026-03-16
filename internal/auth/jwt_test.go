package auth

import (
	"net/http/httptest"
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

// Test GetBearerToken
func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name string
		authHeader string
		expectedToken string
		wantErr bool
	}{
		{
			name: "Valid Bearer Token",
			authHeader: "Bearer validtoken123",
			expectedToken: "validtoken123",
			wantErr: false,
		},
		{
			name: "Invalid Format - Missing Bearer Prefix",
			authHeader: "invalidtoken456",
			expectedToken: "",
			wantErr: true,
		},
		{
			name: "Empty Header",
			authHeader: "",
			expectedToken: "",
			wantErr: true,
		},
		{
			name: "Bearer Prefix only, no token",
			authHeader: "Bearer ",
			expectedToken: "",
			wantErr: true,
		},
		{
			name: "Not a proper Bearer",
			authHeader: "IAmBearer someToken789",
			expectedToken: "",
			wantErr: true,
		},
		{
			name: "Token with extra spaces",
			authHeader: "Bearer        io238022j3oude90823u",
			expectedToken: "io238022j3oude90823u",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock request using httptest
			req := httptest.NewRequest("GET", "/", nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Get bearer token
			getToken, err := GetBearerToken(req.Header)

			// Check return error match wantErr or not
			if (err != nil) != tt.wantErr {
				t.Fatalf("GetBearerToken() want error %v, get %v", tt.wantErr, err)
			}

			// Check return token match expectedToken or not
			if getToken != tt.expectedToken {
				t.Errorf("GetBearerToken() = %v, want %v", getToken, tt.expectedToken)
			}
		})
	}
}

func TestValidAuthenticateRequest(t *testing.T) {
	// prepare userID & tokenString
	userID := uuid.New()
	tokenSecret := "ThisIsABloodySecret"
	tokenString, err := MakeJWT(userID, tokenSecret, time.Hour)
	if err != nil {
		t.Fatalf("failed to make jwt: %v", err)
	}

	// Create a mock request using httptest
	req := httptest.NewRequest("GET", "/", nil)
	if tokenString != "" {
		req.Header.Set("Authorization", "Bearer " + tokenString)
	}

	resultID, err := AuthenticateRequest(req.Header, tokenSecret)
	if err != nil {
		t.Fatalf("failed to authenticate request: %v", err)
	}
	if userID != resultID{
		t.Errorf("user id does not matched")
	}
}

func TestAuthenticateRequestWithMissingHeader(t *testing.T) {
	// Create a mock request with no bearer
	req := httptest.NewRequest("GET", "/", nil)
	_, err := AuthenticateRequest(req.Header, "secret")
	if err == nil {
		t.Fatal("expected an error for missing header but got none")
	}
}