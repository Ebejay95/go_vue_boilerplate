# Development Dockerfile für Vue.js mit gRPC-Web Hot Reload
FROM node:18-alpine

# Install protoc and tools
RUN apk add --no-cache protobuf-dev curl unzip bash

# Install protoc-gen-grpc-web für JavaScript Generation
RUN curl -sSL https://github.com/grpc/grpc-web/releases/download/1.4.2/protoc-gen-grpc-web-1.4.2-linux-x86_64 \
	-o /usr/local/bin/protoc-gen-grpc-web && \
	chmod +x /usr/local/bin/protoc-gen-grpc-web

WORKDIR /app

# Package.json und lock files zuerst kopieren
COPY frontend/package*.json ./

# Install dependencies including grpc-tools for protoc-gen-js
RUN npm install
RUN npm install -g grpc-tools@latest

# Konfigurationsdateien kopieren
COPY frontend/vue.config.js ./
COPY frontend/tailwind.config.js ./

# Proto files kopieren und JavaScript generieren
COPY proto ../proto
RUN mkdir -p src/proto

# STEP 1: Generate JavaScript files from proto
RUN echo "=== Generating JavaScript protobuf files ==="
RUN npx grpc_tools_node_protoc \
	--js_out=import_style=commonjs,binary:src/proto \
	--proto_path=../proto \
	../proto/*.proto

# STEP 2: Generate gRPC-Web files
RUN echo "=== Generating gRPC-Web files ==="
RUN npx grpc_tools_node_protoc \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:src/proto \
	--proto_path=../proto \
	../proto/*.proto

# Verification: Check generated files exist
RUN echo "=== Verifying generated files ==="
RUN ls -la src/proto/
RUN echo "=== Checking user_pb.js ===" && test -f src/proto/user_pb.js && echo "✅ user_pb.js exists" || echo "❌ user_pb.js missing"
RUN echo "=== Checking user_grpc_web_pb.js ===" && test -f src/proto/user_grpc_web_pb.js && echo "✅ user_grpc_web_pb.js exists" || echo "❌ user_grpc_web_pb.js missing"

# Show first few lines of generated files for debugging
RUN echo "=== user_pb.js content preview ==="
RUN head -10 src/proto/user_pb.js 2>/dev/null || echo "Could not read user_pb.js"

RUN echo "=== user_grpc_web_pb.js content preview ==="
RUN head -10 src/proto/user_grpc_web_pb.js 2>/dev/null || echo "Could not read user_grpc_web_pb.js"

# Source files kopieren NACH proto generation
COPY frontend/src ./src
COPY frontend/public ./public

# WICHTIG: Proto files in src/ überschreiben, falls sie durch COPY überschrieben wurden
RUN echo "=== Re-verifying proto files after src copy ==="
RUN ls -la src/proto/
RUN test -f src/proto/user_pb.js && echo "✅ user_pb.js still exists" || (echo "❌ user_pb.js disappeared, regenerating..." && \
	npx grpc_tools_node_protoc --js_out=import_style=commonjs,binary:src/proto --proto_path=../proto ../proto/*.proto)
RUN test -f src/proto/user_grpc_web_pb.js && echo "✅ user_grpc_web_pb.js still exists" || (echo "❌ user_grpc_web_pb.js disappeared, regenerating..." && \
	npx grpc_tools_node_protoc --grpc-web_out=import_style=commonjs,mode=grpcwebtext:src/proto --proto_path=../proto ../proto/*.proto)

# Final verification
RUN echo "=== Final file check ==="
RUN ls -la src/proto/
RUN file src/proto/* || echo "No proto files found"

# KRITISCHE Environment Variablen für File Watching
ENV NODE_ENV=development
ENV CHOKIDAR_USEPOLLING=true
ENV CHOKIDAR_INTERVAL=500
ENV WATCHPACK_POLLING=true
ENV FAST_REFRESH=true

# Vue CLI spezifische Env Vars
ENV VUE_CLI_BABEL_TRANSPILE_MODULES=true
ENV VUE_CLI_BABEL_TARGET_NODE=true

# Ports freigeben
EXPOSE ${FRONTEND_PORT}

# Development command with proto file check
CMD ["sh", "-c", "echo '=== Starting with proto files: ===' && ls -la src/proto/ && npm run serve"]
