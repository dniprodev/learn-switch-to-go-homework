package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/justinas/alice"

	"github.com/dniprodev/learn-switch-to-go-homework/internal/handlers"
	"github.com/dniprodev/learn-switch-to-go-homework/internal/middlewares"
	"github.com/dniprodev/learn-switch-to-go-homework/internal/models/message"
	"github.com/dniprodev/learn-switch-to-go-homework/internal/models/user"
)

func main() {
	context := context.Background()

	uri := os.Getenv("DATABASE_URL")
	usersRepo, err := user.NewRepository(context, uri)
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

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	chain := alice.New(
		middlewares.HttpCallsLoggingHandler,
		middlewares.CreateRecoverHandler(logger),
		middlewares.CreateHttpErrorLoggingHandler(logger),
	)

	// TODOq:
	save := func(user user.User) error {
		return usersRepo.Save(context, user)
	}

	findByUsername := func(name string) (user.User, error) {
		return usersRepo.FindByUsername(context, name)
	}

	http.Handle("/user", chain.Then(handlers.CreateUserHandler(save)))
	http.Handle("/user/login", chain.Then(handlers.LoginUserHandler(findByUsername)))

	log.Fatal(http.ListenAndServe(":8080", nil))

	// TODOq: TestMain, TestSutie
}
