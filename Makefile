.PHONY: build up down logs clean proto

# Build and start all services
up:
	docker-compose up --build

# Start services in background
up-detached:
	docker-compose up --build -d

# Stop all services
down:
	docker-compose down

# Show logs
logs:
	docker-compose logs -f

# Show logs for specific service
logs-server:
	docker-compose logs -f grpc-server

logs-frontend:
	docker-compose logs -f frontend

# Clean up everything
clean:
	docker-compose down -v --rmi all --remove-orphans
	docker system prune -f

# Generate protobuf files locally (optional)
proto:
	mkdir -p server/pb
	protoc --go_out=./server/pb --go-grpc_out=./server/pb --proto_path=./proto proto/*.proto

# Build only specific service
build-server:
	docker-compose build grpc-server

build-frontend:
	docker-compose build frontend

# Restart specific service
restart-server:
	docker-compose restart grpc-server

restart-frontend:
	docker-compose restart frontend
