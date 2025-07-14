#!/bin/bash

echo "🚀 Starting Enhanced Hot-Reload Development Environment..."

# Debug: Show environment
echo "📍 Current directory: $(pwd)"
echo "📍 Environment Variables:"
echo "  - BACKEND_PORT: ${BACKEND_PORT:-50051}"
echo "  - GRPC_WEB_PORT: ${GRPC_WEB_PORT:-8081}"
echo "  - FRONTEND_PORT: ${FRONTEND_PORT:-3000}"

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
    log_info "Available files in /go/bin:"
    ls -la /go/bin/ 2>/dev/null || echo "Directory not found"
    log_info "Searching for air binary recursively..."
    find /go -name "air" -type f 2>/dev/null | head -10 || echo "No air binary found"

    # Try to find any air binary and use it
    AIR_BINARY=$(find /go -name "air" -type f 2>/dev/null | head -1)
    if [ -n "$AIR_BINARY" ]; then
        AIR_COMMAND="$AIR_BINARY"
        log_warning "Found air at $AIR_BINARY, using it directly"
    else
        log_error "No air binary found anywhere!"
        exit 1
    fi
fi

# Test Air command
log_info "Testing Air installation..."
if $AIR_COMMAND -v >/dev/null 2>&1; then
    log_success "Air is working correctly"
else
    log_warning "Air command test failed, but continuing anyway..."
fi

# Check protoc and Go tools
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
    log_info "Searching for protoc-gen-go..."
    find /go -name "protoc-gen-go" -type f 2>/dev/null | head -5
    exit 1
fi

if command_exists protoc-gen-go-grpc; then
    log_success "protoc-gen-go-grpc found at $(which protoc-gen-go-grpc)"
else
    log_error "protoc-gen-go-grpc not found in PATH"
    log_info "Searching for protoc-gen-go-grpc..."
    find /go -name "protoc-gen-go-grpc" -type f 2>/dev/null | head -5
    exit 1
fi

# Create necessary directories
mkdir -p server/pb server/tmp frontend/src/proto

# PID tracking for cleanup
PIDS=()

# Cleanup function
cleanup() {
    log_warning "Shutting down all processes..."
    for pid in "${PIDS[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            log_info "Killing process $pid"
            kill -TERM "$pid" 2>/dev/null || kill -KILL "$pid" 2>/dev/null
        fi
    done

    # Kill any remaining processes
    pkill -f "air" 2>/dev/null || true
    pkill -f "npm run serve" 2>/dev/null || true
    pkill -f "inotifywait" 2>/dev/null || true

    log_success "Cleanup completed"
    exit 0
}

trap cleanup SIGTERM SIGINT

# Function to generate proto files
generate_proto() {
    log_info "🔨 Generating proto files..."

    # Debug: Show what's available
    log_info "Current PATH: $PATH"
    log_info "protoc-gen-go location: $(which protoc-gen-go 2>/dev/null || echo 'NOT FOUND')"
    log_info "protoc-gen-go-grpc location: $(which protoc-gen-go-grpc 2>/dev/null || echo 'NOT FOUND')"

    # Generate Go files
    if protoc --go_out=server/pb --go_opt=paths=source_relative \
             --go-grpc_out=server/pb --go-grpc_opt=paths=source_relative \
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
    log_success "✅ All proto files generated successfully"
}

# Function to start proto file watcher
start_proto_watcher() {
    log_info "🔍 Starting proto file watcher..."

    (
        while inotifywait -e modify,create,delete,move -r proto/ 2>/dev/null; do
            log_warning "📝 Proto file changes detected, regenerating..."
            if generate_proto; then
                log_success "🔄 Proto files regenerated successfully"
                # Touch a Go file to trigger Air reload
                touch server/cmd/server/main.go
            else
                log_error "❌ Proto regeneration failed"
            fi
        done
    ) &

    PROTO_WATCHER_PID=$!
    PIDS+=($PROTO_WATCHER_PID)
    log_success "Proto watcher started (PID: $PROTO_WATCHER_PID)"
}

# Function to start CSS/Asset watcher for frontend
start_asset_watcher() {
    log_info "🎨 Starting asset file watcher..."

    (
        while inotifywait -e modify,create,delete,move -r frontend/src/assets/ 2>/dev/null; do
            log_warning "🎨 Asset file changes detected"
            # Vue's dev server should handle this automatically, but we can trigger explicit reload
            touch frontend/src/App.vue
        done
    ) &

    ASSET_WATCHER_PID=$!
    PIDS+=($ASSET_WATCHER_PID)
    log_success "Asset watcher started (PID: $ASSET_WATCHER_PID)"
}

