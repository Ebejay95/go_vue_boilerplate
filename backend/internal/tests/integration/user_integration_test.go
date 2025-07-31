package integration

import (
	"context"
	"os"
	"testing"

	"backend-grpc-server/internal/server"
	pb "backend-grpc-server/pb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestUserService_Integration(t *testing.T) {
	// Disable auto-migration for tests
	os.Setenv("AUTO_MIGRATE", "false")

	// Setup test server
	srv := server.NewServer()

	// Start server in goroutine
	go func() {
		srv.ServeGRPC("50051")
	}()

	// Connect to server
	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	defer conn.Close()

	client := pb.NewUserServiceClient(conn)

	// Test create user
	createReq := &pb.CreateUserRequest{
		Name:  "Integration Test User",
		Email: "integration@example.com",
		Age:   28,
		Role:  "user",
	}

	createResp, err := client.CreateUser(context.Background(), createReq)
	require.NoError(t, err)
	assert.NotNil(t, createResp.User)

	// Test get user
	getReq := &pb.GetUserRequest{
		Id: createResp.User.Id,
	}

	getResp, err := client.GetUser(context.Background(), getReq)
	require.NoError(t, err)
	assert.Equal(t, createResp.User.Id, getResp.User.Id)
}
