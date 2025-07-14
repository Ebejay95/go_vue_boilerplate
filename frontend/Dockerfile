# Multi-stage build für maximale Effizienz
FROM node:18-alpine AS builder

# Install protoc and protoc-gen-grpc-web
RUN apk add --no-cache protobuf-dev curl unzip

# Install protoc-gen-grpc-web für gRPC-Web Generation
RUN curl -sSL https://github.com/grpc/grpc-web/releases/download/1.4.2/protoc-gen-grpc-web-1.4.2-linux-x86_64 \
	-o /usr/local/bin/protoc-gen-grpc-web && \
	chmod +x /usr/local/bin/protoc-gen-grpc-web

WORKDIR /app

# Dependencies first (besseres Caching)
COPY frontend/package*.json ./

# Smart dependency installation: npm ci if lock file exists, npm install otherwise
RUN if [ -f package-lock.json ]; then \
	echo "=== Using npm ci (package-lock.json found) ===" && \
	npm ci --no-audit; \
	else \
	echo "=== Using npm install (no package-lock.json found) ===" && \
	npm install --no-audit; \
	fi

# Proto generation needs grpc-tools, so ensure it's available
RUN npm list grpc-tools >/dev/null 2>&1 || npm install grpc-tools@latest

# Proto generation
COPY proto ../proto
RUN mkdir -p src/proto

# Generate proto files
RUN npx grpc_tools_node_protoc \
	--js_out=import_style=commonjs,binary:src/proto \
	--proto_path=../proto ../proto/*.proto

RUN npx grpc_tools_node_protoc \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:src/proto \
	--proto_path=../proto ../proto/*.proto

# Copy source code
COPY frontend/ ./

# Build Args - diese werden zur Build-Zeit übergeben
ARG VUE_APP_GRPC_WEB_URL
ARG NODE_ENV=production

# Environment Variables für den Build setzen
ENV VUE_APP_GRPC_WEB_URL=${VUE_APP_GRPC_WEB_URL}
ENV NODE_ENV=${NODE_ENV}

# Debug: Build Args anzeigen
RUN echo "Building with VUE_APP_GRPC_WEB_URL=${VUE_APP_GRPC_WEB_URL}"

# Build app with environment variables
RUN npm run build

# Production image
FROM nginx:alpine

# Remove default nginx config
RUN rm /etc/nginx/conf.d/default.conf

# Copy our template from project root to the templates directory
COPY nginx.conf.template /etc/nginx/templates/default.conf.template

# Copy built app
COPY --from=builder /app/dist /usr/share/nginx/html

# Create nginx cache directories
RUN mkdir -p /var/cache/nginx/client_temp \
	/var/cache/nginx/proxy_temp \
	/var/cache/nginx/fastcgi_temp \
	/var/cache/nginx/uwsgi_temp \
	/var/cache/nginx/scgi_temp

# Set default environment variables (can be overridden by docker-compose)
ENV FRONTEND_PORT=3000
ENV BACKEND_HOST=backend-grpc-server
ENV GRPC_WEB_PORT=8081

EXPOSE ${FRONTEND_PORT}

# Let nginx handle the template processing automatically
CMD ["nginx", "-g", "daemon off;"]
