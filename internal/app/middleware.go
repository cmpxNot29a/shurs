package app

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/cmpxNot29a/shurs/internal/helper"
	"github.com/go-chi/chi/v5"
)

const defaultMiddlewareIDLengthFallback = 8

// ValidateURLMiddleware проверяет URL в теле POST запроса.
func ValidateURLMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("ERROR: Middleware (URL): Failed to read request body: %v", err)
			http.Error(w, "Cannot read request body", http.StatusInternalServerError)
			return
		}
		r.Body.Close()
		originalURL := string(bodyBytes)

		if !helper.IsValidURL(originalURL) {
			log.Printf("WARN: Middleware (URL): Invalid URL received: %s", originalURL)
			http.Error(w, "Invalid URL format", http.StatusBadRequest)
			return
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		log.Println("INFO: Middleware (URL): Validation successful.")
		next.ServeHTTP(w, r)
	})
}

// ValidateIDMiddleware создает middleware, которое проверяет ID указанной длины.
func ValidateIDMiddleware(expectedIDLength int) func(http.Handler) http.Handler {
	localExpectedLength := expectedIDLength
	if localExpectedLength <= 0 {
		log.Printf("WARN (Middleware Factory): Invalid expected ID length (%d) passed, using fallback %d for validation",
			expectedIDLength, defaultMiddlewareIDLengthFallback)
		localExpectedLength = defaultMiddlewareIDLengthFallback
	}

	middlewareFunc := func(next http.Handler) http.Handler {
		requestHandlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idFromURL := chi.URLParam(r, "id")
			if !helper.IsValidBase62String(idFromURL, localExpectedLength) {
				log.Printf("WARN: Middleware (ID): Invalid ID format received (len != %d or invalid chars): %s",
					localExpectedLength, idFromURL)
				http.Error(w, "Invalid ID format", http.StatusBadRequest)
				return
			}

			log.Printf("INFO: Middleware (ID): Validation successful for ID: %s", idFromURL)
			next.ServeHTTP(w, r)
		})
		return requestHandlerFunc
	}
	return middlewareFunc
}
