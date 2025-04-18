package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Handler обрабатывает HTTP запросы, делегируя логику сервису.
type Handler struct {
	service ShortenerUseCase
	baseURL string
}

// NewHandler создает новый экземпляр Handler.
func NewHandler(service ShortenerUseCase, baseURL string) *Handler {
	return &Handler{
		service: service,
		baseURL: baseURL,
	}
}

// CreateShortURL обрабатывает POST /
func (h *Handler) CreateShortURL(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Handler: Failed to read request body: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	originalURL := string(bodyBytes)

	shortID, err := h.service.CreateShortURL(r.Context(), originalURL)
	if err != nil {
		log.Printf("ERROR: Handler: Service failed to create short URL: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.baseURL, shortID)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

// Redirect обрабатывает GET /{id}
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "id")

	originalURL, err := h.service.GetOriginalURL(r.Context(), shortID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			log.Printf("WARN: Handler: ID not found: %s", shortID)
			http.Error(w, "URL not found", http.StatusNotFound)
		} else {
			log.Printf("ERROR: Handler: Service failed to get original URL: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
