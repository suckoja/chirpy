package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/suckoja/chirpy/internal/httpjson"
)

type createUserRequest struct {
	Email string `json:"email"`
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = httpjson.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	user, err := h.db.CreateUser(r.Context(), req.Email)
	if err != nil {
		_ = httpjson.Error(w, http.StatusInternalServerError, "could not create user")
		return
	}

	_ = httpjson.Respond(w, http.StatusCreated, user)
}
