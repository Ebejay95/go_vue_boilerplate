#!/bin/bash

echo "🚀 Starting Enhanced Hot-Reload Development Environment..."

# Set default values for missing environment variables
export BACKEND_PORT=${BACKEND_PORT:-50051}
export WS_PORT=${WS_PORT:-8082}
export GRPC_WEB_PORT=${GRPC_WEB_PORT:-8081}
export FRONTEND_PORT=${FRONTEND_PORT:-3000}
export PORT=${BACKEND_PORT}
export WEB_PORT=${GRPC_WEB_PORT}

# Debug: Show environment
echo "📍 Current directory: $(pwd)"
echo "📍 Environment Variables:"
echo "  - WS_PORT: ${WS_PORT}"
echo "  - BACKEND_PORT: ${BACKEND_PORT}"
echo "  - GRPC_WEB_PORT: ${GRPC_WEB_PORT}"
echo "  - FRONTEND_PORT: ${FRONTEND_PORT}"
echo "  - PORT: ${PORT}"
echo "  - WEB_PORT: ${WEB_PORT}"

# CRITICAL: Ensure Go bin directories are in PATH
export PATH="/go/bin:/go/bin/linux_amd64:$PATH"
echo "📍 Updated PATH: $PATH"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
log_info() { echo -e "${CYAN}[INFO]${NC} $1"; }
log_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
log_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
log_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to kill all related processes
kill_all_processes() {
    log_info "🧹 Cleaning up existing processes..."

    # Kill Air processes
    pkill -f "air" 2>/dev/null || true
    pkill -f "./tmp/main" 2>/dev/null || true
    pkill -f "main" 2>/dev/null || true

    # Kill frontend processes
    pkill -f "npm run serve" 2>/dev/null || true
    pkill -f "vue-cli-service" 2>/dev/null || true

    # Kill file watchers
    pkill -f "inotifywait" 2>/dev/null || true

    # Wait for clean shutdown
    sleep 2

    # Force kill if needed
    pkill -9 -f "air" 2>/dev/null || true
    pkill -9 -f "./tmp/main" 2>/dev/null || true
    pkill -9 -f "vue-cli-service" 2>/dev/null || true

    log_success "Processes cleaned up"
}

# Check required tools
log_info "Checking required tools..."

# Check if air exists in multiple locations
AIR_COMMAND=""
if command_exists air; then
    AIR_COMMAND="air"
    log_success "air found in PATH at $(which air)"
elif [ -f "/go/bin/air" ]; then
    AIR_COMMAND="/go/bin/air"
    export PATH="/go/bin:$PATH"
    log_success "air found at /go/bin/air and added to PATH"
elif [ -f "/go/bin/linux_amd64/air" ]; then
    AIR_COMMAND="/go/bin/linux_amd64/air"
    export PATH="/go/bin/linux_amd64:$PATH"
    log_success "air found at /go/bin/linux_amd64/air and added to PATH"
else
    log_error "air not found in common locations"
    exit 1
fi

# Test Air command
log_info "Testing Air installation..."
if $AIR_COMMAND -v >/dev/null 2>&1; then
    log_success "Air is working correctly"
else
    log_warning "Air command test failed, but continuing anyway..."
fi

# Check other tools
for tool in protoc node npm go; do
    if command_exists "$tool"; then
        log_success "$tool found"
    else
        log_error "$tool not found!"
        exit 1
    fi
done

# Check Go protoc plugins
log_info "Checking Go protoc plugins..."
if command_exists protoc-gen-go; then
    log_success "protoc-gen-go found at $(which protoc-gen-go)"
else
    log_error "protoc-gen-go not found in PATH"
    exit 1
fi

if command_exists protoc-gen-go-grpc; then
    log_success "protoc-gen-go-grpc found at $(which protoc-gen-go-grpc)"
else
    log_error "protoc-gen-go-grpc not found in PATH"
    exit 1
fi

# Clean up any existing processes
kill_all_processes

# Create necessary directories
mkdir -p backend/pb backend/tmp frontend/src/proto

# Check directory structure
log_info "📂 Checking directory structure..."
echo "  - /app/backend exists: $(test -d /app/backend && echo 'YES' || echo 'NO')"
echo "  - /app/frontend exists: $(test -d /app/frontend && echo 'YES' || echo 'NO')"
echo "  - /app/proto exists: $(test -d /app/proto && echo 'YES' || echo 'NO')"

if [ ! -d "/app/backend" ]; then
    log_error "Backend directory not found at /app/backend"
    log_info "Available directories in /app:"
    ls -la /app/
    exit 1
fi

# PID tracking for cleanup
PIDS=()

# Cleanup function
cleanup() {
    log_warning "Shutting down all processes..."
    kill_all_processes
    log_success "Cleanup completed"
    exit 0
}

trap cleanup SIGTERM SIGINT

