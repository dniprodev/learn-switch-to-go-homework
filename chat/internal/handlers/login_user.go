package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	HeaderContentType = "Content-Type"
	ContentTypeJSON   = "application/json"
)

type loginUserRequest struct {
	Password string `json:"password"`
	UserName string `json:"userName"`
}

type LoginUserResponse struct {
	URL string `json:"url"`
}

// TODOq: Is this correct (loginUser)
func LoginUserHandler(loginUser func(context.Context, string, string) (string, error)) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		var request loginUserRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		token, err := loginUser(r.Context(), request.UserName, request.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		response := LoginUserResponse{
			URL: fmt.Sprintf("ws://fancy-chat.io/ws&token=%s", token),
		}

		w.Header().Set(HeaderContentType, ContentTypeJSON)
		w.Header().Set("X-Expires-After", time.Now().Add(2*time.Hour).Format(time.RFC3339))
		w.Header().Set("X-Rate-Limit", "1000")
		json.NewEncoder(w).Encode(response)
	}
	return http.HandlerFunc(fn)
}
