package main

import (
	"log"
	"myapp/internal/database"
	"myapp/internal/handlers"
	pb "myapp/proto/user"
	"net"
	"net/http"
	"os"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := database.Connect()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Register services
	userServer := handlers.NewUserGRPCServer(db)
	pb.RegisterUserServiceServer(grpcServer, userServer)

	// Enable reflection for grpcurl testing
	reflection.Register(grpcServer)

	// Create gRPC-Web wrapper
	wrappedGrpc := grpcweb.WrapServer(grpcServer,
		grpcweb.WithCorsForRegisteredEndpointsOnly(false),
		grpcweb.WithOriginFunc(func(origin string) bool {
			// Allow all origins for development
			return true
		}),
	)

	// HTTP handler that serves both gRPC-Web and health checks
	handler := func(resp http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/health" {
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte(`{"status": "ok"}`))
			return
		}

		// Handle gRPC-Web requests
		if wrappedGrpc.IsGrpcWebRequest(req) {
			wrappedGrpc.ServeHTTP(resp, req)
			return
		}

		// Handle preflight requests
		if req.Method == "OPTIONS" {
			resp.Header().Set("Access-Control-Allow-Origin", "*")
			resp.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			resp.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, X-User-Agent, X-Grpc-Web")
			resp.WriteHeader(http.StatusOK)
			return
		}

		// Default response for other requests
		http.NotFound(resp, req)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Start HTTP server for gRPC-Web
	log.Printf("Starting gRPC-Web server on port %s", port)
	go func() {
		if err := http.ListenAndServe(":"+port, http.HandlerFunc(handler)); err != nil {
			log.Fatalf("Failed to start HTTP server: %v", err)
		}
	}()

	// Start native gRPC server on a different port
	grpcPort := "9090"
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	log.Printf("Starting native gRPC server on port %s", grpcPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}
