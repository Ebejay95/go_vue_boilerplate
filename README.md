# gRPC + Vue.js + Go Boilerplate

A modern, full-stack boilerplate for building applications with gRPC backend (Go) and Vue.js frontend. Features Docker containerization, hot reload for development, and production-ready builds.

## 🚀 Features

- **Backend**: Go with gRPC server
- **Frontend**: Vue.js 3 with Composition API
- **Protocol Buffers**: Type-safe communication between services
- **Docker**: Multi-stage builds for production and development
- **Hot Reload**: Live reloading for both frontend and backend in development
- **Tailwind CSS**: Utility-first CSS framework
- **Express Proxy**: Frontend proxy server for API routing
- **Makefile**: Comprehensive build and development commands

## 📋 Prerequisites

- Docker and Docker Compose
- Make (optional, but recommended for easier commands)

## 🏗️ Project Structure

```
.
├── frontend/           # Vue.js frontend application
│   ├── src/           # Vue.js source code
│   ├── public/        # Static assets
│   ├── Dockerfile     # Production build
│   └── Dockerfile.dev # Development with hot reload
├── server/            # Go gRPC server
│   ├── cmd/server/    # Main server application
│   ├── Dockerfile     # Production build
│   └── Dockerfile.dev # Development with Air live reload
├── proto/             # Protocol Buffer definitions
├── docker-compose.yml # Production configuration
├── docker-compose.dev.yml # Development configuration
├── Makefile          # Build and development commands
└── .env              # Environment variables
```

## ⚡ Quick Start

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd grpc-vue-go-boilerplate
cp .env-sample .env
```

### 2. Configure Environment

Edit `.env` file with your preferred settings:

### 3. Start Development Environment

```bash
# Using Make (recommended)
make dev

# Or using Docker Compose directly
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up --build
```


## 🛠️ Development Mode

Development mode provides hot reload for both frontend and backend, making development fast and efficient.

### Start Development

```bash
make dev                    # Start all services with hot reload
make dev-up-detached       # Start in background
make dev-logs              # View all logs
make dev-logs-frontend     # View only frontend logs
make dev-logs-backend      # View only backend logs
```

### Development Features

- **Frontend Hot Reload**: Vue.js changes are instantly reflected
- **Backend Live Reload**: Go server restarts automatically on code changes
- **File Watching**: Aggressive polling ensures Docker file changes are detected
- **Source Maps**: Debug-friendly builds with source mapping

### Development Commands

```bash
# Restart services
make dev-restart           # Restart all containers
make dev-rebuild-backend   # Rebuild and restart backend
make dev-rebuild-frontend  # Rebuild and restart frontend

# Debugging
make dev-shell-backend     # Open shell in backend container
make dev-shell-frontend    # Open shell in frontend container

# Cleanup
make dev-clean             # Clean development volumes
make dev-down              # Stop development environment
```

## 🚀 Production Mode

Production mode creates optimized builds with multi-stage Docker images.

### Start Production

```bash
make up                    # Start production build
make up-detached          # Start in background
make logs                 # View logs
```

### Production Features

- **Optimized Builds**: Minified Vue.js bundle
- **Multi-stage Docker**: Smaller production images
- **Express Static Serving**: Frontend served via Express
- **Environment Variables**: Production-specific configuration

### Production Commands

```bash
# Container management
make down                  # Stop all services
make restart-server        # Restart backend only
make restart-frontend      # Restart frontend only

# Building
make build-server          # Build backend image
make build-frontend        # Build frontend image

