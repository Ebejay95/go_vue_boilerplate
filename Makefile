.PHONY: dev-up dev-down dev-shell prod-up prod-down prod-shell fclean help migrate-up migrate-down migrate-status build-migrate migrate-clean test test-unit test-integration test-coverage test-race test-docker test-ci test-frontend test-frontend-unit test-frontend-e2e test-frontend-component test-frontend-visual test-frontend-all

GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

TEST_TIMEOUT=30s
COVERAGE_OUT=coverage.out

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
# BACKEND TESTING
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

# Backend test target
test: test-dev test-coverage test-frontend-all

test-dev:
	$(call print_status,"Running backend tests in dev container...")
	make clean-test-db
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/backend && go test -v -timeout $(TEST_TIMEOUT) ./..."
	$(call print_success,"Backend tests completed!")

test-coverage:
	$(call print_status,"Running backend coverage tests sequentially...")
	make clean-test-db
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/backend && go test -v -timeout=90s -p=1 -coverprofile=coverage.out -covermode=atomic ./internal/handlers ./internal/models ./internal/storage ./internal/validation && echo 'Coverage Summary:' && go tool cover -func=coverage.out"
	$(call print_success,"Sequential coverage tests completed!")

test-ci:
	$(call print_status,"Starting CI test pipeline...")
	@rm -rf backend/pb/ backend/go.sum
	docker-compose -f docker-compose.dev.yml up --build -d
	make wait-for-services
	make clean-test-db
	make test
	make dev-down

test-unit: test-dev
test-integration: test-dev

# ===========================================
# FRONTEND TESTING (Updated for Container)
# ===========================================

