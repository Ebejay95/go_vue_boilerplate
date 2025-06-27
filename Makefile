# Makefile for gRPC Fullstack App

.PHONY: help proto build run clean test docker-up docker-down docker-rebuild

# Default target
help:
	@echo "Available commands:"
	@echo "  proto          - Generate protobuf files for Go and JavaScript"
	@echo "  build          - Build the Go backend"
	@echo "  run            - Run the application locally"
	@echo "  clean          - Clean generated files"
	@echo "  test           - Run tests"
	@echo "  docker-up      - Start all services with Docker Compose"
	@echo "  docker-down    - Stop all services"
	@echo "  docker-rebuild - Rebuild and restart services"
	@echo "  grpc-test      - Test gRPC endpoints with grpcurl"

# Generate protobuf files
proto:
	@echo "Generating protobuf files..."
	@mkdir -p proto/user
	@mkdir -p frontend/src/proto
	@protoc --go_out=. --go_opt=paths=source_relative \
	        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	        proto/user.proto
	@echo "Go protobuf files generated successfully"
	@cd frontend && npm run proto
	@echo "JavaScript protobuf files generated successfully"

# Build the Go backend
build: proto
	@echo "Building backend..."
	@cd backend && go mod tidy
	@cd backend && go build -o bin/server ./cmd/server
	@echo "Backend built successfully"

# Run the application locally (requires PostgreSQL running)
run: build
	@echo "Starting backend server..."
	@cd backend && ./bin/server

# Clean generated files
clean:
	@echo "Cleaning generated files..."
	@rm -rf backend/bin/
	@rm -rf backend/proto/
	@rm -rf frontend/src/proto/
	@rm -rf frontend/dist/
	@rm -rf frontend/node_modules/
	@echo "Clean completed"

# Run tests
test:
	@echo "Running Go tests..."
	@cd backend && go test ./...
	@echo "Running frontend tests..."
	@cd frontend && npm test

# Docker commands
docker-up:
	@echo "Starting services with Docker Compose..."
	@docker-compose up -d
	@echo "Services started. Check status with: docker-compose ps"

docker-down:
	@echo "Stopping services..."
	@docker-compose down
	@echo "Services stopped"

docker-rebuild:
	@echo "Rebuilding and restarting services..."
	@docker-compose down
	@docker-compose build --no-cache
	@docker-compose up -d
	@echo "Services rebuilt and restarted"

# Test gRPC endpoints
grpc-test:
	@echo "Testing gRPC endpoints..."
	@echo "1. List available services:"
	@grpcurl -plaintext localhost:9090 list
	@echo ""
	@echo "2. Describe UserService:"
	@grpcurl -plaintext localhost:9090 describe user.UserService
	@echo ""
	@echo "3. Get all users:"
	@grpcurl -plaintext localhost:9090 user.UserService/GetUsers
	@echo ""
	@echo "4. Create a test user:"
	@grpcurl -plaintext -d '{"name": "gRPC Test User", "email": "grpc@test.com"}' \
	         localhost:9090 user.UserService/CreateUser

# Install dependencies
install-deps:
	@echo "Installing Go dependencies..."
	@cd backend && go mod download
	@echo "Installing frontend dependencies..."
	@cd frontend && npm install
	@echo "Dependencies installed"

# Development setup
dev-setup: install-deps proto
	@echo "Development setup completed"
	@echo "Run 'make docker-up' to start the full stack"
	@echo "Or run 'make run' to start only the backend (requires local PostgreSQL)"

# Check if required tools are installed
check-tools:
	@echo "Checking required tools..."
	@command -v protoc >/dev/null 2>&1 || (echo "protoc is required but not installed. Please install Protocol Buffers compiler." && exit 1)
	@command -v go >/dev/null 2>&1 || (echo "Go is required but not installed." && exit 1)
	@command -v node >/dev/null 2>&1 || (echo "Node.js is required but not installed." && exit 1)
	@command -v docker >/dev/null 2>&1 || (echo "Docker is required but not installed." && exit 1)
	@command -v docker-compose >/dev/null 2>&1 || (echo "Docker Compose is required but not installed." && exit 1)
	@echo "All required tools are installed"