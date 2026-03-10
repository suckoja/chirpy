package httpjson

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRespond(t *testing.T) {
	tests := []struct {
		name           string
		status         int
		payload        any
		expectedStatus int
		expectedBody   string
	}{
		{"200 with struct", http.StatusOK, map[string]string{"hello": "world"}, http.StatusOK, `{"hello":"world"}`},
		{"201 created", http.StatusCreated, map[string]string{"id": "123"}, http.StatusCreated, `{"id":"123"}`},
		{"500 with nil", http.StatusInternalServerError, nil, http.StatusInternalServerError, `null`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := Respond(w, tt.status, tt.payload)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			res := w.Result()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}
			if ct := res.Header.Get("Content-Type"); ct != "application/json; charset=utf-8" {
				t.Errorf("expected Content-Type application/json, got %q", ct)
			}
			if res.Header.Get("Access-Control-Allow-Origin") != "*" {
				t.Error("expected CORS header Access-Control-Allow-Origin: *")
			}

			var body any
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatalf("could not decode response body: %v", err)
			}
			got, _ := json.Marshal(body)
			if string(got) != tt.expectedBody {
				t.Errorf("expected body %s, got %s", tt.expectedBody, got)
			}
		})
	}
}

func TestError(t *testing.T) {
	tests := []struct {
		name           string
		code           int
		msg            string
		expectedStatus int
	}{
		{"bad request", http.StatusBadRequest, "invalid input", http.StatusBadRequest},
		{"not found", http.StatusNotFound, "not found", http.StatusNotFound},
		{"internal error", http.StatusInternalServerError, "something went wrong", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			err := Error(w, tt.code, tt.msg)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			res := w.Result()

			if res.StatusCode != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, res.StatusCode)
			}

			var body map[string]string
			if err := json.NewDecoder(res.Body).Decode(&body); err != nil {
				t.Fatalf("could not decode response body: %v", err)
			}
			if body["error"] != tt.msg {
				t.Errorf("expected error %q, got %q", tt.msg, body["error"])
			}
		})
	}
}
