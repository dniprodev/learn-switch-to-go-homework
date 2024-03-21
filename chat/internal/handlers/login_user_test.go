package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/dniprodev/learn-switch-to-go-homework/user_service/models/user"
)

var correctUser = user.NewUser("myUser", "secretPass123")

func loginUserHandlerTest(name, body string, wantStatus int, wantUser user.User, t *testing.T) {
	req, err := http.NewRequest("POST", "/users/login", strings.NewReader(body))
	if err != nil {
		t.Fatalf("could not create request: %v", err)
	}

	w := httptest.NewRecorder()

	findUser := func(ctx context.Context, name, password string) (string, error) {
		if correctUser.Name == name {
			return correctUser.Password, nil
		}

		return "", nil
	}

	LoginUserHandler(findUser).ServeHTTP(w, req)

	resp := w.Result()
	if resp.StatusCode != wantStatus {
		t.Fatalf("expected status %v; got %v", wantStatus, resp.StatusCode)
	}

	if wantStatus != http.StatusOK {
		return
	}

	// if status is OK, check the response body
	responseBody, _ := io.ReadAll(resp.Body)
	var parsedBody LoginUserResponse
	json.Unmarshal(responseBody, &parsedBody)
	if !strings.Contains(parsedBody.URL, "ws://fancy-chat.io/ws&token=") {
		t.Fatalf("Response body does not contain expected URL - got %v", parsedBody.URL)
	}

	// Check the headers
	if contentType := w.Header().Get(HeaderContentType); contentType != ContentTypeJSON {
		t.Fatalf("handler returned wrong content-type: got %v want %v", contentType, ContentTypeJSON)
	}
	const headerKeyExpiresAfter = "X-Expires-After"
	if expiresAfter := w.Header().Get(headerKeyExpiresAfter); expiresAfter == "" {
		t.Fatalf("missing %s header", headerKeyExpiresAfter)
	}
	const headerKeyRateLimit = "X-Rate-Limit"
	if rateLimit := w.Header().Get(headerKeyRateLimit); rateLimit != "1000" {
		t.Fatalf("handler returned wrong %s: got %v want %v", headerKeyRateLimit, rateLimit, "1000")
	}
}
func TestLoginUserHandler(t *testing.T) {
	tests := []struct {
		name       string
		body       string
		wantStatus int
		wantUser   user.User
	}{
		{
			name:       "Valid Input",
			body:       `{"userName":"` + correctUser.Name + `","password":"` + correctUser.Password + `"}`,
			wantStatus: http.StatusOK,
			wantUser: user.User{
				ID:       "1",
				Name:     "myUser",
				Password: correctUser.Password,
			},
		},
		{
			name:       "Invalid UserName",
			body:       `{"userName": "wrong","password":"` + correctUser.Password + `"}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "Invalid Password",
			body:       `{"userName": "` + correctUser.Name + `","password":"wrongPass"}`,
			wantStatus: http.StatusBadRequest,
		},
		// add more test cases as per your requirements
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			loginUserHandlerTest(tc.name, tc.body, tc.wantStatus, tc.wantUser, t)
		})
	}
}
