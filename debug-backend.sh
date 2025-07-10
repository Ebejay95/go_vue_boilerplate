#!/bin/bash

# Debug script für Backend Container
echo "🔍 Backend Development Container Debugging"
echo "=========================================="

# Farben
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if container is running
echo -e "${BLUE}📦 Checking container status...${NC}"
docker-compose -f docker-compose.yml -f docker-compose.dev.yml ps backend-grpc-server

# Show recent logs
echo -e "\n${BLUE}📋 Recent logs (last 50 lines):${NC}"
docker-compose -f docker-compose.yml -f docker-compose.dev.yml logs --tail=50 backend-grpc-server

# Check if CompileDaemon is running
echo -e "\n${BLUE}🔄 Checking if CompileDaemon is running:${NC}"
docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec backend-grpc-server ps aux | grep -E "(CompileDaemon|main)" || echo -e "${RED}❌ CompileDaemon not found${NC}"

# Check Go environment
echo -e "\n${BLUE}🐹 Go environment:${NC}"
docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec backend-grpc-server go version

# Check proto files
echo -e "\n${BLUE}📄 Proto files status:${NC}"
docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec backend-grpc-server ls -la pb/ || echo -e "${RED}❌ No pb/ directory${NC}"

# Check listening ports
echo -e "\n${BLUE}🌐 Checking listening ports:${NC}"
docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec backend-grpc-server netstat -tlnp | grep -E "(50051|8081)" || echo -e "${RED}❌ Ports not listening${NC}"

# Check file watch capability
echo -e "\n${BLUE}📁 Testing file watching:${NC}"
echo "Creating test file..."
docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec backend-grpc-server touch /app/test-file.go
sleep 2
echo "Removing test file..."
docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec backend-grpc-server rm -f /app/test-file.go

# Show environment variables
echo -e "\n${BLUE}🔧 Environment variables:${NC}"
docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec backend-grpc-server env | grep -E "(GO_|PORT|WEB_PORT)" | sort

# Test gRPC health
echo -e "\n${BLUE}🏥 Testing gRPC server health:${NC}"
docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec backend-grpc-server nc -z localhost 50051 && echo -e "${GREEN}✅ gRPC port 50051 responding${NC}" || echo -e "${RED}❌ gRPC port 50051 not responding${NC}"
docker-compose -f docker-compose.yml -f docker-compose.dev.yml exec backend-grpc-server nc -z localhost 8081 && echo -e "${GREEN}✅ gRPC-Web port 8081 responding${NC}" || echo -e "${RED}❌ gRPC-Web port 8081 not responding${NC}"

echo -e "\n${GREEN}🎯 Debugging complete!${NC}"
echo -e "${YELLOW}💡 To see live logs: make dev-logs-backend${NC}"
echo -e "${YELLOW}💡 To restart backend: make dev-restart-backend${NC}"
echo -e "${YELLOW}💡 To rebuild backend: make dev-rebuild-backend${NC}"
