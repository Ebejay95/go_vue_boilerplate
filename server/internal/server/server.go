package server

import (
	"net"
	"net/http"
	"os"

	"backend-grpc-server/internal/handlers"
	"backend-grpc-server/internal/storage"
	pb "backend-grpc-server/pb"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents the gRPC server configuration
type Server struct {
	grpcServer   *grpc.Server
	wrappedGrpc  *grpcweb.WrappedGrpcServer
	userHandler  *handlers.UserHandler
	notifHandler *handlers.NotificationHandler
}

// NewServer creates a new server instance
func NewServer() *Server {
	// Create stores
	userStore := storage.NewInMemoryUserStore()
	notificationStore := storage.NewInMemoryNotificationStore()

	// Create handlers
	userHandler := handlers.NewUserHandler(userStore)
	notifHandler := handlers.NewNotificationHandler(notificationStore)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	pb.RegisterNotificationServiceServer(grpcServer, notifHandler)

	// Enable reflection for grpcurl
	reflection.Register(grpcServer)

	// Create gRPC-Web wrapper
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			allowedOrigins := []string{
				os.Getenv("BASE_URL"),
				os.Getenv("BACKEND_BASE_URL"),
			}
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					return true
				}
			}
			return false
		}),
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			return true
		}),
	)

	return &Server{
		grpcServer:   grpcServer,
		wrappedGrpc:  wrappedGrpc,
		userHandler:  userHandler,
		notifHandler: notifHandler,
	}
}

// ServeGRPC starts the gRPC server
func (s *Server) ServeGRPC(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	return s.grpcServer.Serve(lis)
}

// CreateHTTPHandler creates the HTTP handler for gRPC-Web
func (s *Server) CreateHTTPHandler() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		// CORS Headers für gRPC-Web
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web, grpc-timeout")
		resp.Header().Set("Access-Control-Expose-Headers", "grpc-status, grpc-message")

		if req.Method == "OPTIONS" {
			return
		}

		if s.wrappedGrpc.IsGrpcWebRequest(req) || s.wrappedGrpc.IsGrpcWebSocketRequest(req) {
			s.wrappedGrpc.ServeHTTP(resp, req)
			return
		}

		// Fallback für andere Requests
		http.NotFound(resp, req)
	}
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	s.grpcServer.GracefulStop()
}