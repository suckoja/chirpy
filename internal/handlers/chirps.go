package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/suckoja/chirpy/internal/database"
	"github.com/suckoja/chirpy/internal/httpjson"
)

type createChirpRequest struct {
	Body   string `json:"body"`
	UserID string `json:"user_id"`
}

type chirpResponse struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at`
	UpdatedAt time.Time `json:"updated_at`
	Body      string    `json:"body"`
	UserID    string    `json:"user_id"`
}

var profanes = map[string]struct{}{
	"kerfuffle": {},
	"sharbert":  {},
	"fornax":    {},
}

func ValidateChirpText(body string) (bool, string) {
	const limit = 140
	body = strings.TrimSpace(body)
	if utf8.RuneCountInString(body) > limit {
		return false, "Chirp is too long"
	}
	return true, ""
}

func CleanRequestBody(body string) string {
	const replacement = "****"

	words := strings.Fields(body)

	for i, word := range words {
		if _, exists := profanes[strings.ToLower(word)]; exists {
			words[i] = replacement
		}
	}

	return strings.Join(words, " ")
}

func (h *Handlers) CreateChirp(w http.ResponseWriter, r *http.Request) {
	var req createChirpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		_ = httpjson.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	valid, reason := ValidateChirpText(req.Body)
	if !valid {
		_ = httpjson.Error(w, http.StatusBadRequest, reason)
		return
	}

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		_ = httpjson.Error(w, http.StatusBadRequest, "invalid user_id")
		return
	}

	chirp, err := h.db.CreateChirp(r.Context(), database.CreateChirpParams{
		Body:   CleanRequestBody(req.Body),
		UserID: userID,
	})
	if err != nil {
		_ = httpjson.Error(w, http.StatusInternalServerError, "could not create user")
		return
	}

	resp := chirpResponse{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}

	_ = httpjson.Respond(w, http.StatusCreated, resp)
}
