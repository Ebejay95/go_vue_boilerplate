package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"backend-grpc-server/internal/database"
	"backend-grpc-server/internal/handlers"
	"backend-grpc-server/internal/models"
	"backend-grpc-server/internal/storage"
	pb "backend-grpc-server/pb"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// Server represents the server configuration
type Server struct {
	grpcServer          *grpc.Server
	wrappedGrpc         *grpcweb.WrappedGrpcServer
	userHandler         *handlers.UserHandler
	notificationHandler *handlers.NotificationHandler
	socketHandler       *handlers.SocketHandler
	db                  *database.DB
}

// NewServer creates a new server instance
func NewServer() *Server {
	// Initialize database connection
	db, err := database.NewConnection()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Create stores
	userStore := storage.NewPostgresUserStore(db)
	notificationStore := storage.NewPostgresNotificationStore(db)

	// Create socket handler
	socketHandler := handlers.NewSocketHandler()

	// Create handlers
	userHandler := handlers.NewUserHandler(userStore, socketHandler)
	notificationHandler := handlers.NewNotificationHandler(notificationStore, socketHandler)

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	pb.RegisterUserServiceServer(grpcServer, userHandler)
	pb.RegisterNotificationServiceServer(grpcServer, notificationHandler)

	// Enable reflection for grpcurl
	reflection.Register(grpcServer)

	// Create gRPC-Web wrapper
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool {
			allowedOrigins := []string{
				os.Getenv("BASE_URL"),
				os.Getenv("BACKEND_BASE_URL"),
				"http://localhost:3000",
				"http://localhost:8080",
			}
			for _, allowed := range allowedOrigins {
				if origin == allowed {
					return true
				}
			}
			return origin == "" || origin == "null"
		}),
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketOriginFunc(func(req *http.Request) bool {
			return true
		}),
	)

	return &Server{
		grpcServer:          grpcServer,
		wrappedGrpc:         wrappedGrpc,
		userHandler:         userHandler,
		notificationHandler: notificationHandler,
		socketHandler:       socketHandler,
		db:                  db,
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
		// CORS Headers for gRPC-Web
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

		// Fallback for other requests
		http.NotFound(resp, req)
	}
}

// NewSocketHandler creates a WebSocket handler for a specific endpoint
func (s *Server) NewSocketHandler() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		s.socketHandler.ServeSocket(resp, req)
	}
}

// handleHealthCheck provides a health check endpoint
func (s *Server) handleHealthCheck(resp http.ResponseWriter, req *http.Request) {
	if err := s.db.HealthCheck(); err != nil {
		resp.WriteHeader(http.StatusServiceUnavailable)
		resp.Write([]byte(`{"status":"unhealthy","database":"disconnected"}`))
		return
	}

	// Get socket statistics
	clients := s.socketHandler.GetConnectedClients()
	clientCount := len(clients)

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)
	resp.Write([]byte(fmt.Sprintf(`{
		"status":"healthy",
		"database":"connected",
		"socket_clients":%d,
		"timestamp":"%s"
	}`, clientCount, time.Now().Format(time.RFC3339))))
}

// handleSocketStatus provides socket connection information
func (s *Server) handleSocketStatus(resp http.ResponseWriter, req *http.Request) {
	clients := s.socketHandler.GetConnectedClients()

	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusOK)

	statusData := map[string]interface{}{
		"connected_clients": len(clients),
		"clients":           clients,
		"endpoints": map[string]string{
			"notifications": "/notifications",
		},
	}

	json.NewEncoder(resp).Encode(statusData)
}

// GetNotificationHandler returns the notification handler for use in other parts of the app
func (s *Server) GetNotificationHandler() *handlers.NotificationHandler {
	return s.notificationHandler
}

// GetSocketHandler returns the socket handler for use in other parts of the app
func (s *Server) GetSocketHandler() *handlers.SocketHandler {
	return s.socketHandler
}

