package helper

import (
	"crypto/rand"
	"fmt"
	"log"
	"regexp"
)

func GenerateRandomBase62(length int) ([]byte, error) {

	const base62Charset = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	// Гарантируем, что у нас достаточно случайных байт для преобразования
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return []byte{}, err
	}

	// Преобразуем каждый байт в соответствующий символ из base62Charset
	result := make([]byte, length)
	charsetLen := byte(len(base62Charset))
	for i, b := range bytes {
		result[i] = base62Charset[b%charsetLen]
	}

	return result, nil
}

func IsValidBase62String(s string) bool {
	// Проверка длины
	if len(s) != 8 {
		return false
	}
	matched, err := regexp.MatchString("^[0-9a-zA-Z]{8}$", s)
	if err != nil {
		return false
	}
	return matched
}

func GenUnicID(attempts int, idLength int, storage map[string]string) (string, error) {
	for range attempts {

		randomIDBytes, err := GenerateRandomBase62(idLength)

		if err != nil {
			continue
		}

		randomID := string(randomIDBytes)
		_, exists := storage[randomID]
		if !exists {
			return randomID, nil
		}
		log.Printf("WARN: Handler: Collision detected for ID: %s. Retrying...", randomID)
	}
	return "", fmt.Errorf("failed to generate unique ID after multiple attempts")
}
