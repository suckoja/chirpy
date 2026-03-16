package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/suckoja/chirpy/internal/auth"
	"github.com/suckoja/chirpy/internal/httpjson"
)

type logInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	ExpiresInSeconds int `json:"expires_in_seconds"`
}

type loginResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
	Token string `json:"token"`
}

func (h *Handlers) LogIn(w http.ResponseWriter, r *http.Request) {
	var req logInRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = httpjson.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// 1. Find the user
	user, err := h.db.GetUserByEmail(r.Context(), req.Email)
	if err != nil {
		_ = httpjson.Error(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// 2. Verify the password
	match, err := auth.CheckPasswordHash(req.Password, user.HashedPassword)
	if err != nil || !match {
		_ = httpjson.Error(w, http.StatusUnauthorized, "Incorrect email or password")
		return
	}

	// 3. If it's specified by the client, use it as the expiration time. 
	// If it's not specified, use a default expiration time of 1 hour. 
	// If the client specified a number over 1 hour, use 1 hour as the expiration time.
	expiresIn := time.Hour 
	if req.ExpiresInSeconds > 0 && req.ExpiresInSeconds <= 3600 {
		expiresIn = time.Duration(req.ExpiresInSeconds) * time.Second
	}

	// 4. Make a JWT Token
	jwt, err := auth.MakeJWT(user.ID, h.jwtSecret, expiresIn) 
	if err != nil {
		_ = httpjson.Error(w, http.StatusInternalServerError, "Failed to make jwt")
		return
	}

	_ = httpjson.Respond(w, http.StatusOK, loginResponse{
		ID: user.ID.String(),
		Email: user.Email,
		Token: jwt,
	})
}