package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpCallsLoggingHandler(t *testing.T) {
	tests := []struct {
		method string
		url    string
	}{
		{"GET", "/test"},
		{"POST", "/test2"},
		{"DELETE", "/test3"},
	}

	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, tt.url, nil)
		if err != nil {
			t.Fatal(err)
		}

		// Create a ResponseRecorder to record the response.
		rr := httptest.NewRecorder()

		// Create a handler to test
		handler := HttpCallsLoggingHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))

		handler.ServeHTTP(rr, req)

		expectedStatus := 200
		if status := rr.Code; status != expectedStatus {
			t.Errorf("handler returned wrong status code: got %v want %v", status, expectedStatus)
		}
	}
}
