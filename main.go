package main

import (
	"log"
	"net/http"

	"github.com/dniprodev/learn-switch-to-go-homework/internal/handlers"
	"github.com/dniprodev/learn-switch-to-go-homework/internal/models/user"
)

func main() {
	var userRepository user.Repository

	http.HandleFunc("/user", handlers.CreateUserHandler(userRepository.Save))
	http.HandleFunc("/user/login", handlers.LoginUserHandler(userRepository.FindByUsername))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
