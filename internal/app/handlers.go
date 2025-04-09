package app

import (
	"fmt"
	"github.com/cmpxNot29a/shurs/internal/helper"
	"io"
	"log"
	"net/http"
)

func handleCreateShortURL(w http.ResponseWriter, r *http.Request) {

	const idLength = 8
	const attempts = 10

	defer r.Body.Close()

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Handler: Failed to read request body after middleware: %v", err)
		http.Error(w, "Внутренняя ошибка сервера", http.StatusInternalServerError)
		return
	}
	originalURL := string(bodyBytes)

	randomID, err := helper.GenUnicID(attempts, idLength, memStorage)
	if err != nil {
		log.Printf("ERROR: Handler: Failed to generate unique ID: %v", err)
		http.Error(w, "Внутренняя ошибка сервера при генерации ID", http.StatusInternalServerError)
		return
	}

	memStorage[randomID] = originalURL

	log.Printf("INFO: Handler: Stored URL: %s -> %s", randomID, originalURL)

	baseURL := fmt.Sprintf("http://%s", r.Host) // TODO: Config
	shortURL := fmt.Sprintf("%s/%s", baseURL, randomID)

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
}

func handleRedirect(w http.ResponseWriter, r *http.Request) {

	idFromURL := r.PathValue("id")

	originalURL, exists := memStorage[idFromURL]

	if exists {
		log.Printf("INFO: Handler: Redirecting: %s -> %s", idFromURL, originalURL)
		http.Redirect(w, r, originalURL, http.StatusTemporaryRedirect)
		return
	}

	log.Printf("WARN: Handler: ID not found: %s", idFromURL)
	http.Error(w, "URL не найден", http.StatusBadRequest)
}
