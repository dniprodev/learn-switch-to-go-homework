package handlers

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dniprodev/learn-switch-to-go-homework/chat/internal/models/user"
)

func CreateUserHandlerTest(name, body string, wantStatus int, wantUserName string, t *testing.T) {
	req, err := http.NewRequest("POST", "/users", strings.NewReader(body))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	w := httptest.NewRecorder()

	var userToBeCreated user.User
	createUser := func(user user.User) error {
		userToBeCreated = user
		return nil
	}

	CreateUserHandler(createUser).ServeHTTP(w, req)

	resp := w.Result()

	if resp.StatusCode != wantStatus {
		t.Fatalf("expected status %v; got %v", wantStatus, resp.StatusCode)
	}

	if resp.StatusCode == http.StatusOK {
		responseBody, _ := io.ReadAll(resp.Body)
		if !strings.Contains(string(responseBody), "userName") {
			t.Fatalf("Response body does not contain userName fields")
		}

		if userToBeCreated.Name != wantUserName {
			t.Fatalf("expected user created %v", wantUserName)
		}
	}
}

func TestCreateUserHandler(t *testing.T) {
	tests := []struct {
		name         string
		body         string
		wantStatus   int
		wantUserName string
	}{
		{
			name:         "Valid Input",
			body:         `{"userName": "myUser","password":"secretPass123"}`,
			wantStatus:   http.StatusOK,
			wantUserName: "myUser",
		},
		{
			name:       "Missing Password",
			body:       `{"userName": "myUser"}`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			CreateUserHandlerTest(tc.name, tc.body, tc.wantStatus, tc.wantUserName, t)
		})
	}
}
