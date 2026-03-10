package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/suckoja/chirpy/internal/auth"
	"github.com/suckoja/chirpy/internal/database"
	"github.com/suckoja/chirpy/internal/httpjson"
)

type createUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func toUserResponse(user database.User) userResponse {
	return userResponse{
		ID:        user.ID.String(),
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
}

func (h *Handlers) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req createUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = httpjson.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		_ = httpjson.Error(w, http.StatusInternalServerError, "could not hash a password")
		return
	}

	user, err := h.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          req.Email,
		HashedPassword: hashedPassword,
	})
	if err != nil {
		_ = httpjson.Error(w, http.StatusInternalServerError, "could not create user")
		return
	}

	_ = httpjson.Respond(w, http.StatusCreated, toUserResponse(user))
}