# Function to start Vue.js file watcher (enhanced)
start_vue_watcher() {
    log_info "🖼️ Starting Vue.js component watcher..."

    (
        while inotifywait -e modify,create,delete,move -r frontend/src/ --exclude 'node_modules|dist|proto' 2>/dev/null; do
            log_warning "🖼️ Vue.js file changes detected"
            # Vue dev server handles this automatically, but we can add custom logic here
        done
    ) &

    VUE_WATCHER_PID=$!
    PIDS+=($VUE_WATCHER_PID)
    log_success "Vue.js watcher started (PID: $VUE_WATCHER_PID)"
}

# Function to start Go dependencies watcher
start_go_deps_watcher() {
    log_info "📦 Starting Go dependencies watcher..."

    (
        while inotifywait -e modify server/go.mod server/go.sum 2>/dev/null; do
            log_warning "📦 Go dependencies changed, running go mod tidy..."
            cd server
            go mod tidy
            cd ..
            log_success "🔄 Go dependencies updated"
        done
    ) &

    GO_DEPS_WATCHER_PID=$!
    PIDS+=($GO_DEPS_WATCHER_PID)
    log_success "Go dependencies watcher started (PID: $GO_DEPS_WATCHER_PID)"
}

# Function to start npm dependencies watcher
start_npm_deps_watcher() {
    log_info "📦 Starting npm dependencies watcher..."

    (
        while inotifywait -e modify frontend/package.json frontend/package-lock.json 2>/dev/null; do
            log_warning "📦 npm dependencies changed, running npm install..."
            cd frontend
            npm install
            cd ..
            log_success "🔄 npm dependencies updated"
        done
    ) &

    NPM_DEPS_WATCHER_PID=$!
    PIDS+=($NPM_DEPS_WATCHER_PID)
    log_success "npm dependencies watcher started (PID: $NPM_DEPS_WATCHER_PID)"
}

# Initial proto generation
log_info "🔧 Initial setup..."
if ! generate_proto; then
    log_error "Initial proto generation failed"
    exit 1
fi

# Start all watchers
start_proto_watcher
start_asset_watcher
start_vue_watcher
start_go_deps_watcher
start_npm_deps_watcher

# Prepare Go environment
log_info "🔧 Preparing Go environment..."
cd server
if [ -f "go.mod" ]; then
    go mod tidy
    log_success "Go modules prepared"
else
    log_error "go.mod not found"
    exit 1
fi
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

# Start Go backend with Air (hot reload)
log_info "🔥 Starting Go backend with Air hot-reload..."
cd server

# Use the detected Air command
log_info "Using Air command: $AIR_COMMAND"
$AIR_COMMAND -c .air.toml > ../backend.log 2>&1 &
BACKEND_PID=$!
PIDS+=($BACKEND_PID)
cd ..
log_success "Backend started with Air (PID: $BACKEND_PID)"

# Wait for backend to start
log_info "⏳ Waiting for backend to start..."
for i in {1..30}; do
    if curl -s http://localhost:${GRPC_WEB_PORT:-8081} >/dev/null 2>&1; then
        log_success "✅ Backend is responding"
        break
    fi
    if [ $i -eq 30 ]; then
        log_error "❌ Backend failed to start within 30 seconds"
        log_info "Backend log:"
        tail -20 backend.log 2>/dev/null || echo "No backend log available"
        exit 1
    fi
    sleep 1
done

# Start Vue.js frontend with hot-reload
log_info "🎨 Starting Vue.js frontend with hot-reload..."
cd frontend
export CHOKIDAR_USEPOLLING=true
export WATCHPACK_POLLING=true
npm run serve -- --host 0.0.0.0 --port ${FRONTEND_PORT:-3000} > ../frontend.log 2>&1 &
FRONTEND_PID=$!
PIDS+=($FRONTEND_PID)
cd ..
log_success "Frontend started with hot-reload (PID: $FRONTEND_PID)"

