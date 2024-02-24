package main

import (
	"log"
	"net/http"

	"github.com/justinas/alice"

	"github.com/dniprodev/learn-switch-to-go-homework/internal/handlers"
	"github.com/dniprodev/learn-switch-to-go-homework/internal/middlewares"
	"github.com/dniprodev/learn-switch-to-go-homework/internal/models/user"
)

func main() {
	var userRepository user.Repository

	chain := alice.New(
		middlewares.RecoverHandler,
		middlewares.HttpCallsLoggingHandler,
		middlewares.HttpErrorLoggingHandler,
	)

	http.Handle("/user", chain.Then(handlers.CreateUserHandler(userRepository.Save)))
	http.Handle("/user/login", chain.Then(handlers.LoginUserHandler(userRepository.FindByUsername)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
