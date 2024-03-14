package main

import (
	"log"
	"net/http"
	"os"

	"github.com/justinas/alice"

	"github.com/dniprodev/learn-switch-to-go-homework/internal/handlers"
	"github.com/dniprodev/learn-switch-to-go-homework/internal/middlewares"
	"github.com/dniprodev/learn-switch-to-go-homework/internal/models/message"
	"github.com/dniprodev/learn-switch-to-go-homework/internal/models/user"
)

func main() {
	uri := os.Getenv("DATABASE_URL")
	usersRepo, err := user.NewRepository(uri)
	if err != nil {
		log.Fatalf("Failed to create the users repository: %s", err)
	}

	defer usersRepo.Close()

	err = usersRepo.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %s", err)
	}

	mongoURI := os.Getenv("MONGO_URI")
	messagesRepo, err := message.NewRepository(mongoURI)
	if err != nil {
		log.Fatalf("Failed to create the messages repository: %s", err)
	}

	messagesRepo.Save(message.Message{Id: "1", Text: "Hello, world!"})
	messsages, err := messagesRepo.FindAll()
	if err != nil {
		log.Println("Find all error")
	}
	log.Println(messsages)

	chain := alice.New(
		middlewares.RecoverHandler,
		middlewares.HttpCallsLoggingHandler,
		middlewares.HttpErrorLoggingHandler,
	)

	http.Handle("/user", chain.Then(handlers.CreateUserHandler(usersRepo.Save)))
	http.Handle("/user/login", chain.Then(handlers.LoginUserHandler(usersRepo.FindByUsername)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