# Wait for frontend to start
log_info "⏳ Waiting for frontend to start..."
for i in {1..60}; do
    if curl -s http://localhost:${FRONTEND_PORT:-3000} >/dev/null 2>&1; then
        log_success "✅ Frontend is responding"
        break
    fi
    if [ $i -eq 60 ]; then
        log_error "❌ Frontend failed to start within 60 seconds"
        log_info "Frontend log:"
        tail -20 frontend.log 2>/dev/null || echo "No frontend log available"
        exit 1
    fi
    sleep 1
done

# Display service information
echo ""
log_success "🎉 Enhanced Hot-Reload Development Environment is ready!"
echo ""
echo -e "${PURPLE}📋 Services:${NC}"
echo -e "   🔗 Backend (gRPC): http://localhost:${BACKEND_PORT:-50051}"
echo -e "   🌐 gRPC-Web: http://localhost:${GRPC_WEB_PORT:-8081}"
echo -e "   🎨 Frontend: http://localhost:${FRONTEND_PORT:-3000}"
echo ""
echo -e "${PURPLE}🔍 Hot-Reload Features:${NC}"
echo -e "   📝 .proto files → Auto-regeneration + backend restart"
echo -e "   🐹 .go files → Air hot-reload"
echo -e "   🖼️ .vue/.js files → Vue dev server hot-reload"
echo -e "   🎨 .css files → Vue dev server hot-reload"
echo -e "   📦 package.json → Auto npm install"
echo -e "   📦 go.mod → Auto go mod tidy"
echo ""
echo -e "${PURPLE}📋 Logs:${NC}"
echo -e "   Backend: tail -f backend.log"
echo -e "   Frontend: tail -f frontend.log"
echo ""
echo -e "${YELLOW}💡 Tips:${NC}"
echo -e "   • Edit any .proto file to see automatic regeneration"
echo -e "   • Edit any .go file to see Air hot-reload"
echo -e "   • Edit any .vue/.js/.css file to see instant frontend updates"
echo -e "   • Press Ctrl+C to stop all services"
echo ""

# Health monitoring loop
log_info "👀 Starting health monitoring..."
HEALTH_CHECK_INTERVAL=30
LAST_HEALTH_CHECK=0

while true; do
    current_time=$(date +%s)

    # Health check every 30 seconds
    if [ $((current_time - LAST_HEALTH_CHECK)) -ge $HEALTH_CHECK_INTERVAL ]; then
        # Check if all main processes are still running
        if ! kill -0 $BACKEND_PID 2>/dev/null; then
            log_error "❌ Backend process died!"
            log_info "Backend log:"
            tail -20 backend.log 2>/dev/null || echo "No backend log available"
            exit 1
        fi

        if ! kill -0 $FRONTEND_PID 2>/dev/null; then
            log_error "❌ Frontend process died!"
            log_info "Frontend log:"
            tail -20 frontend.log 2>/dev/null || echo "No frontend log available"
            exit 1
        fi

        # Check if services are responding
        if ! curl -s http://localhost:${GRPC_WEB_PORT:-8081} >/dev/null 2>&1; then
            log_warning "⚠️ Backend not responding"
        fi

        if ! curl -s http://localhost:${FRONTEND_PORT:-3000} >/dev/null 2>&1; then
            log_warning "⚠️ Frontend not responding"
        fi

        LAST_HEALTH_CHECK=$current_time
        log_info "💚 Health check passed - All services running"
    fi

    sleep 5
done

# Create necessary directories
mkdir -p server/pb server/tmp frontend/src/proto

# PID tracking for cleanup
PIDS=()

# Cleanup function
cleanup() {
    log_warning "Shutting down all processes..."
    for pid in "${PIDS[@]}"; do
        if kill -0 "$pid" 2>/dev/null; then
            log_info "Killing process $pid"
            kill -TERM "$pid" 2>/dev/null || kill -KILL "$pid" 2>/dev/null
        fi
    done

    # Kill any remaining processes
    pkill -f "air" 2>/dev/null || true
    pkill -f "npm run serve" 2>/dev/null || true
    pkill -f "inotifywait" 2>/dev/null || true

    log_success "Cleanup completed"
    exit 0
}

trap cleanup SIGTERM SIGINT

# Function to generate proto files
generate_proto() {
    log_info "🔨 Generating proto files..."

    # Generate Go files
    if protoc --go_out=server/pb --go_opt=paths=source_relative \
             --go-grpc_out=server/pb --go-grpc_opt=paths=source_relative \
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
    log_success "✅ All proto files generated successfully"
}

