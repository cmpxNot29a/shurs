package middleware

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

// responseData хранит код статуса и размер ответа.
type responseData struct {
	status int
	size   int
}

// loggingResponseWriter реализует http.ResponseWriter и перехватывает статус и размер.
type loggingResponseWriter struct {
	http.ResponseWriter
	responseData *responseData
}

// WriteHeader перехватывает код статуса.
func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode // Записываем статус
}

// Write перехватывает количество записанных байт.
func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size
	return size, err
}

func LoggingMiddleware(log *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		// Возвращаем сам обработчик middleware
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{}
			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}

			uri := r.RequestURI
			method := r.Method

			// Вызываем следующий обработчик в цепочке с обернутым ResponseWriter
			next.ServeHTTP(&lw, r)

			duration := time.Since(start)

			if responseData.status == 0 {
				responseData.status = http.StatusOK
			}

			log.Info("Request processed",
				zap.String("uri", uri),
				zap.String("method", method),
				zap.Duration("duration", duration),
				zap.Int("status", responseData.status),
				zap.Int("size", responseData.size),
			)
		})
	}
}
