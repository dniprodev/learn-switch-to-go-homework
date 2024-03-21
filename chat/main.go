package main

import (
	"context"
	"fmt"
	"io"
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

func uploadImage(ctx context.Context, client pb.UserServiceClient, image io.Reader) error {
	// Create a new stream for uploading the image.
	stream, err := client.UploadImage(ctx)
	if err != nil {
		return fmt.Errorf("failed to create upload stream: %v", err)
	}

	// Create a buffer to hold the image chunks.
	buf := make([]byte, 1024)

	// Read the image data in chunks and send each chunk to the server.
	for {
		n, err := image.Read(buf)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read image data: %v", err)
		}
		if n == 0 {
			break
		}

		// Create a new ImageData message with the chunk of image data.
		req := &pb.ImageData{
			Data: buf[:n],
		}

		// Send the ImageData message to the server.
		if err := stream.Send(req); err != nil {
			return fmt.Errorf("failed to send image data: %v", err)
		}
	}

	// Close the stream and receive the response from the server.
	_, err = stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("failed to close stream: %v", err)
	}

	return nil
}

func fetchImage(ctx context.Context, client pb.UserServiceClient, imageID string) (io.Reader, error) {
	// Create a new FetchImageRequest
	req := &pb.FetchImageRequest{
		ImageId: imageID,
	}

	// Call the FetchImage method on the client
	stream, err := client.FetchImage(ctx, req)
	if err != nil {
		return nil, err
	}

	// Create a pipe to read the image data
	pr, pw := io.Pipe()

	// Start a new goroutine to read the image data from the stream and write it to the pipe
	go func() {
		defer pw.Close()

		for {
			data, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("Failed to receive image data: %v", err)
				return
			}

			if _, err := pw.Write(data.Data); err != nil {
				log.Printf("Failed to write image data: %v", err)
				return
			}
		}
	}()

	return pr, nil
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

	uploadImageWrapper := func(ctx context.Context, image io.Reader) error {
		return uploadImage(ctx, userServiceClient, image)
	}

	fetchImageWrapper := func(ctx context.Context) (io.Reader, error) {
		return fetchImage(ctx, userServiceClient, "imageID")
	}

	imageHandler := func() http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.Method {
			case http.MethodGet:
				handlers.FetchImageHandler(fetchImageWrapper).ServeHTTP(w, r)
			case http.MethodPost:
				handlers.ImageUploadHandler(uploadImageWrapper).ServeHTTP(w, r)
			default:
				http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			}
		})
	}

	http.Handle("/user", chain.Then(handlers.CreateUserHandler(storeUserWrapper)))
	http.Handle("/user/login", chain.Then(handlers.LoginUserHandler(loginUserWrapper)))
	http.Handle("/image", chain.Then(imageHandler()))

	httpPort := os.Getenv("HTTP_PORT")

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", httpPort), nil))
}

// TODOq: Better place for gRPC client functions
// TODOq: testcontainers, ryuk is dsiabled
// https://medium.com/twodigits/testcontainers-on-podman-a090c348b9d8
// https://github.com/testcontainers/testcontainers-java/issues/2088
// https://github.com/containers/podman/blob/main/docs/tutorials/socket_activation.md
