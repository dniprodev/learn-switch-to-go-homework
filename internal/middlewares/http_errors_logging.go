package middlewares

import (
	"log/slog"
	"net/http"
)

// Custom response writer
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Override WriteHeader
func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func HttpErrorLoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		respWriter := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(respWriter, r)

		if respWriter.statusCode >= 400 {
			slog.Error("HTTP error occurred", "errorCode", respWriter.statusCode, "url", r.URL)
		}
	}
	return http.HandlerFunc(fn)
}