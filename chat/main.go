package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/justinas/alice"

	"github.com/dniprodev/learn-switch-to-go-homework/chat/internal/handlers"
	"github.com/dniprodev/learn-switch-to-go-homework/chat/internal/middlewares"
	"github.com/dniprodev/learn-switch-to-go-homework/chat/internal/models/user"
)

func main() {
	context := context.Background()

	uri := os.Getenv("DATABASE_URL")
	fmt.Println(uri)
	usersRepo, err := user.NewRepository(context, uri)
	if err != nil {
		log.Fatalf("Failed to create the users repository: %s", err)
	}

	defer usersRepo.Close()

	err = usersRepo.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %s", err)
	}

	// mongoURI := os.Getenv("MONGO_URI")
	// messagesRepo, err := message.NewRepository(mongoURI)
	// if err != nil {
	// log.Fatalf("Failed to create the messages repository: %s", err)
	// }

	// messagesRepo.Save(message.Message{Id: "1", Text: "Hello, world!"})
	// messsages, err := messagesRepo.FindAll()
	// if err != nil {
	// 	log.Println("Find all error")
	// }
	// log.Println(messsages)

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

	httpPort := os.Getenv("HTTP_PORT")

	fmt.Println("!!!!!!!!")
	fmt.Println(uri)
	fmt.Println(httpPort)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil))

	// TODOq: testcontainers, ryuk is dsiabled
	// https://medium.com/twodigits/testcontainers-on-podman-a090c348b9d8
	// https://github.com/testcontainers/testcontainers-java/issues/2088
	// https://github.com/containers/podman/blob/main/docs/tutorials/socket_activation.md
	// TODOq: TestMain, TestSutie
	// TODOq: converting existed structure into workspace
	// TODOq: ene-to-end
}
