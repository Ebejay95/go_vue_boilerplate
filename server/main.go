package main

import (
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
		port = "8080" // default port
	}

	webPort := os.Getenv("WEB_PORT")
	if webPort == "" {
		webPort = os.Getenv("GRPC_WEB_PORT")
	}
	if webPort == "" {
		webPort = "8081" // default web port
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

	// Start HTTP server for gRPC-Web
	log.Printf("Starting gRPC-Web server on port %s", webPort)
	log.Printf("Health check available at: http://localhost:%s/health", webPort)

	if err := http.ListenAndServe(":"+webPort, srv.CreateHTTPHandler()); err != nil {
		log.Fatalf("Failed to serve gRPC-Web: %v", err)
	}
}
