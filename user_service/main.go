package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"

	usermanagement "github.com/dniprodev/learn-switch-to-go-homework/generated"
	"github.com/dniprodev/learn-switch-to-go-homework/user_service/models/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type server struct {
	usermanagement.UnimplementedUserServiceServer
	usersRepo *user.Repository
}

func (s *server) LoginUser(ctx context.Context, req *usermanagement.LoginUserRequest) (*usermanagement.LoginUserResponse, error) {
	log.Printf("Received login request for user: %v", req.UserName)
	user, err := s.usersRepo.FindByUsername(ctx, req.UserName)
	if err != nil {
		return nil, fmt.Errorf("user %s not found", req.UserName)
	}

	if user.Password != req.Password {
		log.Printf("%v != %v", user.Password, req.Password)
		return nil, fmt.Errorf("invalid password")
	}

	// token, err := generateRandomToken()

	if err != nil {
		return nil, fmt.Errorf("can't generate token")
	}

	return &usermanagement.LoginUserResponse{
		Token: "token",
	}, nil
}

func (s *server) GetUser(ctx context.Context, in *usermanagement.GetUserRequest) (*usermanagement.GetUserResponse, error) {
	user, err := s.usersRepo.FindByUsername(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	log.Printf("Received: %v", in.UserId)
	return &usermanagement.GetUserResponse{
		User: &usermanagement.User{
			UserId: user.ID,
			Name:   user.Name,
		},
	}, nil
}

func (s *server) StoreUser(ctx context.Context, req *usermanagement.StoreUserRequest) (*usermanagement.StoreUserResponse, error) {
	log.Printf("Received store user request for user: %v", req.User.Name)

	user := user.User{
		ID:       req.User.UserId,
		Name:     req.User.Name,
		Password: req.User.Password,
	}

	err := s.usersRepo.Save(ctx, user)
	if err != nil {
		return &usermanagement.StoreUserResponse{Result: "Failed to store user"}, err
	}
	return &usermanagement.StoreUserResponse{Result: "User stored successfully!!!"}, nil
}

// Global variable to hold the image data.
var imageData []byte

func (s *server) UploadImage(stream usermanagement.UserService_UploadImageServer) error {
	for {
		// Receive a message from the stream.
		req, err := stream.Recv()
		if err == io.EOF {
			// The client has finished sending data.
			break
		}
		if err != nil {
			return status.Errorf(codes.Unknown, "failed to receive data: %v", err)
		}

		// Append the received data to the global variable.
		imageData = append(imageData, req.Data...)
	}

	// Send a response to the client.
	res := &usermanagement.UploadImageResponse{
		Result: "Image uploaded successfully",
	}
	if err := stream.SendAndClose(res); err != nil {
		return status.Errorf(codes.Unknown, "failed to send response: %v", err)
	}

	return nil
}

func main() {
	portStr := os.Getenv("GRPC_PORT")
	if portStr == "" {
		portStr = "50051" // default port
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	uri := os.Getenv("DATABASE_URL")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	usersRepo, err := user.NewRepository(ctx, uri)
	if err != nil {
		log.Fatalf("Failed to create the users repository: %s", err)
	}
	defer usersRepo.Close()
	err = usersRepo.Initialize()
	if err != nil {
		log.Fatalf("Failed to initialize the database: %s", err)
	}

	s := grpc.NewServer()
	usermanagement.RegisterUserServiceServer(s, &server{usersRepo: usersRepo})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
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