// Shutdown gracefully shuts down the server
func (s *Server) Shutdown() {
	log.Println("Shutting down server...")

	// Stop socket handler
	s.socketHandler.Shutdown()

	// Stop gRPC server
	s.grpcServer.GracefulStop()

	// Close database connection
	if err := s.db.Close(); err != nil {
		log.Printf("Error closing database connection: %v", err)
	}

	log.Println("Server shutdown complete")
}

// Example usage in other parts of your application:

// Background Jobs / Scheduled Tasks
func (s *Server) StartBackgroundJobs() {
	// Example: Daily maintenance notification
	go func() {
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()

		for range ticker.C {
			s.notificationHandler.NotifyAll(
				"Daily system maintenance completed successfully",
				"info",
				false, // not persistent, just info
			)
		}
	}()

	// Example: Monitor system health
	go func() {
		ticker := time.NewTicker(5 * time.Minute)
		defer ticker.Stop()

		for range ticker.C {
			if err := s.db.HealthCheck(); err != nil {
				// Notify admins about database issues
				s.notificationHandler.NotifyAll(
					"Database connection issue detected",
					"error",
					true, // persistent
				)
			}
		}
	}()
}

// Business Logic Integration Examples

// Example: User Registration Handler
func (s *Server) OnUserRegistered(user *models.User) {
	// Welcome the new user
	s.notificationHandler.NotifyUser(
		user.ID,
		fmt.Sprintf("Welcome %s! Your account has been created successfully.", user.Name),
		"success",
		true, // persistent welcome message
	)

	// Notify admins about new user (real-time only)
	s.notificationHandler.NotifyAll(
		fmt.Sprintf("New user registered: %s", user.Name),
		"info",
		false,
	)
}

// Example: Payment Processing
func (s *Server) OnPaymentProcessed(userID int32, amount float64, success bool) {
	if success {
		s.notificationHandler.NotifyUser(
			userID,
			fmt.Sprintf("Payment of $%.2f processed successfully", amount),
			"success",
			true,
		)
	} else {
		s.notificationHandler.NotifyUserWithData(
			userID,
			fmt.Sprintf("Payment of $%.2f failed. Please update your payment method.", amount),
			"error",
			true, // persistent error
			map[string]interface{}{
				"amount":          amount,
				"action_required": true,
				"retry_url":       "/payment/retry",
			},
		)
	}
}

// Example: System Events
func (s *Server) OnSystemEvent(eventType string, message string) {
	var notificationType string
	var persistent bool

	switch eventType {
	case "maintenance_start":
		notificationType = "warning"
		persistent = true
	case "maintenance_end":
		notificationType = "success"
		persistent = false
	case "critical_error":
		notificationType = "error"
		persistent = true
	default:
		notificationType = "info"
		persistent = false
	}

	s.notificationHandler.NotifyAll(message, notificationType, persistent)
}

// Example: Custom Event Handling
func (s *Server) SetupCustomEventHandlers() {
	// Handle custom events from frontend
	s.socketHandler.OnEvent("user_typing", func(client *handlers.SocketClient, data interface{}) {
		// Broadcast typing indicator to other users
		if client.UserID != nil {
			s.socketHandler.EmitToAll("user_typing", map[string]interface{}{
				"userId": *client.UserID,
				"typing": true,
			})
		}
	})

	s.socketHandler.OnEvent("user_stopped_typing", func(client *handlers.SocketClient, data interface{}) {
		if client.UserID != nil {
			s.socketHandler.EmitToAll("user_typing", map[string]interface{}{
				"userId": *client.UserID,
				"typing": false,
			})
		}
	})

	// Handle chat messages
	s.socketHandler.OnEvent("chat_message", func(client *handlers.SocketClient, data interface{}) {
		// Process and broadcast chat message
		// ... your chat logic here ...
		s.socketHandler.EmitToAll("new_chat_message", data)
	})
}