# Function to generate proto files
generate_proto() {
    log_info "🔨 Generating proto files..."

    # Generate Go files
    if protoc --go_out=backend/pb --go_opt=paths=source_relative \
             --go-grpc_out=backend/pb --go-grpc_opt=paths=source_relative \
             --proto_path=proto proto/*.proto; then
        log_success "Go proto files generated"
    else
        log_error "Failed to generate Go proto files"
        return 1
    fi

    # Generate JavaScript files
    cd frontend
    if npx grpc_tools_node_protoc \
        --js_out=import_style=commonjs,binary:src/proto \
        --proto_path=../proto ../proto/*.proto; then
        log_success "JavaScript proto files generated"
    else
        log_error "Failed to generate JavaScript proto files"
        cd ..
        return 1
    fi

    # Generate gRPC-Web files
    if npx grpc_tools_node_protoc \
        --grpc-web_out=import_style=commonjs,mode=grpcwebtext:src/proto \
        --proto_path=../proto ../proto/*.proto; then
        log_success "gRPC-Web proto files generated"
    else
        log_error "Failed to generate gRPC-Web proto files"
        cd ..
        return 1
    fi

    cd ..
    log_success "All proto files generated successfully"
}

# Start a simple Go file watcher that restarts Air completely
start_simple_go_watcher() {
    log_info "🐹 Starting simple Go file watcher..."

    if ! command -v inotifywait >/dev/null 2>&1; then
        log_error "inotifywait not found - Go file watcher disabled"
        return 1
    fi

    if [ ! -d "backend/" ]; then
        log_error "backend/ directory not found - Go file watcher disabled"
        return 1
    fi

    (
        while inotifywait -e modify,create,delete,move -r backend/ --exclude '(tmp|pb|node_modules|dist|\.log)' 2>/dev/null; do
            log_warning "🐹 Go file changed - Restarting backend..."

            # Kill backend process
            if [ ! -z "${BACKEND_PID}" ] && kill -0 $BACKEND_PID 2>/dev/null; then
                kill -TERM $BACKEND_PID 2>/dev/null || true
                sleep 1
                kill -KILL $BACKEND_PID 2>/dev/null || true
            fi

            # Kill any Air processes
            pkill -f "air" 2>/dev/null || true
            pkill -f "./tmp/main" 2>/dev/null || true

            # Wait and restart
            sleep 2
            cd backend
            $AIR_COMMAND -c .air.toml > ../backend.log 2>&1 &
            BACKEND_PID=$!
            cd ..

            log_success "🔄 Backend restarted (PID: $BACKEND_PID)"
        done
    ) &

    GO_WATCHER_PID=$!
    PIDS+=($GO_WATCHER_PID)
    log_success "Go file watcher started (PID: $GO_WATCHER_PID)"
}

# Initial proto generation
log_info "🔧 Initial setup..."
if ! generate_proto; then
    log_error "Initial proto generation failed"
    exit 1
fi

# Start simple Go watcher
start_simple_go_watcher

# Prepare Go environment
log_info "🔧 Preparing Go environment..."
cd backend
if [ -f "go.mod" ]; then
    go mod tidy
    log_success "Go modules prepared"
else
    log_error "go.mod not found in backend directory"
    ls -la
    exit 1
fi

# Ensure tmp directory exists and is clean
mkdir -p tmp
rm -f tmp/main

cd ..

# Prepare Node.js environment
log_info "🔧 Preparing Node.js environment..."
cd frontend
if [ ! -d "node_modules" ]; then
    log_info "Installing npm dependencies..."
    npm install
fi

# Ensure grpc-tools is available
if ! npm list grpc-tools >/dev/null 2>&1; then
    log_info "Installing grpc-tools..."
    npm install grpc-tools@latest
fi
cd ..

# Create Air config that points to the correct directory
log_info "🔧 Creating Air configuration..."
cat > backend/.air.toml << 'EOF'
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./"
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "node_modules", "dist", ".git", "pb"]
  exclude_file = []
  exclude_regex = ["_test.go", ".*\\.log$", ".*\\.tmp$", ".*\\.pb\\.go$"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = ["cmd", "internal", "handler", "service", "pkg"]
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "3s"
  log = "build-errors.log"
  poll = true
  poll_interval = 1000
  pre_cmd = ["echo '🔨 Building Go application...'"]
  post_cmd = ["echo 'Go server started successfully!'"]
  rerun = false
  rerun_delay = 500
  send_interrupt = false
  stop_on_error = false

[color]
  app = "blue"
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = true

[misc]
  clean_on_exit = true

[screen]
  clear_on_rebuild = true
  keep_scroll = false
EOF

# Start Go backend with Air
log_info "🔥 Starting Go backend with Air hot-reload..."
cd backend

log_info "Using Air command: $AIR_COMMAND"
log_info "Environment check - PORT: $PORT, WEB_PORT: $WEB_PORT"
log_info "Working directory: $(pwd)"
log_info "Go files in directory:"
find . -name "*.go" -type f | head -5

# Start Air in background
$AIR_COMMAND -c .air.toml > ../backend.log 2>&1 &
BACKEND_PID=$!
PIDS+=($BACKEND_PID)
cd ..
log_success "Backend started with Air (PID: $BACKEND_PID)"

# Wait for backend to start with enhanced checking
log_info "⏳ Waiting for backend to start..."
for i in {1..60}; do
    # Check if process is still running
    if ! kill -0 $BACKEND_PID 2>/dev/null; then
        log_error "❌ Backend process died during startup"
        log_info "Backend log:"
        tail -50 backend.log 2>/dev/null || echo "No backend log available"
        exit 1
    fi

    # Check if services are responding
    if curl -s http://localhost:${GRPC_WEB_PORT} >/dev/null 2>&1; then
        log_success "Backend is responding on port ${GRPC_WEB_PORT}"
        break
    fi

    if [ $i -eq 60 ]; then
        log_error "❌ Backend failed to start within 60 seconds"
        log_info "Backend log:"
        tail -50 backend.log 2>/dev/null || echo "No backend log available"
        log_info "Process status:"
        ps aux | grep -E "(air|main)" | grep -v grep || echo "No relevant processes found"
        exit 1
    fi

    if [ $((i % 10)) -eq 0 ]; then
        echo "Waiting... ($i/60)"
    else
        echo -n "."
    fi
    sleep 1
done

# Start Vue.js frontend with hot-reload
log_info "🎨 Starting Vue.js frontend with hot-reload..."
cd frontend
export CHOKIDAR_USEPOLLING=true
export WATCHPACK_POLLING=true
npm run serve -- --host 0.0.0.0 --port ${FRONTEND_PORT} > ../frontend.log 2>&1 &
FRONTEND_PID=$!
PIDS+=($FRONTEND_PID)
cd ..
log_success "Frontend started with hot-reload (PID: $FRONTEND_PID)"

# Wait for frontend to start
log_info "⏳ Waiting for frontend to start..."
for i in {1..60}; do
    if curl -s http://localhost:${FRONTEND_PORT} >/dev/null 2>&1; then
        log_success "Frontend is responding on port ${FRONTEND_PORT}"
        break
    fi
    if [ $i -eq 60 ]; then
        log_error "❌ Frontend failed to start within 60 seconds"
        log_info "Frontend log:"
        tail -20 frontend.log 2>/dev/null || echo "No frontend log available"
        exit 1
    fi
    if [ $((i % 10)) -eq 0 ]; then
        echo "Waiting... ($i/60)"
    else
        echo -n "."
    fi
    sleep 1
done

# Display service information
echo ""
log_success "🎉 Enhanced Hot-Reload Development Environment is ready!"
echo ""
echo -e "${PURPLE}📋 Services:${NC}"
echo -e "   🔗 Backend (gRPC): http://localhost:${BACKEND_PORT}"
echo -e "   🔗 Sockets (ws): http://localhost:${WS_PORT}"
echo -e "   🌐 gRPC-Web: http://localhost:${GRPC_WEB_PORT}"
echo -e "   🎨 Frontend: http://localhost:${FRONTEND_PORT}"
echo ""
echo -e "${PURPLE}🔍 Hot-Reload Features:${NC}"
echo -e "   📝 .proto files → Auto-regeneration + backend restart"
echo -e "   🐹 .go files → Air hot-reload with automatic restart"
echo -e "   🖼️ .vue/.js files → Vue dev server hot-reload"
echo -e "   🎨 .css files → Vue dev server hot-reload"
echo ""
echo -e "${PURPLE}📋 Logs:${NC}"
echo -e "   Backend: tail -f backend.log"
echo -e "   Frontend: tail -f frontend.log"
echo ""
echo -e "${YELLOW}💡 Tips:${NC}"
echo -e "   • Edit any .go file to see Air hot-reload"
echo -e "   • Edit any .vue/.js/.css file to see instant frontend updates"
echo -e "   • Press Ctrl+C to stop all services"
echo ""

# Simple health monitoring loop
log_info "👀 Starting health monitoring..."

while true; do
    sleep 30

    # Check if main processes are still running
    if ! kill -0 $BACKEND_PID 2>/dev/null; then
        log_error "❌ Backend process died!"
        exit 1
    fi

    if ! kill -0 $FRONTEND_PID 2>/dev/null; then
        log_error "❌ Frontend process died!"
        exit 1
    fi

    # Check if services are responding
    if ! curl -s http://localhost:${GRPC_WEB_PORT} >/dev/null 2>&1; then
        log_warning "⚠️ Backend not responding"
    fi

    if ! curl -s http://localhost:${FRONTEND_PORT} >/dev/null 2>&1; then
        log_warning "⚠️ Frontend not responding"
    fi

    log_info "💚 Health check passed - All services running"
done