# Frontend setup - ensure all deps are installed in container
frontend-setup:
	$(call print_status,"Setting up frontend dependencies in container...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && rm -f package-lock.json && npm install --include=dev"
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npm run proto:generate"
	# Verify vitest installation
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npx vitest --version"
	$(call print_success,"Frontend setup completed!")

# Frontend unit tests - run in container
test-frontend-unit:
	$(call print_status,"Running frontend unit tests in container...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npx vitest run"
	$(call print_success,"Frontend unit tests completed!")

test-frontend-unit-watch:
	$(call print_status,"Running frontend unit tests in watch mode...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npx vitest"

test-frontend-unit-ui:
	$(call print_status,"Running frontend unit tests with UI...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npx vitest --ui"

test-frontend-coverage:
	$(call print_status,"Running frontend coverage tests in container...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npx vitest run --coverage"
	$(call print_success,"Frontend coverage report generated!")

# Frontend component tests - run in container
test-frontend-component:
	$(call print_status,"Running frontend component tests in container...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npx cypress run --component"
	$(call print_success,"Frontend component tests completed!")

test-frontend-component-open:
	$(call print_status,"Opening frontend component tests...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npx cypress open --component"

# Frontend E2E tests (requires running backend)
test-frontend-e2e:
	$(call print_status,"Running frontend E2E tests...")
	@if ! docker-compose -f docker-compose.dev.yml ps | grep -q "Up"; then \
		echo "Backend services not running. Starting them..."; \
		make dev-up-detached; \
		make wait-for-services; \
	fi
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npx cypress run"

test-frontend-e2e-open:
	$(call print_status,"Opening frontend E2E tests...")
	@if ! docker-compose -f docker-compose.dev.yml ps | grep -q "Up"; then \
		echo "Backend services not running. Starting them..."; \
		make dev-up-detached; \
		make wait-for-services; \
	fi
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npx cypress open"

# Linting and formatting - run in container
test-frontend-lint:
	$(call print_status,"Running frontend linting in container...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npm run lint"
	$(call print_success,"Frontend linting completed!")

test-frontend-lint-fix:
	$(call print_status,"Fixing frontend linting issues...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npm run lint:fix"
	$(call print_success,"Frontend linting fixes applied!")

test-frontend-format:
	$(call print_status,"Checking frontend formatting...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npx prettier --check \"src/**/*.{vue,js,ts,json,css,scss}\""
	$(call print_success,"Frontend formatting check completed!")

test-frontend-format-fix:
	$(call print_status,"Fixing frontend formatting...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npm run format"
	$(call print_success,"Frontend formatting fixes applied!")

# Debug helpers
debug-frontend-deps:
	$(call print_status,"Checking frontend dependencies in container...")
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && npm list vitest || echo 'vitest not found'"
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && which npx"
	docker-compose -f docker-compose.dev.yml exec -T dev-environment sh -c "cd /app/frontend && ls -la node_modules/.bin/ | grep vitest"

# Local alternatives (for faster development)
test-frontend-unit-local:
	$(call print_status,"Running frontend unit tests locally...")
	cd frontend && npm run test:unit
	$(call print_success,"Local frontend unit tests completed!")

test-frontend-coverage-local:
	$(call print_status,"Running frontend coverage tests locally...")
	cd frontend && npm run test:coverage
	$(call print_success,"Local frontend coverage completed!")

# Comprehensive frontend test suites
test-frontend-all:
	$(call print_status,"Running all frontend tests...")
	make frontend-setup
	# make test-frontend-lint
	# make test-frontend-format
	make test-frontend-unit
	make test-frontend-coverage
	make test-frontend-component
	make test-frontend-e2e
	$(call print_success,"All frontend tests completed!")

# Alternative: run tests locally for speed
test-frontend-all-local:
	$(call print_status,"Running all frontend tests locally...")
	cd frontend && npm ci
	cd frontend && npm run lint
	cd frontend && npx prettier --check "src/**/*.{vue,js,ts,json,css,scss}"
	cd frontend && npm run test:coverage
	$(call print_success,"All local frontend tests completed!")

# Combined test suites
test-all:
	$(call print_status,"Running all tests (backend + frontend)...")
	make test-ci
	make test-frontend-ci
	$(call print_success,"All tests completed!")

test-quick:
	$(call print_status,"Running quick test suite...")
	make test-frontend-unit
	make test-dev
	$(call print_success,"Quick tests completed!")

# Development helpers
test-watch:
	$(call print_status,"Running tests in watch mode...")
	make test-frontend-unit-watch &
	make test-dev

# ===========================================
# CLEANUP
# ===========================================

fclean:
	$(call print_status,"Complete cleanup...")
	docker-compose down -v --rmi all --remove-orphans 2>/dev/null || true
	docker-compose -f docker-compose.dev.yml down -v --rmi all --remove-orphans 2>/dev/null || true
	docker system prune -f
	@rm -rf backend/pb/ backend/go.sum $(COVERAGE_OUT) coverage.html
	@rm -rf frontend/node_modules frontend/dist frontend/coverage frontend/test-results
	@rm -rf frontend/cypress/videos frontend/cypress/screenshots
	$(call print_success,"Complete cleanup finished!")

clean-frontend-cache:
	$(call print_status,"Cleaning frontend cache...")
	@rm -rf frontend/node_modules/.cache
	@rm -rf frontend/dist
	@rm -rf frontend/coverage
	@rm -rf frontend/test-results
	@rm -rf frontend/cypress/videos
	@rm -rf frontend/cypress/screenshots
	$(call print_success,"Frontend cache cleaned!")

# ===========================================
# DATABASE MIGRATIONS
# ===========================================

migrate-up:
	@echo "🔄 Running database migrations..."
	@cd backend && go run cmd/migrate/main.go -action=up

migrate-down:
	@echo "⬇️  Rolling back last migration..."
	@cd backend && go run cmd/migrate/main.go -action=down

migrate-status:
	@echo "📊 Checking migration status..."
	@cd backend && go run cmd/migrate/main.go -action=status

build-migrate:
	@echo "🔨 Building migration tool..."
	@cd backend && go build -o bin/migrate cmd/migrate/main.go

migrate-reset:
	@echo "⚠️  Resetting database (this will drop all data)..."
	@read -p "Are you sure? Type 'yes' to continue: " confirm && [ "$confirm" = "yes" ]
	@docker-compose exec postgres psql -U postgres -d grpc_server_db -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public;"
	@$(MAKE) migrate-up

dev-setup: migrate-up
	@echo "✅ Development database setup complete"

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
# DEVELOPMENT UTILITIES
# ===========================================

# Proto generation
proto-generate:
	$(call print_status,"Generating Proto files...")
	cd frontend && npm run proto:generate
	$(call print_success,"Proto files generated!")

# Install dependencies
install-deps:
	$(call print_status,"Installing all dependencies...")
	cd frontend && npm ci
	cd backend && go mod download
	$(call print_success,"Dependencies installed!")

# Build everything
build-all:
	$(call print_status,"Building all components...")
	cd frontend && npm run build
	cd backend && go build -o bin/server main.go
	$(call print_success,"Build completed!")

# Development server start (frontend only)
dev-frontend:
	$(call print_status,"Starting frontend development server...")
	cd frontend && npm run serve

# Check status of all services
status:
	$(call print_status,"Checking service status...")
	@echo "Docker Compose Status:"
	@docker-compose -f docker-compose.dev.yml ps || echo "Development environment not running"
	@echo ""
	@echo "Service Health Checks:"
	@curl -s http://localhost:8081 >/dev/null 2>&1 && echo "✅ Backend gRPC-Web: Running" || echo "❌ Backend gRPC-Web: Not responding"
	@curl -s http://localhost:3000 >/dev/null 2>&1 && echo "✅ Frontend: Running" || echo "❌ Frontend: Not responding"
	@curl -s http://localhost:50051 >/dev/null 2>&1 && echo "✅ Backend gRPC: Running" || echo "❌ Backend gRPC: Not responding"
	@docker-compose -f docker-compose.dev.yml exec -T postgres-dev pg_isready >/dev/null 2>&1 && echo "✅ PostgreSQL: Running" || echo "❌ PostgreSQL: Not responding"

# ===========================================
# HELP
# ===========================================

help:
	@echo "$(BOLD)Available commands:$(RESET)"
	@echo ""
	@echo "$(YELLOW)Development:$(RESET)"
	@echo "  make dev-up              Start development environment"
	@echo "  make dev-down            Stop development environment"
	@echo "  make dev-shell           Open development container shell"
	@echo "  make dev-frontend        Start frontend development server only"
	@echo "  make status              Check status of all services"
	@echo ""
	@echo "$(YELLOW)Production:$(RESET)"
	@echo "  make prod-up             Start production environment"
	@echo "  make prod-down           Stop production environment"
	@echo "  make prod-shell          Open production container shell"
	@echo ""
	@echo "$(YELLOW)Backend Testing:$(RESET)"
	@echo "  make test                Run backend tests"
	@echo "  make test-coverage       Run backend coverage tests"
	@echo "  make test-race           Run backend race condition tests"
	@echo "  make test-ci             Backend CI-friendly test pipeline"
	@echo ""
	@echo "$(YELLOW)Frontend Testing:$(RESET)"
	@echo "  make test-frontend-unit      Run frontend unit tests"
	@echo "  make test-frontend-coverage  Run frontend coverage tests"
	@echo "  make test-frontend-component Run frontend component tests"
	@echo "  make test-frontend-e2e       Run frontend E2E tests"
	@echo "  make test-frontend-visual    Run frontend visual regression tests"
	@echo "  make test-frontend-lint      Run frontend linting"
	@echo "  make test-frontend-format    Check frontend formatting"
	@echo "  make test-frontend-all       Run all frontend tests"
	@echo "  make test-frontend-ci        Frontend CI-friendly test pipeline"
	@echo ""
	@echo "$(YELLOW)Combined Testing:$(RESET)"
	@echo "  make test-all            Run all tests (backend + frontend)"
	@echo "  make test-quick          Run quick test suite"
	@echo ""
	@echo "$(YELLOW)Database:$(RESET)"
	@echo "  make migrate-up          Run database migrations"
	@echo "  make migrate-down        Rollback migrations"
	@echo "  make migrate-status      Check migration status"
	@echo "  make fix-db              Apply database fixes"
	@echo ""
	@echo "$(YELLOW)Utilities:$(RESET)"
	@echo "  make install-deps        Install all dependencies"
	@echo "  make proto-generate      Generate Proto files"
	@echo "  make build-all           Build all components"
	@echo "  make frontend-setup      Setup frontend dependencies"
	@echo ""
	@echo "$(YELLOW)Cleanup:$(RESET)"
	@echo "  make fclean              Complete cleanup (containers, images, volumes)"
	@echo "  make clean-frontend-cache Clean frontend cache only"
	@echo ""
	@echo "$(CYAN)Default: make help$(RESET)"

test-frontend-e2e-open:
	$(call print_status,"Opening frontend E2E tests...")
	@if ! docker-compose -f docker-compose.dev.yml ps | grep -q "Up"; then \
		echo "Backend services not running. Starting them..."; \
		make dev-up-detached; \
		make wait-for-services; \
	fi
	cd frontend && npm run test:e
