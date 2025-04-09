package app

import (
	"fmt"
	"github.com/cmpxNot29a/shurs/internal/helper"
	"io"
	"net/http"
)

var memStorage map[string]string

func init() {
	memStorage = make(map[string]string)
}

func getUrlfromPost(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Невозможно прочитать тело запроса", http.StatusBadRequest)
		return
	}

	if helper.IsValidURL(string(body)) {
		randomID, err := helper.GenerateRandomBase62(8)
		if err != nil {
			http.Error(w, "Внутренняя огик", http.StatusBadRequest)
			return
		}
		memStorage[string(randomID)] = string(body)

		w.WriteHeader(http.StatusCreated)

		w.Write(fmt.Appendf(nil, "http://%s/%s", r.Host, string(randomID)))
		return
	}
	http.Error(w, "No valid URL", http.StatusBadRequest)

}

func getUrlfromGet(w http.ResponseWriter, r *http.Request) {

	idFromURL := r.URL.Path[1:]

	if helper.IsValidBase62String(idFromURL) {
		url, exists := memStorage[idFromURL]
		if exists {
			http.Redirect(w, r, url, http.StatusTemporaryRedirect)
			return
		}
	}
	http.Error(w, "No valid url id", http.StatusBadRequest)
}

func MethodRouter(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		getUrlfromGet(w, r)
	case http.MethodPost:
		getUrlfromPost(w, r)
	default:
		http.Error(w, "Invalid method", http.StatusBadRequest)
	}
}

func App() {

	mux := http.NewServeMux()
	mux.Handle(`/`, helper.MethodPipe(
		http.HandlerFunc(MethodRouter),
		helper.GetPostOnly()),
	)

	err := http.ListenAndServe(`:8080`, mux)

	if err != nil {
		panic(err)
	}

}
