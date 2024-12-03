package hashing

import (
	"testing"
)

// TestHashPassword проверяет, что хеширование пароля возвращает корректный результат.
func TestHashPassword(t *testing.T) {
	password := "securePassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(hash) == 0 {
		t.Fatalf("expected non-empty hash, got empty hash")
	}
}

// TestVerifyPassword проверяет, что пароль проходит проверку для корректного хэша.
func TestVerifyPassword(t *testing.T) {
	password := "securePassword123"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Убедимся, что VerifyPassword возвращает true для правильного пароля.
	if !VerifyPassword(password, hash) {
		t.Fatalf("expected password to match hash, but it did not")
	}
}

// TestVerifyPasswordIncorrect проверяет, что некорректный пароль не проходит проверку.
func TestVerifyPasswordIncorrect(t *testing.T) {
	password := "securePassword123"
	incorrectPassword := "wrongPassword"

	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// Убедимся, что VerifyPassword возвращает false для неправильного пароля.
	if VerifyPassword(incorrectPassword, hash) {
		t.Fatalf("expected password not to match hash, but it did")
	}
}
