package hashing

import (
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	secretKey := "a16byteslongkey!"
	plaintext := "This is a secret message."

	crypter := NewCrypter(secretKey)

	encrypted, err := crypter.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if len(encrypted) == 0 {
		t.Fatalf("expected encrypted data, got empty result")
	}

	decrypted, err := crypter.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Fatalf("expected decrypted text to be '%s', got '%s'", plaintext, decrypted)
	}
}

func TestDecryptWithInvalidData(t *testing.T) {
	secretKey := "a16byteslongkey!"
	invalidCiphertext := []byte("invalid base64 data")

	crypter := NewCrypter(secretKey)

	_, err := crypter.Decrypt(invalidCiphertext)
	if err == nil {
		t.Fatalf("expected an error when decrypting invalid data, got none")
	}
}

func TestEncryptWithInvalidKey(t *testing.T) {
	invalidKey := "shortkey"
	plaintext := "This is a secret message."

	crypter := NewCrypter(invalidKey)

	_, err := crypter.Encrypt([]byte(plaintext))
	if err == nil {
		t.Fatalf("expected an error when encrypting with an invalid key, got none")
	}
}

func TestDecryptWithIncorrectKey(t *testing.T) {
	secretKey := "a16byteslongkey!"
	incorrectKey := "aDifferent16Key!"
	plaintext := "This is a secret message."

	crypter := NewCrypter(secretKey)
	encrypted, err := crypter.Encrypt([]byte(plaintext))
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	crypterWithIncorrectKey := NewCrypter(incorrectKey)
	_, err = crypterWithIncorrectKey.Decrypt(encrypted)
	if err == nil {
		t.Fatalf("expected an error when decrypting with an incorrect key, got none")
	}
}
