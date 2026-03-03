package httpjson

import (
	"encoding/json"
	"net/http"
)

func Respond(w http.ResponseWriter, status int, payload any) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(payload)
}

func Error(w http.ResponseWriter, code int, msg string) error {
	return Respond(w, code, map[string]string{"error": msg})
}