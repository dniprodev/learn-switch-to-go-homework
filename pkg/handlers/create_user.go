package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/dniprodev/learn-switch-to-go-homework/pkg/models/user"
)

type createUserRequest struct {
	Password string `json:"password"`
	UserName string `json:"userName"`
}

type CreateUserResponse struct {
	ID       string `json:"id"`
	UserName string `json:"userName"`
}

func (request createUserRequest) validate() error {
	if request.UserName == "" {
		return errors.New("username is required")
	}
	if len(request.UserName) < 4 {
		return errors.New("username must be at least 4 characters long")
	}
	if request.Password == "" {
		return errors.New("password is required")
	}
	if len(request.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

func CreateUserHandler(repo UserRepositoryInterface) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request createUserRequest
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		err = request.validate()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		newUser := user.User{
			Name:     request.UserName,
			Password: request.Password,
		}

		newUser = repo.Save(newUser)

		response := CreateUserResponse{
			ID:       newUser.ID,
			UserName: newUser.Name,
		}

		w.Header().Set(HeaderContentType, ContentTypeJSON)
		json.NewEncoder(w).Encode(response)
	}
}
