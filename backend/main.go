package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"backend-grpc-server/internal/server"
)

func main() {
	// Get ports from environment
	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("BACKEND_PORT")
	}
	if port == "" {
		port = "8080"
	}

	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		webPort = os.Getenv("GRPC_WEB_PORT")
	}
	if webPort == "" {
		webPort = "8081" // default web port
	}

	socketPort := os.Getenv("WS_PORT")
	if socketPort == "" {
		socketPort = "8082" // default socket port
	}

	// Create server
	srv := server.NewServer()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Received shutdown signal")
		srv.Shutdown()
		os.Exit(0)
	}()

	// Start gRPC server in goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s", port)
		if err := srv.ServeGRPC(port); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start WebSocket server in goroutine
	go func() {
		log.Printf("Starting WebSocket server on port %s", socketPort)

		// Create a new HTTP mux for WebSocket routes
		mux := http.NewServeMux()

		// Register WebSocket routes - alle verwenden denselben Socket Handler
		mux.HandleFunc("/notifications", srv.NewSocketHandler())
		mux.HandleFunc("/users", srv.NewSocketHandler())
		//mux.HandleFunc("/chat", srv.NewSocketHandler())

		// Separate Route for Socket-Status
		mux.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			clients := srv.GetSocketHandler().GetConnectedClients()
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(fmt.Sprintf(`{
				"connected_clients": %d,
				"clients": %v,
				"endpoints": {
					"notifications": "/notifications",
					"users": "/users",
					"chat": "/chat"
				}
			}`, len(clients), clients)))
		})

		// Add a catch-all handler for the root path
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{
					"service": "WebSocket Server",
					"status": "running",
					"endpoints": {
						"notifications": "/notifications",
						"users": "/users",
						"status": "/status"
					}
				}`))
				return
			}
			http.NotFound(w, r)
		})

		if err := http.ListenAndServe(":"+socketPort, mux); err != nil {
			log.Fatalf("Failed to serve WebSocket: %v", err)
		}
	}()

	// Start HTTP server for gRPC-Web
	log.Printf("Starting gRPC-Web server on port %s", webPort)
	log.Printf("Health check available at: http://localhost:%s/health", webPort)
	log.Printf("WebSocket endpoints available at:")
	log.Printf("  - Notifications: ws://localhost:%s/notifications", socketPort)
	log.Printf("  - Users: ws://localhost:%s/users", socketPort)
	log.Printf("  - Chat: ws://localhost:%s/chat", socketPort)
	log.Printf("  - Status: http://localhost:%s/status", socketPort)

	if err := http.ListenAndServe(":"+webPort, srv.CreateHTTPHandler()); err != nil {
		log.Fatalf("Failed to serve gRPC-Web: %v", err)
	}
}
