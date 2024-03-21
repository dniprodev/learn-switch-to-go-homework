package middlewares

import (
	"log/slog"
	"net/http"
)

func HttpCallsLoggingHandler(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		slog.Info("HTTP Call", "method", r.Method, "url", r.URL.Path)

		next.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}
