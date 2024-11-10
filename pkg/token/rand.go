package token

import (
	"crypto/rand"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandHash(length int) (string, error) {
	bytes := make([]byte, length)

	// Заполняем массив случайными байтами
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	for i, b := range bytes {
		// Преобразуем байты в индекс символа из charset
		bytes[i] = charset[b%byte(len(charset))]
	}

	return string(bytes), nil
}
