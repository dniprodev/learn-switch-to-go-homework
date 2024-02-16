package handlers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

// TODO: Q - does it exist in Go library
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

func LoginUserHandler(repo UserRepositoryInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request loginUserRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		// TODO: Q - it is not required import )))
		user, ok := repo.FindByUsername(request.UserName)
		if !ok {
			http.Error(w, "Invalid username/password", http.StatusBadRequest)
			return
		}
		if user.Password != request.Password {
			http.Error(w, "Invalid username/password", http.StatusBadRequest)
			return
		}

		token, err := generateRandomToken()
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		response := LoginUserResponse{
			URL: fmt.Sprintf("ws://fancy-chat.io/ws&token=%s", token),
		}

		w.Header().Set(HeaderContentType, ContentTypeJSON)
		// TODO: Q - Who control compliance? Server? how it is usually impemented?
		w.Header().Set("X-Expires-After", time.Now().Add(2*time.Hour).Format(time.RFC3339))
		w.Header().Set("X-Rate-Limit", "1000")
		json.NewEncoder(w).Encode(response)
	}
}

func generateRandomToken() (string, error) {
	n := 10 // change to the length of the token you want
	bytes := make([]byte, n)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	for i := 0; i < n; i++ {
		bytes[i] = byte(65 + rand.Intn(85)) // A=65 and Z = 65+25
	}
	return string(bytes), nil
}
