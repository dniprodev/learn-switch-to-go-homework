package main

import (
	"log"
	"net/http"

	// TODO: Q - is ok to import in such verbose way?
	"github.com/dniprodev/learn-switch-to-go-homework/pkg/handlers"
	"github.com/dniprodev/learn-switch-to-go-homework/pkg/models/user"
)

// TODO: Q - What is better approach
var userRepository user.Repository

func main() {
	// TODO: Q - Not sure i am correctly implment  Mat Ryer advise
	http.HandleFunc("/user", handlers.CreateUserHandler(&userRepository))
	http.HandleFunc("/user/login", handlers.LoginUserHandler(&userRepository))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// TODO: Q -
// ❯ go test
// ?       github.com/dniprodev/learn-switch-to-go-homework        [no test files]
// but next is ok
// ❯ go test -v ./...

// TODO: Q -
// So many noisy files in root
