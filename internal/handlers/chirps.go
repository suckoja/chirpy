package handlers

import (
	"database/sql"
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

func toChirpResponse(chirp database.Chirp) chirpResponse {
	return chirpResponse{
		ID:        chirp.ID.String(),
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID.String(),
	}
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

	_ = httpjson.Respond(w, http.StatusCreated, toChirpResponse(chirp))
}

func (h *Handlers) ListChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := h.db.ListChirps(r.Context())
	if err != nil {
		_ = httpjson.Error(w, http.StatusInternalServerError, "could not fetch chirps")
		return
	}

	resp := make([]chirpResponse, 0, len(chirps))
	for _, chirp := range chirps {
		resp = append(resp, toChirpResponse(chirp))
	}

	_ = httpjson.Respond(w, http.StatusOK, resp)
}

func (h *Handlers) GetChirp(w http.ResponseWriter, r *http.Request) {
	chirpIDStr := r.PathValue("chirpID")
	chirpID, err := uuid.Parse(chirpIDStr)
	if err != nil {
		_ = httpjson.Error(w, http.StatusBadRequest, "invalid chirp id")
		return
	}

	chirp, err := h.db.GetChirp(r.Context(), chirpID)
	if err != nil {
		if err == sql.ErrNoRows {
			_ = httpjson.Error(w, http.StatusNotFound, "chirp not found")
			return
		}
		_ = httpjson.Error(w, http.StatusInternalServerError, "could not fetch chirp")
		return
	}

	_ = httpjson.Respond(w, http.StatusOK, toChirpResponse(chirp))
}
