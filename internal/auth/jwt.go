package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {

	// create the claims (payload) for the token
	claims := CustomClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer: "chirpy-access",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresIn)),
			Subject: userID.String(),
		},
	}

	// create a new token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// sign the token with the secret key
	tokenString, err := token.SignedString([]byte(tokenSecret))
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}

	return tokenString, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	// Parse and validate the token
	parsedToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Ensure the algorithm is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]) 
		}

		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing token: %w", err)
	}

	claims, ok := parsedToken.Claims.(*CustomClaims)
	if !ok || !parsedToken.Valid {
		return uuid.Nil, errors.New("token is invalid")
	}

	return claims.UserID, nil
}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader := headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("authorization header not present")
	}

	if !strings.HasPrefix(authHeader, "Bearer") {
		return "", errors.New("authorization does not begin with Bearer")
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")
	if token == "" {
		return "", errors.New("token should not be empty")
	} 

	return strings.Trim(token, " "), nil
}

func AuthenticateRequest(headers http.Header, tokenSecret string) (uuid.UUID, error) {
	// Authetication request with JWT
	// 1. Extract the token from the Authorization header
	tokenString, err := GetBearerToken(headers)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed at get bearer token: %w", err)
	}

	// 2. Validate it to get the userID
	userID, err := ValidateJWT(tokenString, tokenSecret)
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed at jwt validate: %w", err)
	}

	return userID, nil
}