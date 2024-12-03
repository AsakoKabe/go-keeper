package utils

import (
	"math/rand"
	"os"

	"github.com/google/uuid"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandStringRunes Генерация случайной строки длинны n
func RandStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

// GetEnv Получить значение из env или значение по умолчанию
func GetEnv(key, defaultValue string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		return defaultValue
	}
	return value
}

// IsValidUUID Проверяет строку на валидность типа UUID
func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
