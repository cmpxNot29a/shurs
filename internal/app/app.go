package app

import (
	"github.com/cmpxNot29a/shurs/internal/helper"
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

var memStorage map[string]string

func init() {
	memStorage = make(map[string]string)
}

func App() {

	address := ":8080"
	r := chi.NewRouter()

	validatedPostHandler := helper.ValidateURLMiddleware(http.HandlerFunc(handleCreateShortURL))
	validatedGetHandler := helper.ValidateIDMiddleware(http.HandlerFunc(handleRedirect))

	r.Post("/", validatedPostHandler.ServeHTTP)
	r.Get("/{id}", validatedGetHandler.ServeHTTP)

	log.Printf("INFO: Starting server on address %s", address)
	err := http.ListenAndServe(address, r)

	if err != nil {
		log.Fatalf("FATAL: Server failed to start: %v", err)
	}

}
