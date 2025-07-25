# Build stage
FROM golang:1.21-alpine AS builder

# Install protoc, protoc-gen-go und protoc-gen-grpc-web
RUN apk add --no-cache protobuf-dev git curl unzip
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

# Install protoc-gen-grpc-web für Frontend Proto Generation
RUN curl -sSL https://github.com/grpc/grpc-web/releases/download/1.4.2/protoc-gen-grpc-web-1.4.2-linux-x86_64 \
	-o /usr/local/bin/protoc-gen-grpc-web && \
	chmod +x /usr/local/bin/protoc-gen-grpc-web

WORKDIR /app

# Copy server code first (including go.mod)
COPY backend/ ./

# Copy proto files and generate Go code
COPY proto ./proto
RUN mkdir -p pb
RUN protoc --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --proto_path=proto proto/*.proto

# Debug output
RUN echo "=== Generated files ===" && ls -la pb/
RUN echo "=== Package declaration ===" && head -1 pb/user.pb.go

# Initialize Go modules and build
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

# Expose both gRPC and gRPC-Web ports
EXPOSE ${WS_PORT} ${BACKEND_PORT} ${GRPC_WEB_PORT}

CMD ["./main"]
