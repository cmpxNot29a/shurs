package app

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cmpxNot29a/shurs/internal/helper"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock" // Импорт mock
	"github.com/stretchr/testify/require"
)

// MockShortenerService мок для интерфейса ShortenerUseCase.
type MockShortenerService struct {
	mock.Mock
}

// Убедимся, что мок реализует интерфейс (статическая проверка)
var _ ShortenerUseCase = (*MockShortenerService)(nil)

func (m *MockShortenerService) CreateShortURL(ctx context.Context, originalURL string) (string, error) {
	args := m.Called(ctx, originalURL)
	return args.String(0), args.Error(1)
}

func (m *MockShortenerService) GetOriginalURL(ctx context.Context, id string) (string, error) {
	args := m.Called(ctx, id)
	return args.String(0), args.Error(1)
}

// --- Тесты ---

func TestHandler_CreateShortURL(t *testing.T) {
	testCases := []struct {
		name                string
		method              string
		body                string
		mockReturnID        string
		mockReturnErr       error
		expectedStatus      int
		expectedContentType string
		expectedBodyPrefix  string
		testBaseURL         string
	}{
		// ... (тестовые случаи как раньше) ...
		{
			name:                "Valid URL - Success",
			method:              http.MethodPost,
			body:                "https://yandex.ru",
			mockReturnID:        "aBcDeF12",
			mockReturnErr:       nil,
			expectedStatus:      http.StatusCreated,
			expectedContentType: "text/plain; charset=utf-8",
			expectedBodyPrefix:  "http://test.co/",
			testBaseURL:         "http://test.co",
		},
		{
			name:                "Service Error on Create",
			method:              http.MethodPost,
			body:                "https://google.com",
			mockReturnID:        "",
			mockReturnErr:       errors.New("failed to generate unique ID"),
			expectedStatus:      http.StatusInternalServerError,
			expectedContentType: "text/plain; charset=utf-8",
			expectedBodyPrefix:  "",
			testBaseURL:         "http://test.co",
		},
		{
			name:           "Invalid URL (Handled by Middleware)",
			method:         http.MethodPost,
			body:           "invalid-url",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty Body (Handled by Middleware)",
			method:         http.MethodPost,
			body:           "",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Создаем мок и присваиваем интерфейсной переменной
			mockServicePtr := new(MockShortenerService)       // Указатель на мок
			var mockService ShortenerUseCase = mockServicePtr // Присваиваем интерфейсу

			// 2. Создаем хендлер, передавая интерфейс
			handler := NewHandler(mockService, tc.testBaseURL)

			// 3. Настраиваем ожидания мока (используя указатель на мок)
			if tc.expectedStatus != http.StatusBadRequest {
				mockServicePtr.On("CreateShortURL", mock.Anything, tc.body).
					Return(tc.mockReturnID, tc.mockReturnErr).
					Once()
			}

			// 4. Создаем запрос и рекордер
			requestBody := strings.NewReader(tc.body)
			req := httptest.NewRequest(tc.method, "/", requestBody)
			rr := httptest.NewRecorder()

			// 5. Вызываем middleware + handler
			handlerWithMiddleware := ValidateURLMiddleware(http.HandlerFunc(handler.CreateShortURL))
			handlerWithMiddleware.ServeHTTP(rr, req)

			// 6. Проверяем HTTP ответ
			result := rr.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.expectedStatus, result.StatusCode, "Неверный статус код")
			if tc.expectedStatus != http.StatusBadRequest {
				assert.Contains(t, result.Header.Get("Content-Type"), "text/plain", "Неверный Content-Type")
			}
			if tc.expectedStatus == http.StatusCreated {
				responseBodyBytes, err := io.ReadAll(result.Body)
				require.NoError(t, err, "Не удалось прочитать тело ответа")
				expectedBody := tc.testBaseURL + "/" + tc.mockReturnID
				assert.Equal(t, expectedBody, string(responseBodyBytes), "Неверное тело ответа")
			}

			// 7. Проверяем вызовы мока (используя указатель на мок)
			if tc.expectedStatus != http.StatusBadRequest {
				mockServicePtr.AssertExpectations(t)
			} else {
				mockServicePtr.AssertNotCalled(t, "CreateShortURL", mock.Anything, tc.body)
			}
		})
	}
}

