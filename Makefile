#(ls -Rla && echo "=== FILE CONTENTS ===" && find . -type f -not -path "./.git/*" -exec sh -c 'echo "=== {} ==="; cat "{}"' \;) > ./../out

.PHONY: build up down logs clean proto clean-proto clean-all help

RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[0;33m
BLUE := \033[0;34m
MAGENTA := \033[0;35m
CYAN := \033[0;36m
WHITE := \033[0;37m
BOLD := \033[1m
RESET := \033[0m

define print_status
	@echo "$(CYAN)$(BOLD)[INFO]$(RESET) $(1)"
endef

define print_success
	@echo "$(GREEN)$(BOLD)[SUCCESS]$(RESET) $(1)"
endef

define print_error
	@echo "$(RED)$(BOLD)[ERROR]$(RESET) $(1)"
endef

define print_warning
	@echo "$(YELLOW)$(BOLD)[WARNING]$(RESET) $(1)"
endef

all: up

up: clean-proto
	$(call print_status,"Building and starting all services...")
	$(call print_warning,"Proto files will be generated by Docker during build...")
	docker-compose up --build
	$(call print_success,"Services started successfully!")

# Start all services in background (detached) mode
up-detached: clean-proto
	$(call print_status,"Building and starting services in detached mode...")
	$(call print_warning,"Proto files will be generated by Docker during build...")
	docker-compose up --build -d
	$(call print_success,"Services running in background!")

down:
	$(call print_status,"Stopping all services...")
	docker-compose down
	$(call print_success,"Services gestoppt!")

logs:
	$(call print_status,"Showing logs...")
	docker-compose logs -f

logs-server:
	$(call print_status,"Showing gRPC server logs...")
	docker-compose logs -f backend-grpc-server

logs-frontend:
	$(call print_status,"Showing frontend logs...")
	docker-compose logs -f frontend

proto:
	$(call print_status,"Proto files are generated automatically by Docker during build")
	$(call print_warning,"Run 'make up' or 'make build-server' to generate proto files")

proto-extract:
	$(call print_status,"Extracting proto files from Docker container...")
	@mkdir -p server/pb
	@if docker ps -a --format "table {{.Names}}" | grep -q "backend-grpc-server"; then \
		docker cp $$(docker-compose ps -q backend-grpc-server):/app/pb/. server/pb/ && \
		$(call print_success,"Proto files extracted successfully!"); \
	else \
		$(call print_error,"gRPC server container not found. Run 'make up' first."); \
	fi

clean-proto:
	$(call print_status,"Cleaning generated proto files...")
	@rm -rf server/pb/
	@rm -f server/go.sum
	$(call print_success,"Proto files cleaned!")

clean-go:
	$(call print_status,"Cleaning Go modules...")
	@rm -f server/go.sum
	@if command -v go >/dev/null 2>&1; then \
		cd server && go clean -modcache; \
	fi
	$(call print_success,"Go modules cleaned!")

clean-docker:
	$(call print_status,"Cleaning Docker artifacts...")
	docker-compose down -v --rmi all --remove-orphans 2>/dev/null || true
	docker system prune -f
	$(call print_success,"Docker artifacts cleaned!")

clean-all: clean-proto clean-go clean-docker
	$(call print_success,"Complete cleanup finished!")

clean: clean-proto clean-go
	$(call print_success,"Development cleanup finished!")

build-server: clean-proto
	$(call print_status,"Building gRPC server...")
	docker-compose build backend-grpc-server
	$(call print_success,"gRPC server built!")

build-frontend:
	$(call print_status,"Building frontend...")
	docker-compose build frontend
	$(call print_success,"Frontend built!")

restart-server:
	$(call print_status,"Restarting gRPC server...")
	docker-compose restart backend-grpc-server
	$(call print_success,"gRPC server restarted!")

restart-frontend:
	$(call print_status,"Restarting frontend...")
	docker-compose restart frontend
	$(call print_success,"Frontend restarted!")

dev: down clean up-detached logs
	$(call print_success,"Development environment ready!")

debug:
	$(call print_status,"Starting debug container...")
	@if docker ps -a --format "table {{.Names}}" | grep -q "go-debug-container"; then \
		docker rm -f go-debug-container; \
	fi
	docker build -f debug.Dockerfile -t go-debug .
	docker run -it --rm -v $$(pwd):/app --name go-debug-container go-debug sh

status:
	$(call print_status,"Service status:")
	@docker-compose ps
	@echo ""
	$(call print_status,"Proto files:")
	@ls -la server/pb/ 2>/dev/null || echo "$(YELLOW)No proto files generated$(RESET)"
	@echo ""
	$(call print_status,"Go modules:")
	@ls -la server/go.sum 2>/dev/null && echo "$(GREEN)go.sum exists$(RESET)" || echo "$(YELLOW)go.sum missing$(RESET)"

test:
	$(call print_status,"Testing API endpoints...")
	@echo "$(YELLOW)Testing frontend...$(RESET)"
	@curl -s $FRONTEND_PROTOCOL://$FRONTEND_HOST:$FRONTEND_PORT >/dev/null && echo "$(GREEN)✓ Frontend responding$(RESET)" || echo "$(RED)✗ Frontend not responding$(RESET)"
	@echo "$(YELLOW)Testing API...$(RESET)"
	@curl -s $FRONTEND_PROTOCOL://$FRONTEND_HOST:$FRONTEND_PORT/api/users >/dev/null && echo "$(GREEN)✓ API responding$(RESET)" || echo "$(RED)✗ API not responding$(RESET)"

dev: dev-up

dev-up: clean-proto
	$(call print_status,"Starting development environment with live reload...")
	$(call print_warning,"Proto files will be generated by Docker during build...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
	$(call print_success,"Development environment started with live reload!")

dev-up-detached: clean-proto
	$(call print_status,"Starting development environment in detached mode...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build -d
	$(call print_success,"Development environment running in background!")

dev-down:
	$(call print_status,"Stopping development environment...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down
	$(call print_success,"Development environment stopped!")

dev-logs:
	$(call print_status,"Showing development logs...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f

dev-logs-backend:
	$(call print_status,"Showing backend development logs...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f backend-grpc-server

dev-logs-frontend:
	$(call print_status,"Showing frontend development logs...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs -f frontend

dev-rebuild-backend:
	$(call print_status,"Rebuilding backend for development...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml build backend-grpc-server
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml restart backend-grpc-server
	$(call print_success,"Backend rebuilt and restarted!")

dev-rebuild-frontend:
	$(call print_status,"Rebuilding frontend for development...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml build frontend
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml restart frontend
	$(call print_success,"Frontend rebuilt and restarted!")

dev-shell-backend:
	$(call print_status,"Opening shell in backend container...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec backend-grpc-server sh

dev-shell-frontend:
	$(call print_status,"Opening shell in frontend container...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec frontend sh

# Quick development restart (without rebuild)
dev-restart:
	$(call print_status,"Restarting development containers...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml restart
	$(call print_success,"Development containers restarted!")

# Clean development volumes
dev-clean:
	$(call print_status,"Cleaning development volumes...")
	docker-compose -f docker-compose.yml -f docker-compose.dev.yml down -v
	docker volume rm $$(docker volume ls -q | grep $(shell grep APP_NAME .env | cut -d'=' -f2)) 2>/dev/null || true
	$(call print_success,"Development volumes cleaned!")

# Development status
dev-status:
	$(call print_status,"Development environment status:")
	@docker-compose -f docker-compose.yml -f docker-compose.dev.yml ps
