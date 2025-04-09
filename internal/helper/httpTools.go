package helper

import (
	"fmt"
	"net/http"
	"net/url"
	"slices"
)

type Middleware func(http.Handler) http.Handler

func MethodPipe(h http.Handler, middlewares ...Middleware) http.Handler {
	for _, middleware := range middlewares {
		h = middleware(h)
	}
	return h
}

func methodsCheck(method []string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !slices.Contains(method, r.Method) {
				http.Error(w, fmt.Sprintf("Метод %s не поддерживается", r.Method), http.StatusMethodNotAllowed)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}

func IsValidURL(testURL string) bool {
	parsedURL, err := url.Parse(testURL)
	if err != nil {
		return false
	}

	// Проверка на наличие схемы и хоста
	return parsedURL.Scheme != "" && parsedURL.Host != ""
}

func GetOnly() Middleware {
	return methodsCheck([]string{http.MethodGet})
}

func PostOnly() Middleware {
	return methodsCheck([]string{http.MethodPost})
}

func GetPostOnly() Middleware {
	return methodsCheck([]string{http.MethodGet, http.MethodPost})
}
