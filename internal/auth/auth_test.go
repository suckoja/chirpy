package auth

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if hash == "" {
		t.Fatal("expected non-empty hash")
	}
	if hash == "password123" {
		t.Fatal("hash should not equal plaintext password")
	}
}

func TestHashPasswordIsUnique(t *testing.T) {
	hash1, _ := HashPassword("password123")
	hash2, _ := HashPassword("password123")
	if hash1 == hash2 {
		t.Fatal("expected different hashes for same password (argon2id uses random salt)")
	}
}

func TestCheckPasswordHash(t *testing.T) {
	hash, _ := HashPassword("password123")

	match, err := CheckPasswordHash("password123", hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !match {
		t.Fatal("expected password to match hash")
	}
}

func TestCheckPasswordHashWrongPassword(t *testing.T) {
	hash, _ := HashPassword("password123")

	match, err := CheckPasswordHash("wrongpassword", hash)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if match {
		t.Fatal("expected password to not match hash")
	}
}
