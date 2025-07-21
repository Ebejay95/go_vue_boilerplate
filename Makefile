.PHONY: dev-up dev-down dev-shell prod-up prod-down prod-shell fclean help migrate-up migrate-down migrate-status build-migrate migrate-clean

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
# CLEANUP
# ===========================================

fclean:
	$(call print_status,"Complete cleanup...")
	docker-compose down -v --rmi all --remove-orphans 2>/dev/null || true
	docker-compose -f docker-compose.dev.yml down -v --rmi all --remove-orphans 2>/dev/null || true
	docker system prune -f
	@rm -rf backend/pb/ backend/go.sum
	$(call print_success,"Complete cleanup finished!")

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
	@echo "$(YELLOW)Cleanup:$(RESET)"
	@echo "  make fclean    Complete cleanup (containers, images, volumes)"
	@echo ""
	@echo "$(CYAN)Default: make help$(RESET)"

# Default target
.DEFAULT_GOAL := help


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