# Monitoring
make status               # Show service status
make test                 # Test API endpoints
```

## 🔧 Available Make Commands

### Core Commands

```bash
make dev                  # Start development environment
make up                   # Start production environment
make down                 # Stop all services
make logs                 # Show all logs
make clean                # Clean development files
make clean-all            # Complete cleanup (Docker + files)
```

### Development Specific

```bash
make dev-up              # Start dev with hot reload
make dev-logs            # Show dev logs
make dev-restart         # Restart dev containers
make dev-shell-backend   # Backend container shell
make dev-shell-frontend  # Frontend container shell
```

### Build & Maintenance

```bash
make build-server        # Build backend image
make build-frontend      # Build frontend image
make proto              # Generate protocol buffers
make proto-extract      # Extract proto files from container
make status             # Show service status
make test               # Test API endpoints
```

## 📡 API Endpoints

The frontend Express server provides REST API endpoints that proxy to the gRPC backend:

- `GET /api/users` - List all users
- `GET /api/users/:id` - Get user by ID
- `POST /api/users` - Create new user
- `GET /api/health` - Health check

### Example API Usage

```javascript
// Create a new user
const response = await fetch("/api/users", {
  method: "POST",
  headers: { "Content-Type": "application/json" },
  body: JSON.stringify({
    name: "John Doe",
    email: "john@example.com",
    age: 30
  })
});

// List all users
const users = await fetch("/api/users").then(r => r.json());
```

## 🔄 Protocol Buffers

Protocol Buffer definitions are in the `proto/` directory. The build process automatically generates Go code and makes it available to both services.

### Adding New Services

1. Define your service in `proto/your-service.proto`
2. Rebuild containers: `make dev-rebuild-backend`
3. Protocol buffers are automatically generated during build

### Example Proto Definition

```protobuf
syntax = "proto3";
package user;

service UserService {
  rpc GetUser(GetUserRequest) returns (GetUserResponse);
  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);
}

message User {
  int32 id = 1;
  string name = 2;
  string email = 3;
}
```

## 🐳 Docker Configuration

### Development Setup

- Uses `Dockerfile.dev` files for hot reload capability
- Mounts source directories as volumes
- Includes development tools (Air for Go, Vue CLI dev server)
- Aggressive file watching for Docker compatibility

### Production Setup

- Multi-stage builds for optimized images
- Minimal runtime containers (Alpine Linux)
- Built assets served efficiently
- Production environment variables

## 🔧 Customization

### Adding Frontend Dependencies

```bash
# Enter frontend container
make dev-shell-frontend

# Install packages
npm install your-package

# Rebuild to persist changes
make dev-rebuild-frontend
```

### Adding Backend Dependencies

```bash
# Enter backend container
make dev-shell-backend

# Add Go modules
go get your-package

# Rebuild
make dev-rebuild-backend
```

### Environment Variables

Key environment variables in `.env`:

- `APP_NAME`: Project name (affects container names)
- `BACKEND_PORT`: gRPC server port
- `FRONTEND_PORT`: Express server port
- `BACKEND_HOST`: Internal Docker hostname for backend

## 🚨 Troubleshooting

### Hot Reload Not Working

1. Ensure `CHOKIDAR_USEPOLLING=true` in development
2. Check file permissions and Docker volume mounts
3. Restart development environment: `make dev-restart`

### gRPC Connection Issues

1. Verify backend container is running: `docker ps`
2. Check backend logs: `make dev-logs-backend`
3. Ensure correct `GRPC_SERVER_URL` in frontend environment

### Build Failures

1. Clean everything: `make clean-all`
2. Check Docker disk space: `docker system df`
3. Rebuild from scratch: `make dev`

## 📚 Technology Stack

- **Backend**: Go 1.21, gRPC, Protocol Buffers
- **Frontend**: Vue.js 3, Composition API, Axios
- **Styling**: Tailwind CSS
- **Build**: Docker, Docker Compose
- **Development**: Air (Go live reload), Vue CLI Dev Server
- **Proxy**: Express.js server for API routing

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test in both development and production modes
5. Submit a pull request

## 📄 License

This project is open source and available under the [MIT License](LICENSE).

---

**Happy coding! 🎉**

For questions or issues, please check the troubleshooting section or open an issue in the repository.
