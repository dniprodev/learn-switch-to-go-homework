package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/justinas/alice"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/dniprodev/learn-switch-to-go-homework/chat/internal/handlers"
	"github.com/dniprodev/learn-switch-to-go-homework/chat/internal/middlewares"
	"github.com/dniprodev/learn-switch-to-go-homework/user_service/models/user"

	pb "github.com/dniprodev/learn-switch-to-go-homework/generated"
)

func createGRPCClient(ctx context.Context, address string, opts ...grpc.DialOption) (pb.UserServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %w", err)
	}
	return pb.NewUserServiceClient(conn), conn, nil
}

func storeUser(ctx context.Context, client pb.UserServiceClient, user user.User) (string, error) {
	resp, err := client.StoreUser(ctx, &pb.StoreUserRequest{User: &pb.User{
		UserId:   user.ID,
		Name:     user.Name,
		Password: user.Password,
	}})
	log.Printf("Store user resp: %v", resp)
	if err != nil {
		return "", fmt.Errorf("failed to store user: %v", err)
	}
	return resp.Result, nil
}

func loginUser(ctx context.Context, client pb.UserServiceClient, username string, password string) (string, error) {
	resp, err := client.LoginUser(ctx, &pb.LoginUserRequest{
		UserName: username,
		Password: password,
	})
	if err != nil {
		return "", fmt.Errorf("failed to login user: %v", err)
	}
	return resp.Token, nil
}

func main() {
	userServiceAddress := os.Getenv("USER_SERVICE_ADDRESS")
	mongoURI := os.Getenv("MONGO_URI")
	ctx := context.Background()
	fmt.Println("USER_SERVICE_ADDRESS: ", userServiceAddress)
	fmt.Println("MONGO_URI: ", mongoURI)
	userServiceClient, conn, err := createGRPCClient(ctx, userServiceAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to create user service client: %v", err)
	}
	defer conn.Close()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	chain := alice.New(
		middlewares.HttpCallsLoggingHandler,
		middlewares.CreateRecoverHandler(logger),
		middlewares.CreateHttpErrorLoggingHandler(logger),
	)

	storeUserWrapper := func(ctx context.Context, user user.User) error {
		_, err := storeUser(ctx, userServiceClient, user)
		return err
	}

	loginUserWrapper := func(ctx context.Context, username string, password string) (string, error) {
		return loginUser(ctx, userServiceClient, username, password)
	}

	http.Handle("/user", chain.Then(handlers.CreateUserHandler(storeUserWrapper)))
	http.Handle("/user/login", chain.Then(handlers.LoginUserHandler(loginUserWrapper)))

	httpPort := os.Getenv("HTTP_PORT")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil))

	// TODOq: testcontainers, ryuk is dsiabled
	// https://medium.com/twodigits/testcontainers-on-podman-a090c348b9d8
	// https://github.com/testcontainers/testcontainers-java/issues/2088
	// https://github.com/containers/podman/blob/main/docs/tutorials/socket_activation.md
	// TODOq: TestMain, TestSutie
	// TODOq: converting existed structure into workspace
	// TODOq: ene-to-end
}
