package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	pb "backend-grpc-server/pb"

	"google.golang.org/grpc"
)

type User struct {
	ID    int32  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Age   int32  `json:"age"`
	Role  string `json:"role"`
}

type UserServer struct {
	pb.UnimplementedUserServiceServer
	users  map[int32]*User
	mutex  sync.RWMutex
	nextID int32
}

func NewUserServer() *UserServer {
	server := &UserServer{
		users:  make(map[int32]*User),
		nextID: 1,
	}

	// Add some sample data
	server.users[1] = &User{ID: 1, Name: "Max Mustermann", Email: "max@example.com", Age: 30, Role: "admin"}
	server.users[2] = &User{ID: 2, Name: "Anna Schmidt", Email: "anna@example.com", Age: 25, Role: "user"}
	server.nextID = 3

	return server
}

func (s *UserServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	user, exists := s.users[req.Id]
	if !exists {
		return nil, fmt.Errorf("user with ID %d not found", req.Id)
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
			Role:   user.Role,
		},
	}, nil
}

func (s *UserServer) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	user := &User{
		ID:    s.nextID,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
		Role:  req.Role,
	}

	s.users[s.nextID] = user
	s.nextID++

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
			Role:   user.Role,
		},
	}, nil
}

func (s *UserServer) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	var users []*pb.User
	for _, user := range s.users {
		users = append(users, &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
			Role:  user.Role,
		})
	}

	return &pb.ListUsersResponse{
		Users: users,
	}, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, NewUserServer())

	log.Printf("gRPC server listening on port %s", port)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
