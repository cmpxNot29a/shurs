package app

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/cmpxNot29a/shurs/internal/helper"
)

func TestHandleCreateShortURL(t *testing.T) {
	testCases := []struct {
		name                string
		method              string
		body                string
		expectedStatus      int
		expectedContentType string
		expectInStorage     bool
		hostHeader          string
		testBaseURL         string
	}{
		{
			name:                "Valid URL - POST",
			method:              http.MethodPost,
			body:                "https://yandex.ru",
			expectedStatus:      http.StatusCreated,
			expectedContentType: "text/plain; charset=utf-8",
			expectInStorage:     true,
			hostHeader:          "localhost:8080",
			testBaseURL:         "http://testbaseurl.com",
		},
		{
			name:                "Empty Body - POST",
			method:              http.MethodPost,
			body:                "",
			expectedStatus:      http.StatusBadRequest,
			expectedContentType: "text/plain; charset=utf-8",
			expectInStorage:     false,
			hostHeader:          "localhost:8080",
			testBaseURL:         "http://testbaseurl.com",
		},
		{
			name:                "Invalid URL - POST",
			method:              http.MethodPost,
			body:                "invalid-url",
			expectedStatus:      http.StatusBadRequest,
			expectedContentType: "text/plain; charset=utf-8",
			expectInStorage:     false,
			hostHeader:          "localhost:8080",
			testBaseURL:         "http://testbaseurl.com",
		},
	}

	originalBaseURL := currentBaseURL

	t.Cleanup(func() {
		currentBaseURL = originalBaseURL
	})

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			memStorage = make(map[string]string)

			currentBaseURL = tc.testBaseURL

			requestBody := strings.NewReader(tc.body)
			req := httptest.NewRequest(tc.method, "/", requestBody)
			rr := httptest.NewRecorder()

			handlerWithMiddleware := helper.ValidateURLMiddleware(http.HandlerFunc(handleCreateShortURL))
			handlerWithMiddleware.ServeHTTP(rr, req)

			result := rr.Result()
			defer result.Body.Close()

			assert.Equal(t, tc.expectedStatus, result.StatusCode, "Неверный статус код")

			if tc.expectedStatus < 400 {
				assert.Equal(t, tc.expectedContentType, result.Header.Get("Content-Type"), "Неверный Content-Type")
			} else {
				assert.Contains(t, result.Header.Get("Content-Type"), "text/plain", "Неверный Content-Type для ошибки")
			}

			if tc.expectInStorage {
				responseBodyBytes, err := io.ReadAll(result.Body)
				require.NoError(t, err, "Не удалось прочитать тело ответа")
				responseBody := string(responseBodyBytes)
				expectedPrefix := tc.testBaseURL + "/"
				assert.True(t, strings.HasPrefix(responseBody, expectedPrefix), "Тело ответа не начинается с ожидаемого префикса")
				shortID := strings.TrimPrefix(responseBody, expectedPrefix)
				assert.NotEmpty(t, shortID, "Короткий ID в ответе не должен быть пустым")
				storedURL, exists := memStorage[shortID]
				assert.True(t, exists, "URL не был сохранен в хранилище")
				assert.Equal(t, tc.body, storedURL, "В хранилище сохранен неверный URL")
			} else {

				assert.Empty(t, memStorage, "Хранилище должно быть пустым при ошибке")
			}
		})
	}
}

func TestHandleRedirect(t *testing.T) {
	const testID = "abcdef12"
	const testURL = "https://yandex.ru"
	const invalidFormatID = "invalid!"
	const notFoundID = "notfound"

	testCases := []struct {
		name             string
		requestID        string
		setupStorage     bool
		expectedStatus   int
		expectedLocation string
	}{
		{
			name:             "Valid ID - Redirect",
			requestID:        testID,
			setupStorage:     true,
			expectedStatus:   http.StatusTemporaryRedirect, // 307
			expectedLocation: testURL,
		},
		{
			name:             "ID Not Found",
			requestID:        notFoundID,
			setupStorage:     true,                  // Добавляем testID, но запрашиваем notFoundID
			expectedStatus:   http.StatusBadRequest, // 400
			expectedLocation: "",
		},
		{
			name:             "Invalid ID Format",
			requestID:        invalidFormatID,
			setupStorage:     false,                 // Не важно, что в хранилище, формат неверен
			expectedStatus:   http.StatusBadRequest, // 400
			expectedLocation: "",
		},
		{
			name:             "Empty ID in path (handled by PathValue)", // r.PathValue вернет пустую строку
			requestID:        "",
			setupStorage:     false,
			expectedStatus:   http.StatusBadRequest, // 400 (т.к. IsValidBase62String вернет false для "")
			expectedLocation: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			memStorage = make(map[string]string)
			if tc.setupStorage {
				memStorage[testID] = testURL
			}

			req := httptest.NewRequest(http.MethodGet, "/"+tc.requestID, nil)

			req.SetPathValue("id", tc.requestID)

			rr := httptest.NewRecorder()

			handlerWithMiddleware := helper.ValidateIDMiddleware(http.HandlerFunc(handleRedirect))
			handlerWithMiddleware.ServeHTTP(rr, req)

			result := rr.Result()
			defer result.Body.Close()

			assert.Equal(t, tc.expectedStatus, result.StatusCode, "Неверный статус код")

			locationHeader := result.Header.Get("Location")
			assert.Equal(t, tc.expectedLocation, locationHeader, "Неверный заголовок Location")
		})
	}
}

func TestIsValidBase62String(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Valid 8 chars", "Abc123Xy", true},
		{"Too short", "Abc12", false},
		{"Too long", "Abc123Xyz", false},
		{"Invalid chars", "Abc123X!", false},
		{"Empty string", "", false},
		{"Only numbers", "12345678", true},
		{"Only letters", "aBcDeFgH", true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := helper.IsValidBase62String(tc.input)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
