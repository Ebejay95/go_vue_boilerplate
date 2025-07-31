.PHONY: dev-up dev-down dev-shell prod-up prod-down prod-shell fclean help migrate-up migrate-down migrate-status build-migrate migrate-clean test test-unit test-integration test-coverage test-race test-docker test-ci

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

TEST_TIMEOUT=30s
COVERAGE_OUT=coverage.out

# Test configuration

# Colors
GREEN := \033[0;32m
YELLOW := \033[0;33m
CYAN := \033[0;36m
BOLD := \033[1m
RESET := \033[0m

define print_status
	@echo "$(CYAN)$(BOLD)[INFO]$(RESET) $(1)"
endef

define print_success
	@echo "$(GREEN)$(BOLD)[SUCCESS]$(RESET) $(1)"
endef

# ===========================================
# DEVELOPMENT
# ===========================================

dev: dev-up
dev-up:
	$(call print_status,"Starting development environment...")
	@rm -rf backend/pb/ backend/go.sum
	docker-compose -f docker-compose.dev.yml up --build
	$(call print_success,"Development environment started!")

dev-up-detached:
	$(call print_status,"Starting development environment...")
	@rm -rf backend/pb/ backend/go.sum
	docker-compose -f docker-compose.dev.yml up --build -d
	$(call print_success,"Development environment started!")

dev-down:
	$(call print_status,"Stopping development environment...")
	docker-compose -f docker-compose.dev.yml down
	$(call print_success,"Development environment stopped!")

dev-shell:
	$(call print_status,"Opening development container shell...")
	docker-compose -f docker-compose.dev.yml exec dev-environment bash

dev-logs:
	$(call print_status,"Showing development logs...")
	docker-compose -f docker-compose.dev.yml logs -f

# ===========================================
# PRODUCTION
# ===========================================
all: up
up: prod-up
prod-up:
	$(call print_status,"Starting production environment...")
	@rm -rf backend/pb/ backend/go.sum
	docker-compose up --build
	$(call print_success,"Production environment started!")

prod-down:
	$(call print_status,"Stopping production environment...")
	docker-compose down
	$(call print_success,"Production environment stopped!")

prod-shell:
	$(call print_status,"Opening production container shell...")
	docker-compose exec backend-grpc-server bash

prod-logs:
	$(call print_status,"Showing production logs...")
	docker-compose logs -f

# ===========================================
# TESTING (using existing containers)
# ===========================================
clean-test-db:
	$(call print_status,"Cleaning test database...")
	docker-compose -f docker-compose.dev.yml exec -T postgres-dev psql -U postgres -d grpc_server_db -c "\
		TRUNCATE TABLE notifications, users RESTART IDENTITY CASCADE;"
	$(call print_success,"Test database cleaned!")
wait-for-services:
	$(call print_status,"Waiting for services to be ready...")
	@timeout=120; \
	echo "Checking if containers are running..."; \
	while [ $$timeout -gt 0 ]; do \
		if docker-compose -f docker-compose.dev.yml ps | grep -q "Up"; then \
			echo "Containers are up, checking database..."; \
			break; \
		fi; \
		echo "Waiting for containers to start... ($$timeout seconds left)"; \
		sleep 2; \
		timeout=$$((timeout-2)); \
	done; \
	if [ $$timeout -le 0 ]; then \
		echo "Timeout waiting for containers to start"; \
		exit 1; \
	fi; \
	echo "Waiting for PostgreSQL to be ready..."; \
	timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		if docker-compose -f docker-compose.dev.yml exec -T postgres-dev pg_isready -h localhost -p 5432 -U postgres >/dev/null 2>&1; then \
			echo "✅ PostgreSQL is ready!"; \
			break; \
		fi; \
		echo "Waiting for PostgreSQL... ($$timeout seconds left)"; \
		sleep 2; \
		timeout=$$((timeout-2)); \
	done; \
	if [ $$timeout -le 0 ]; then \
		echo "❌ Timeout waiting for PostgreSQL"; \
		exit 1; \
	fi; \
	echo "Waiting for backend to be ready..."; \
	timeout=90; \
	while [ $$timeout -gt 0 ]; do \
		if curl -s http://localhost:8081 >/dev/null 2>&1; then \
			echo "✅ Backend is ready!"; \
			break; \
		fi; \
		echo "Waiting for backend... ($$timeout seconds left)"; \
		sleep 3; \
		timeout=$$((timeout-3)); \
	done; \
	if [ $$timeout -le 0 ]; then \
		echo "❌ Timeout waiting for backend"; \
		exit 1; \
	fi; \
	echo "Waiting for frontend to be ready..."; \
	timeout=60; \
	while [ $$timeout -gt 0 ]; do \
		if curl -s http://localhost:3000 >/dev/null 2>&1; then \
			echo "✅ Frontend is ready!"; \
			break; \
		fi; \
		echo "Waiting for frontend... ($$timeout seconds left)"; \
		sleep 3; \
		timeout=$$((timeout-3)); \
	done; \
	if [ $$timeout -le 0 ]; then \
		echo "❌ Timeout waiting for frontend"; \
		exit 1; \
	fi
	$(call print_success,"All services are ready!")
	sleep 60

# Main test target - runs all tests
test: test-dev

# Run tests in dev container
test-dev:
	$(call print_status,"Running tests in dev container...")
	make clean-test-db
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/backend && go test -v -timeout $(TEST_TIMEOUT) ./..."
	$(call print_success,"Tests completed!")

# Run tests with coverage in dev container
test-coverage:
	$(call print_status,"Running coverage tests...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app && go test -v -coverprofile=coverage.out -covermode=atomic ./... && go tool cover -html=coverage.out -o coverage.html"
	$(call print_success,"Coverage report generated!")

# Run race condition tests
test-race:
	$(call print_status,"Running race condition tests...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app && go test -race -short ./..."
	$(call print_success,"Race condition tests completed!")

# CI-friendly test command
test-ci:
	$(call print_status,"Starting CI test pipeline...")
	@rm -rf backend/pb/ backend/go.sum
	docker-compose -f docker-compose.dev.yml up --build -d
	make wait-for-services
	make clean-test-db
	make test-dev
	make dev-down

test-unit: test-dev
test-integration: test-dev

# ===========================================
# CLEANUP
# ===========================================

fclean:
	$(call print_status,"Complete cleanup...")
	docker-compose down -v --rmi all --remove-orphans 2>/dev/null || true
	docker-compose -f docker-compose.dev.yml down -v --rmi all --remove-orphans 2>/dev/null || true
	docker system prune -f
	@rm -rf backend/pb/ backend/go.sum $(COVERAGE_OUT) coverage.html
	$(call print_success,"Complete cleanup finished!")

# ===========================================
# DATABASE MIGRATIONS
# ===========================================

# Database Migration Commands
migrate-up:
	@echo "🔄 Running database migrations..."
	@cd backend && go run cmd/migrate/main.go -action=up

migrate-down:
	@echo "⬇️  Rolling back last migration..."
	@cd backend && go run cmd/migrate/main.go -action=down

migrate-status:
	@echo "📊 Checking migration status..."
	@cd backend && go run cmd/migrate/main.go -action=status

# Build migration tool
build-migrate:
	@echo "🔨 Building migration tool..."
	@cd backend && go build -o bin/migrate cmd/migrate/main.go

# Reset database (careful!)
migrate-reset:
	@echo "⚠️  Resetting database (this will drop all data)..."
	@read -p "Are you sure? Type 'yes' to continue: " confirm && [ "$$confirm" = "yes" ]
	@docker-compose exec postgres psql -U postgres -d grpc_server_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@$(MAKE) migrate-up

# Development helpers
dev-setup: migrate-up
	@echo "✅ Development database setup complete"

# Quick fix for current issue
fix-db:
	@echo "🔧 Applying immediate fix for notifications table..."
	@docker-compose exec postgres psql -U postgres -d grpc_server_db -c "\
		ALTER TABLE notifications \
		ADD COLUMN IF NOT EXISTS user_id INTEGER REFERENCES users(id) ON DELETE CASCADE, \
		ADD COLUMN IF NOT EXISTS read BOOLEAN DEFAULT FALSE, \
		ADD COLUMN IF NOT EXISTS persistent BOOLEAN DEFAULT TRUE; \
		UPDATE notifications SET persistent = TRUE WHERE persistent IS NULL; \
		UPDATE notifications SET read = FALSE WHERE read IS NULL; \
		CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id); \
		CREATE INDEX IF NOT EXISTS idx_notifications_read ON notifications(read); \
		CREATE INDEX IF NOT EXISTS idx_notifications_persistent ON notifications(persistent);"
	@echo "✅ Database fix applied"

migrate-clean:
	@echo "🧹 Cleaning up..."
	@cd backend && rm -f bin/migrate

# ===========================================
# HELP
# ===========================================

help:
	@echo "$(BOLD)Available commands:$(RESET)"
	@echo ""
	@echo "$(YELLOW)Development:$(RESET)"
	@echo "  make dev-up    Start development environment"
	@echo "  make dev-down  Stop development environment"
	@echo "  make dev-shell Open development container shell"
	@echo ""
	@echo "$(YELLOW)Production:$(RESET)"
	@echo "  make prod-up   Start production environment"
	@echo "  make prod-down Stop production environment"
	@echo "  make prod-shell Open production container shell"
	@echo ""
	@echo "$(YELLOW)Testing:$(RESET)"
	@echo "  make test              Run all tests in dev container"
	@echo "  make test-coverage     Run coverage tests"
	@echo "  make test-race         Run race condition tests"
	@echo "  make test-ci           CI-friendly test pipeline"
	@echo ""
	@echo "$(YELLOW)Database:$(RESET)"
	@echo "  make migrate-up    Run database migrations"
	@echo "  make migrate-down  Rollback migrations"
	@echo "  make migrate-status Check migration status"
	@echo ""
	@echo "$(YELLOW)Cleanup:$(RESET)"
	@echo "  make fclean    Complete cleanup (containers, images, volumes)"
	@echo ""
	@echo "$(CYAN)Default: make help$(RESET)"

# Default target
.DEFAULT_GOAL := help
