package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/suckoja/chirpy/internal/auth"
	"github.com/suckoja/chirpy/internal/httpjson"
)

type logInRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
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

	_ = httpjson.Respond(w, http.StatusOK, toUserResponse(user))
}