# Function to start proto file watcher
start_proto_watcher() {
    log_info "🔍 Starting proto file watcher..."

    (
        while inotifywait -e modify,create,delete,move -r proto/ 2>/dev/null; do
            log_warning "📝 Proto file changes detected, regenerating..."
            if generate_proto; then
                log_success "🔄 Proto files regenerated successfully"
                # Touch a Go file to trigger Air reload
                touch server/cmd/server/main.go
            else
                log_error "❌ Proto regeneration failed"
            fi
        done
    ) &

    PROTO_WATCHER_PID=$!
    PIDS+=($PROTO_WATCHER_PID)
    log_success "Proto watcher started (PID: $PROTO_WATCHER_PID)"
}

# Function to start CSS/Asset watcher for frontend
start_asset_watcher() {
    log_info "🎨 Starting asset file watcher..."

    (
        while inotifywait -e modify,create,delete,move -r frontend/src/assets/ 2>/dev/null; do
            log_warning "🎨 Asset file changes detected"
            # Vue's dev server should handle this automatically, but we can trigger explicit reload
            touch frontend/src/App.vue
        done
    ) &

    ASSET_WATCHER_PID=$!
    PIDS+=($ASSET_WATCHER_PID)
    log_success "Asset watcher started (PID: $ASSET_WATCHER_PID)"
}

# Function to start Vue.js file watcher (enhanced)
start_vue_watcher() {
    log_info "🖼️ Starting Vue.js component watcher..."

    (
        while inotifywait -e modify,create,delete,move -r frontend/src/ --exclude 'node_modules|dist|proto' 2>/dev/null; do
            log_warning "🖼️ Vue.js file changes detected"
            # Vue dev server handles this automatically, but we can add custom logic here
        done
    ) &

    VUE_WATCHER_PID=$!
    PIDS+=($VUE_WATCHER_PID)
    log_success "Vue.js watcher started (PID: $VUE_WATCHER_PID)"
}

# Function to start Go dependencies watcher
start_go_deps_watcher() {
    log_info "📦 Starting Go dependencies watcher..."

    (
        while inotifywait -e modify server/go.mod server/go.sum 2>/dev/null; do
            log_warning "📦 Go dependencies changed, running go mod tidy..."
            cd server
            go mod tidy
            cd ..
            log_success "🔄 Go dependencies updated"
        done
    ) &

    GO_DEPS_WATCHER_PID=$!
    PIDS+=($GO_DEPS_WATCHER_PID)
    log_success "Go dependencies watcher started (PID: $GO_DEPS_WATCHER_PID)"
}

# Function to start npm dependencies watcher
start_npm_deps_watcher() {
    log_info "📦 Starting npm dependencies watcher..."

    (
        while inotifywait -e modify frontend/package.json frontend/package-lock.json 2>/dev/null; do
            log_warning "📦 npm dependencies changed, running npm install..."
            cd frontend
            npm install
            cd ..
            log_success "🔄 npm dependencies updated"
        done
    ) &

    NPM_DEPS_WATCHER_PID=$!
    PIDS+=($NPM_DEPS_WATCHER_PID)
    log_success "npm dependencies watcher started (PID: $NPM_DEPS_WATCHER_PID)"
}

# Initial proto generation
log_info "🔧 Initial setup..."
if ! generate_proto; then
    log_error "Initial proto generation failed"
    exit 1
fi

# Start all watchers
start_proto_watcher
start_asset_watcher
start_vue_watcher
start_go_deps_watcher
start_npm_deps_watcher

# Prepare Go environment
log_info "🔧 Preparing Go environment..."
cd server
if [ -f "go.mod" ]; then
    go mod tidy
    log_success "Go modules prepared"
else
    log_error "go.mod not found"
    exit 1
fi
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

# Start Go backend with Air (hot reload)
log_info "🔥 Starting Go backend with Air hot-reload..."
cd server

# Use the detected Air command
log_info "Using Air command: $AIR_COMMAND"
$AIR_COMMAND -c .air.toml > ../backend.log 2>&1 &
BACKEND_PID=$!
PIDS+=($BACKEND_PID)
cd ..
log_success "Backend started with Air (PID: $BACKEND_PID)"

