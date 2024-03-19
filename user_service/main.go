package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"

	pb "github.com/dniprodev/learn-switch-to-go-homework/generated"
	"github.com/dniprodev/learn-switch-to-go-homework/user_service/models/user"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedUserServiceServer
	usersRepo *user.Repository
}

func (s *server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	log.Printf("Received login request for user: %v", req.UserName)
	user, err := s.usersRepo.FindByUsername(ctx, req.UserName)
	if err != nil {
		return nil, err
	}

	if user.Password != req.Password {
		return nil, fmt.Errorf("invalid username or password")
	}

	log.Printf("Received login request for user: %v", req.UserName)

	token, err := generateRandomToken()

	if err != nil {
		return nil, fmt.Errorf("can't generate token")
	}

	return &pb.LoginUserResponse{
		Token: token,
	}, nil
}

func (s *server) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.usersRepo.FindByUsername(ctx, in.UserId)
	if err != nil {
		return nil, err
	}
	log.Printf("Received: %v", in.UserId)
	return &pb.GetUserResponse{
		User: &pb.User{
			UserId: user.ID,
			Name:   user.Name,
		},
	}, nil
}

func (s *server) StoreUser(ctx context.Context, req *pb.StoreUserRequest) (*pb.StoreUserResponse, error) {
	user := user.User{
		ID:       req.User.UserId,
		Name:     req.User.Name,
		Password: req.User.Password,
	}

	err := s.usersRepo.Save(ctx, user)
	if err != nil {
		return &pb.StoreUserResponse{Result: "Failed to store user"}, err
	}
	return &pb.StoreUserResponse{Result: "User stored successfully"}, nil
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
	pb.RegisterUserServiceServer(s, &server{usersRepo: usersRepo})
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
