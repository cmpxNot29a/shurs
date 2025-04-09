package app

import (
	"github.com/cmpxNot29a/shurs/internal/helper"
	"log"
	"net/http"
)

var memStorage map[string]string

func init() {
	memStorage = make(map[string]string)
}

func App() {

	address := ":8080"
	mux := http.NewServeMux()

	validatedPostHandler := helper.ValidateURLMiddleware(http.HandlerFunc(handleCreateShortURL))
	validatedGetHandler := helper.ValidateIDMiddleware(http.HandlerFunc(handleRedirect))

	mux.Handle("POST /", validatedPostHandler)
	mux.Handle("GET /{id}", validatedGetHandler)

	log.Printf("INFO: Starting server on address %s", address)
	err := http.ListenAndServe(address, mux)

	if err != nil {
		log.Fatalf("FATAL: Server failed to start: %v", err)
	}

}
