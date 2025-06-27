FROM golang:1.22-alpine AS builder

RUN apk add --no-cache protobuf-dev git
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
RUN go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0

WORKDIR /app

# Copy backend source
COPY . .

# Initialize module
RUN go mod init myapp

# Create proto directory and file
RUN mkdir -p proto

# Create proto file step by step
RUN echo 'syntax = "proto3";' > proto/user.proto && \
    echo '' >> proto/user.proto && \
    echo 'package user;' >> proto/user.proto && \
    echo '' >> proto/user.proto && \
    echo 'option go_package = "myapp/proto/user";' >> proto/user.proto && \
    echo '' >> proto/user.proto && \
    echo 'import "google/protobuf/timestamp.proto";' >> proto/user.proto && \
    echo '' >> proto/user.proto && \
    echo 'message User {' >> proto/user.proto && \
    echo '  int32 id = 1;' >> proto/user.proto && \
    echo '  string name = 2;' >> proto/user.proto && \
    echo '  string email = 3;' >> proto/user.proto && \
    echo '  google.protobuf.Timestamp created_at = 4;' >> proto/user.proto && \
    echo '  google.protobuf.Timestamp updated_at = 5;' >> proto/user.proto && \
    echo '}' >> proto/user.proto && \
    echo '' >> proto/user.proto && \
    echo 'message GetUsersRequest {}' >> proto/user.proto && \
    echo 'message GetUserRequest { int32 id = 1; }' >> proto/user.proto && \
    echo 'message CreateUserRequest { string name = 1; string email = 2; }' >> proto/user.proto && \
    echo 'message UpdateUserRequest { int32 id = 1; string name = 2; string email = 3; }' >> proto/user.proto && \
    echo 'message DeleteUserRequest { int32 id = 1; }' >> proto/user.proto && \
    echo '' >> proto/user.proto && \
    echo 'message GetUsersResponse { repeated User users = 1; }' >> proto/user.proto && \
    echo 'message GetUserResponse { User user = 1; }' >> proto/user.proto && \
    echo 'message CreateUserResponse { User user = 1; }' >> proto/user.proto && \
    echo 'message UpdateUserResponse { User user = 1; }' >> proto/user.proto && \
    echo 'message DeleteUserResponse { string message = 1; }' >> proto/user.proto && \
    echo '' >> proto/user.proto && \
    echo 'service UserService {' >> proto/user.proto && \
    echo '  rpc GetUsers(GetUsersRequest) returns (GetUsersResponse);' >> proto/user.proto && \
    echo '  rpc GetUser(GetUserRequest) returns (GetUserResponse);' >> proto/user.proto && \
    echo '  rpc CreateUser(CreateUserRequest) returns (CreateUserResponse);' >> proto/user.proto && \
    echo '  rpc UpdateUser(UpdateUserRequest) returns (UpdateUserResponse);' >> proto/user.proto && \
    echo '  rpc DeleteUser(DeleteUserRequest) returns (DeleteUserResponse);' >> proto/user.proto && \
    echo '}' >> proto/user.proto

# Generate proto files
RUN protoc --go_out=. --go_opt=paths=source_relative \
           --go-grpc_out=. --go-grpc_opt=paths=source_relative \
           proto/user.proto && \
    echo "=== Generated proto files ===" && \
    find . -name "*.pb.go"

# Install dependencies and build
RUN go get google.golang.org/grpc@v1.58.3 && \
    go get google.golang.org/protobuf@v1.31.0 && \
    go get github.com/lib/pq@v1.10.9 && \
    go get github.com/improbable-eng/grpc-web/go/grpcweb@v0.15.0 && \
    go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080 9090
CMD ["./main"]
