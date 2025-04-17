package helper

import (
	"crypto/rand"
	"fmt"
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

func IsValidBase62String(s string, expectedLength int) bool {
	// Проверка длины
	if len(s) != expectedLength {
		return false
	}

	pattern := fmt.Sprintf("^[0-9a-zA-Z]{%d}$", expectedLength)
	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}
	return matched
}
