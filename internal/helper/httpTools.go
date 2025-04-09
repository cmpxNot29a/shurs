package helper

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
)

// Middleware для валидации URL в теле POST запроса
func ValidateURLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Читаем тело запроса. Важно: чтение "потребляет" тело.
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR: Middleware (URL): Failed to read request body: %v", err)
			// Ошибка чтения тела - скорее внутренняя ошибка сервера или проблема сети
			http.Error(w, "Невозможно прочитать тело запроса", http.StatusInternalServerError)
			return
		}
		// Закрываем оригинальное тело запроса СРАЗУ после чтения
		r.Body.Close()

		originalURL := string(bodyBytes)

		// Выполняем валидацию URL
		if !IsValidURL(originalURL) {
			log.Printf("WARN: Middleware (URL): Invalid URL received: %s", originalURL)
			http.Error(w, "Неверный формат URL", http.StatusBadRequest)
			return // Прерываем цепочку, не вызываем next
		}

		// Валидация прошла.
		// ВОССТАНАВЛИВАЕМ ТЕЛО ЗАПРОСА, чтобы следующий обработчик мог его прочитать.
		// Создаем новый io.ReadCloser из прочитанных байт.
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		// Передаем управление следующему обработчику в цепочке
		log.Println("INFO: Middleware (URL): Validation successful.")
		next.ServeHTTP(w, r)
	})
}

// Middleware для валидации формата ID в пути GET запроса
func ValidateIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		idFromURL := r.PathValue("id")

		if !IsValidBase62String(idFromURL) {
			log.Printf("WARN: Middleware (ID): Invalid ID format received: %s", idFromURL)
			http.Error(w, "Неверный формат ID", http.StatusBadRequest)
			return
		}

		log.Printf("INFO: Middleware (ID): Validation successful for ID: %s", idFromURL)
		next.ServeHTTP(w, r)
	})
}

func IsValidURL(testURL string) bool {
	parsedURL, err := url.Parse(testURL)
	if err != nil {
		return false
	}

	// Проверка на наличие схемы и хоста
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}
