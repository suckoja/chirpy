package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"unicode/utf8"

	"github.com/suckoja/chirpy/internal/httpjson"
)

type validateChirpRequest struct {
	Body string `json:"body"`
}
type validateChirpResponse struct {
	CleanBody string `json:"cleaned_body"`
}

var profanes = map[string]struct{}{
	"kerfuffle": {}, 
	"sharbert": {},
	"fornax": {},
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

func (h *Handlers) ValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var req validateChirpRequest
	dec := json.NewDecoder(r.Body)

	if err := dec.Decode(&req); err != nil {
		_ = httpjson.Error(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	valid, reason := ValidateChirpText(req.Body)
	if !valid {
		_ = httpjson.Error(w, http.StatusBadRequest, reason)
		return
	}

	cleanBody := CleanRequestBody(req.Body)
	_ = httpjson.Respond(w, http.StatusOK, validateChirpResponse{CleanBody: cleanBody})
}
