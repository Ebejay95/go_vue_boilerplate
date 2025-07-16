package server

import (
	"log"
	"net"
	"net/http"
	"os"

	"backend-grpc-server/internal/database"
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
	db           *database.DB
}

// NewServer creates a new server instance with PostgreSQL
func NewServer() *Server {
	// Initialize database connection
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create stores with PostgreSQL backend
	userStore := storage.NewPostgresUserStore(db)
	notificationStore := storage.NewPostgresNotificationStore(db)

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
				"http://localhost:3000", // Default frontend
				"http://localhost:8080", // Default alternative
			}
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					return true
				}
			}
			// Allow localhost for development
			return origin == "" || origin == "null"
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
		db:           db,
	}
}

// ServeGRPC starts the gRPC server
func (s *Server) ServeGRPC(port string) error {
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		return err
	}

	log.Printf("gRPC server listening on port %s", port)
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

		// Health check endpoint
		if req.URL.Path == "/health" {
			s.handleHealthCheck(resp, req)
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

// handleHealthCheck provides a health check endpoint
func (s *Server) handleHealthCheck(resp http.ResponseWriter, req *http.Request) {
	if err := s.db.HealthCheck(); err != nil {
		resp.WriteHeader(http.StatusServiceUnavailable)
		resp.Write([]byte(`{"status":"unhealthy","database":"disconnected"}`))
		return
	}

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(`{"status":"healthy","database":"connected"}`))
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	log.Println("Shutting down server...")

	// Stop gRPC server
	s.grpcServer.GracefulStop()

	// Close database connection
	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Println("Server shutdown complete")
}
