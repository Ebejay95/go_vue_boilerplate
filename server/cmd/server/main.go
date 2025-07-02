package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"sync"

	pb "backend-grpc-server/pb"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
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

	log.Printf("📋 GetUser: Found user %s (ID: %d)", user.Name, user.ID)

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
			Role:  user.Role,
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

	log.Printf("➕ CreateUser: Created user %s (ID: %d)", user.Name, user.ID)

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
			Role:  user.Role,
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

	log.Printf("📋 ListUsers: Returning %d users", len(users))

	return &pb.ListUsersResponse{
		Users: users,
	}, nil
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "50051"
	}

	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		webPort = "8081"
	}

	// Standard gRPC Server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	userServer := NewUserServer()
	pb.RegisterUserServiceServer(grpcServer, userServer)

	// gRPC Reflection für grpcurl
	reflection.Register(grpcServer)

	// gRPC-Web Wrapper
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			// Erlaube alle Origins für Development
			// In Production solltest du spezifische Origins definieren
			return true
		}),
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			return true
		}),
	)

	// HTTP Handler für gRPC-Web
	httpHandler := func(resp http.ResponseWriter, req *http.Request) {
		// CORS Headers für gRPC-Web
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web, grpc-timeout")
		resp.Header().Set("Access-Control-Expose-Headers", "grpc-status, grpc-message")

		if req.Method == "OPTIONS" {
			return
		}

		if wrappedGrpc.IsGrpcWebRequest(req) || wrappedGrpc.IsGrpcWebSocketRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}

		// Fallback für andere Requests
		http.NotFound(resp, req)
	}

	// Standard gRPC Server starten
	go func() {
		log.Printf("🚀 gRPC server listening on port %s", port)
		log.Printf("📡 gRPC reflection enabled - you can use grpcurl")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// gRPC-Web Server starten
	log.Printf("🌐 gRPC-Web server listening on port %s", webPort)
	log.Printf("🔗 Frontend can connect to: http://localhost:%s", webPort)

	if err := http.ListenAndServe(":"+webPort, http.HandlerFunc(httpHandler)); err != nil {
		log.Fatalf("Failed to serve gRPC-Web: %v", err)
	}
}
