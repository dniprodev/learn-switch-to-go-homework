package middlewares

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateHttpErrorLoggingHandler(t *testing.T) {
	tests := []struct {
		method      string
		url         string
		statusCode  int
		expectedLog string
	}{
		{"GET", "/test1", 200, ""},  // Expected log for this case
		{"POST", "/test2", 201, ""}, // Expected log for this case
		{"PUT", "/test3", 500, "\"level\":\"ERROR\",\"msg\":\"HTTP error occurred\",\"errorCode\":500"},    // Expected log for this case
		{"DELETE", "/test4", 404, "\"level\":\"ERROR\",\"msg\":\"HTTP error occurred\",\"errorCode\":404"}, // Expected log for this case
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Creating the buffer and pass it to JSONHandler.
		var buf bytes.Buffer
		logger := slog.New(slog.NewJSONHandler(&buf, nil))

		// Create a ResponseRecorder to record the response.
		rr := httptest.NewRecorder()

		// Create a handler to test
		handler := CreateHttpErrorLoggingHandler(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(tt.statusCode)
		}))

		handler.ServeHTTP(rr, req)

		if !strings.Contains(buf.String(), tt.expectedLog) {
			t.Errorf("Logger output does not match expected. Got: %v, Expected to contain: %v", buf.String(), tt.expectedLog)
		}
	}
}