func TestHandler_Redirect(t *testing.T) {
	const testURL = "https://yandex.ru"
	const validID = "abcdef12"
	const notFoundID = "notexist"
	const invalidFormatID = "invalid!"
	const serviceErrorID = "dberror1"
	const expectedIDLength = 8 // Длина для middleware

	testCases := []struct {
		name             string
		requestID        string
		mockReturnURL    string
		mockReturnErr    error
		expectedStatus   int
		expectedLocation string
	}{
		// ... (тестовые случаи как раньше) ...
		{
			name:             "Valid ID - Success",
			requestID:        validID,
			mockReturnURL:    testURL,
			mockReturnErr:    nil,
			expectedStatus:   http.StatusTemporaryRedirect,
			expectedLocation: testURL,
		},
		{
			name:             "ID Not Found (Service returns ErrNotFound)",
			requestID:        notFoundID,
			mockReturnURL:    "",
			mockReturnErr:    ErrNotFound, // Используем нашу ошибку
			expectedStatus:   http.StatusNotFound,
			expectedLocation: "",
		},
		{
			name:           "Invalid ID Format (Handled by Middleware)",
			requestID:      invalidFormatID,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Empty ID in path (Handled by Middleware)",
			requestID:      "",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:             "Internal Service Error",
			requestID:        serviceErrorID,
			mockReturnURL:    "",
			mockReturnErr:    errors.New("database connection failed"),
			expectedStatus:   http.StatusInternalServerError,
			expectedLocation: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1. Создаем мок и присваиваем интерфейсу
			mockServicePtr := new(MockShortenerService)
			var mockService ShortenerUseCase = mockServicePtr

			// 2. Создаем хендлер
			handler := NewHandler(mockService, "http://dummy.base")

			// 3. Настраиваем ожидания мока
			if tc.expectedStatus != http.StatusBadRequest {
				mockServicePtr.On("GetOriginalURL", mock.Anything, tc.requestID).
					Return(tc.mockReturnURL, tc.mockReturnErr).
					Once()
			}

			// 4. Создаем запрос, рекордер, контекст chi
			req := httptest.NewRequest(http.MethodGet, "/"+tc.requestID, nil)
			rr := httptest.NewRecorder()
			routeCtx := chi.NewRouteContext()
			routeCtx.URLParams.Add("id", tc.requestID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))

			// 5. Вызываем middleware + handler
			middlewareFunc := ValidateIDMiddleware(expectedIDLength)
			handlerWithMiddleware := middlewareFunc(http.HandlerFunc(handler.Redirect))
			handlerWithMiddleware.ServeHTTP(rr, req)

			// 6. Проверяем ответ
			result := rr.Result()
			defer result.Body.Close()
			assert.Equal(t, tc.expectedStatus, result.StatusCode, "Неверный статус код")
			locationHeader := result.Header.Get("Location")
			assert.Equal(t, tc.expectedLocation, locationHeader, "Неверный заголовок Location")

			// 7. Проверяем вызовы мока
			if tc.expectedStatus != http.StatusBadRequest {
				mockServicePtr.AssertExpectations(t)
			} else {
				mockServicePtr.AssertNotCalled(t, "GetOriginalURL", mock.Anything, tc.requestID)
			}
		})
	}
}

// Тест для хелпера IsValidBase62String
func TestIsValidBase62String(t *testing.T) {
	const testLength = 8 // Длина, используемая в приложении
	testCases := []struct {
		name     string
		input    string
		length   int
		expected bool
	}{
		{"Valid 8 chars", "Abc123Xy", testLength, true},
		{"Too short", "Abc12", testLength, false},
		{"Too long", "Abc123Xyz", testLength, false},
		{"Invalid chars", "Abc123X!", testLength, false},
		{"Empty string", "", testLength, false},
		{"Valid 5 chars (wrong length)", "abc12", 5, true}, // Доп. тест на другую длину
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Вызываем хелпер с нужной длиной
			actual := helper.IsValidBase62String(tc.input, tc.length)
			assert.Equal(t, tc.expected, actual)
		})
	}
}