# Wait for backend to start
log_info "⏳ Waiting for backend to start..."
for i in {1..30}; do
    if curl -s http://localhost:${GRPC_WEB_PORT:-8081} >/dev/null 2>&1; then
        log_success "✅ Backend is responding"
        break
    fi
    if [ $i -eq 30 ]; then
        log_error "❌ Backend failed to start within 30 seconds"
        log_info "Backend log:"
        tail -20 backend.log 2>/dev/null || echo "No backend log available"
        exit 1
    fi
    sleep 1
done

# Start Vue.js frontend with hot-reload
log_info "🎨 Starting Vue.js frontend with hot-reload..."
cd frontend
export CHOKIDAR_USEPOLLING=true
export WATCHPACK_POLLING=true
npm run serve -- --host 0.0.0.0 --port ${FRONTEND_PORT:-3000} > ../frontend.log 2>&1 &
FRONTEND_PID=$!
PIDS+=($FRONTEND_PID)
cd ..
log_success "Frontend started with hot-reload (PID: $FRONTEND_PID)"

# Wait for frontend to start
log_info "⏳ Waiting for frontend to start..."
for i in {1..60}; do
    if curl -s http://localhost:${FRONTEND_PORT:-3000} >/dev/null 2>&1; then
        log_success "✅ Frontend is responding"
        break
    fi
    if [ $i -eq 60 ]; then
        log_error "❌ Frontend failed to start within 60 seconds"
        log_info "Frontend log:"
        tail -20 frontend.log 2>/dev/null || echo "No frontend log available"
        exit 1
    fi
    sleep 1
done

# Display service information
echo ""
log_success "🎉 Enhanced Hot-Reload Development Environment is ready!"
echo ""
echo -e "${PURPLE}📋 Services:${NC}"
echo -e "   🔗 Backend (gRPC): http://localhost:${BACKEND_PORT:-50051}"
echo -e "   🌐 gRPC-Web: http://localhost:${GRPC_WEB_PORT:-8081}"
echo -e "   🎨 Frontend: http://localhost:${FRONTEND_PORT:-3000}"
echo ""
echo -e "${PURPLE}🔍 Hot-Reload Features:${NC}"
echo -e "   📝 .proto files → Auto-regeneration + backend restart"
echo -e "   🐹 .go files → Air hot-reload"
echo -e "   🖼️ .vue/.js files → Vue dev server hot-reload"
echo -e "   🎨 .css files → Vue dev server hot-reload"
echo -e "   📦 package.json → Auto npm install"
echo -e "   📦 go.mod → Auto go mod tidy"
echo ""
echo -e "${PURPLE}📋 Logs:${NC}"
echo -e "   Backend: tail -f backend.log"
echo -e "   Frontend: tail -f frontend.log"
echo ""
echo -e "${YELLOW}💡 Tips:${NC}"
echo -e "   • Edit any .proto file to see automatic regeneration"
echo -e "   • Edit any .go file to see Air hot-reload"
echo -e "   • Edit any .vue/.js/.css file to see instant frontend updates"
echo -e "   • Press Ctrl+C to stop all services"
echo ""

# Health monitoring loop
log_info "👀 Starting health monitoring..."
HEALTH_CHECK_INTERVAL=30
LAST_HEALTH_CHECK=0

while true; do
    current_time=$(date +%s)

    # Health check every 30 seconds
    if [ $((current_time - LAST_HEALTH_CHECK)) -ge $HEALTH_CHECK_INTERVAL ]; then
        # Check if all main processes are still running
        if ! kill -0 $BACKEND_PID 2>/dev/null; then
            log_error "❌ Backend process died!"
            log_info "Backend log:"
            tail -20 backend.log 2>/dev/null || echo "No backend log available"
            exit 1
        fi

        if ! kill -0 $FRONTEND_PID 2>/dev/null; then
            log_error "❌ Frontend process died!"
            log_info "Frontend log:"
            tail -20 frontend.log 2>/dev/null || echo "No frontend log available"
            exit 1
        fi

        # Check if services are responding
        if ! curl -s http://localhost:${GRPC_WEB_PORT:-8081} >/dev/null 2>&1; then
            log_warning "⚠️ Backend not responding"
        fi

        if ! curl -s http://localhost:${FRONTEND_PORT:-3000} >/dev/null 2>&1; then
            log_warning "⚠️ Frontend not responding"
        fi

        LAST_HEALTH_CHECK=$current_time
        log_info "💚 Health check passed - All services running"
    fi

    sleep 5
done
