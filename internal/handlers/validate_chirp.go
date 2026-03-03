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
	Valid bool `json:"valid"`
}

// Pure function: easy to reuse + test
func ValidateChirpText(body string) (bool, string) {
	const limit = 140
	body = strings.TrimSpace(body)
	if utf8.RuneCountInString(body) > limit {
		return false, "Chirp is too long"
	}
	return true, ""
}

func (h *Handlers) ValidateChirp(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	
	var req validateChirpRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	if err := dec.Decode(&req); err != nil {
		_ = httpjson.Error(w, http.StatusBadRequest, "Invalid JSON body")
		return
	}

	valid, reason := ValidateChirpText(req.Body)
	if !valid {
		_ = httpjson.Error(w, http.StatusBadRequest, reason)
		return
	} 

	_ = httpjson.Respond(w, http.StatusOK, validateChirpResponse{Valid: true})
}
