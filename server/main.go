package main

import (
	"log"
	"net/http"
	"os"

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

	// Start gRPC server in goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s", port)
		if err := srv.ServeGRPC(port); err != nil {
			log.Fatalf("Failed to serve gRPC: %v", err)
		}
	}()

	// Start HTTP server for gRPC-Web
	log.Printf("Starting gRPC-Web server on port %s", webPort)
	if err := http.ListenAndServe(":"+webPort, srv.CreateHTTPHandler()); err != nil {
		log.Fatalf("Failed to serve gRPC-Web: %v", err)
	}
}