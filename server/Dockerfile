# Build stage
FROM golang:1.21-alpine AS builder

# Install protoc and protoc-gen-go
RUN apk add --no-cache protobuf-dev git
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

WORKDIR /app

# Copy server code first (including go.mod)
COPY server/ ./

# Copy proto files and generate Go code
COPY proto ./proto
RUN mkdir -p pb
RUN protoc --go_out=pb --go_opt=paths=source_relative --go-grpc_out=pb --go-grpc_opt=paths=source_relative --proto_path=proto proto/*.proto

# Debug output
RUN echo "=== Generated files ===" && ls -la pb/
RUN echo "=== Package declaration ===" && head -1 pb/user.pb.go

# Initialize Go modules and build
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -v -o main ./cmd/server

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/main .

EXPOSE 50051

CMD ["./main"]
