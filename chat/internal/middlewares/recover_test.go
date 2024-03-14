package middlewares

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type TestResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (trw *TestResponseWriter) WriteHeader(code int) {
	trw.statusCode = code
	trw.ResponseWriter.WriteHeader(code)
}

func TestCreateRecoverHandler(t *testing.T) {
	tests := []struct {
		method      string
		url         string
		panic       bool
		expectedLog string
	}{
		{"GET", "/test1", false, ""},                               // No panic, no log expected
		{"POST", "/test2", true, `"method":"POST","url":"/test2"`}, // Panic, log expected
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
		rr := &TestResponseWriter{ResponseWriter: httptest.NewRecorder()}

		// Create a handler to test
		handler := CreateRecoverHandler(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if tt.panic {
				panic("Test panic")
			}
		}))

		handler.ServeHTTP(rr, req)

		// Check the logger output
		if tt.panic {
			if tt.expectedLog != "" && !strings.Contains(buf.String(), tt.expectedLog) {
				t.Errorf("Logger output does not match expected. Got: %v, Expected to contain: %v", buf.String(), tt.expectedLog)
			}
			if rr.statusCode != http.StatusInternalServerError {
				t.Errorf("Incorrect HTTP status code. Got: %d, Expected: %d", rr.statusCode, http.StatusInternalServerError)
			}
		} else {
			if buf.String() != "" {
				t.Errorf("Did not expect any log output, got: %s", buf.String())
			}
		}
	}
}
