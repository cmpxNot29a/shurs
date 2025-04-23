package handler

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	repo "github.com/cmpxNot29a/shurs/internal/domain/repository"
	domainUC "github.com/cmpxNot29a/shurs/internal/domain/usecase"
	"github.com/go-chi/chi/v5"
)

// Handler обрабатывает HTTP запросы.
type Handler struct {
	useCase domainUC.ShortenerUseCase // Зависимость от интерфейса из domain
	baseURL string
}

// NewHandler создает новый HTTP хендлер.
func NewHandler(uc domainUC.ShortenerUseCase, baseURL string) *Handler {
	return &Handler{useCase: uc, baseURL: baseURL}
}

// Create обрабатывает POST / запросы.
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR (Handler): Failed to read request body: %v", err)
		http.Error(w, "Failed to read request", http.StatusInternalServerError)
		return
	}
	originalURL := string(bodyBytes)

	shortID, err := h.useCase.CreateShortURL(r.Context(), originalURL)
	if err != nil {
		log.Printf("ERROR (Handler): UseCase failed to create short URL: %v", err)
		http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
		return
	}

	shortURL := fmt.Sprintf("%s/%s", h.baseURL, shortID)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

// Redirect обрабатывает GET /{id} запросы.
func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	shortID := chi.URLParam(r, "id")

	originalURL, err := h.useCase.GetOriginalURL(r.Context(), shortID)
	if err != nil {
		if errors.Is(err, repo.ErrNotFound) {
			log.Printf("WARN (Handler): ID not found: %s", shortID)
			http.Error(w, "URL not found", http.StatusNotFound)
		} else {
			log.Printf("ERROR (Handler): UseCase failed to get original URL: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
}
