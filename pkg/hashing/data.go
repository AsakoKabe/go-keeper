package hashing

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
)

// CryptInterface определяет контракт для шифрования и расшифровки данных.
type CryptInterface interface {
	Encrypt(data []byte) ([]byte, error)

	Decrypt(ciphertext []byte) (string, error)
}

type Crypt struct {
	SecretKey string
}

// NewCrypter создаёт новый объект Crypt
func NewCrypter(secretKey string) *Crypt {
	return &Crypt{SecretKey: secretKey}
}

// Encrypt шифрует данные с использованием AES-GCM.
//
// plaintext: Данные для шифрования.
//
// Возвращает зашифрованные данные в кодировке base64 и ошибку, если таковая возникла.
func (c *Crypt) Encrypt(plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher([]byte(c.SecretKey))
	if err != nil {
		return []byte{}, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return []byte{}, err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return []byte{}, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	b64 := base64.StdEncoding.EncodeToString(ciphertext)

	return []byte(b64), nil
}

// Decrypt расшифровывает данные, зашифрованные с использованием AES-GCM.
//
// ciphertext: Зашифрованные данные в кодировке base64.
//
// Возвращает расшифрованные данные и ошибку, если таковая возникла.
func (c *Crypt) Decrypt(ciphertext []byte) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(string(ciphertext))
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(c.SecretKey))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
