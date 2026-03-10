package handlers

import "testing"

func TestValidateChirpText(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		valid   bool
		message string
	}{
		{"empty string", "", true, ""},
		{"normal chirp", "Hello world!", true, ""},
		{"exactly 140 chars", string(make([]byte, 140)), true, ""},
		{"141 chars", string(make([]byte, 141)), false, "Chirp is too long"},
		{"whitespace trimmed under limit", "  hello  ", true, ""},
		{"unicode characters", "こんにちは世界", true, ""},
		{"unicode over limit", string(make([]rune, 141)), false, "Chirp is too long"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			valid, msg := ValidateChirpText(tt.input)
			if valid != tt.valid {
				t.Errorf("expected valid=%v, got %v", tt.valid, valid)
			}
			if msg != tt.message {
				t.Errorf("expected message=%q, got %q", tt.message, msg)
			}
		})
	}
}

func TestCleanRequestBody(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"no profanity", "Hello world", "Hello world"},
		{"kerfuffle lowercase", "what a kerfuffle today", "what a **** today"},
		{"sharbert lowercase", "sharbert is bad", "**** is bad"},
		{"fornax lowercase", "fornax ruins things", "**** ruins things"},
		{"mixed case profanity", "What a Kerfuffle!", "What a Kerfuffle!"},
		{"multiple profanities", "kerfuffle sharbert fornax", "**** **** ****"},
		{"profanity mid-sentence", "this kerfuffle is real", "this **** is real"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CleanRequestBody(tt.input)
			if result != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, result)
			}
		})
	}
}
