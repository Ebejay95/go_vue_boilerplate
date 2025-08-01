package integration

import (
	"context"
	"os"
	"testing"
	"time"

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

	// Create server
	srv := server.NewServer()

	// Start server in goroutine
	go func() {
		if err := srv.ServeGRPC("50051"); err != nil {
			t.Logf("gRPC server error: %v", err)
		}
	}()

	// Wait for server to start
	time.Sleep(2 * time.Second)

	// Connect to server with retry logic
	var conn *grpc.ClientConn
	var err error

	for i := 0; i < 5; i++ {
		conn, err = grpc.Dial("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err == nil {
			break
		}
		t.Logf("Attempt %d: Failed to connect, retrying... %v", i+1, err)
		time.Sleep(1 * time.Second)
	}
	require.NoError(t, err, "Failed to connect to gRPC server after retries")
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
	require.NoError(t, err, "Failed to create user")
	require.NotNil(t, createResp.User, "Created user should not be nil")
	require.NotZero(t, createResp.User.Id, "Created user ID should not be zero")

	t.Logf("Created user with ID: %d", createResp.User.Id)

	// Test get user
	getReq := &pb.GetUserRequest{
		Id: createResp.User.Id,
	}

	getResp, err := client.GetUser(context.Background(), getReq)
	require.NoError(t, err, "Failed to get user with ID %d", createResp.User.Id)
	require.NotNil(t, getResp.User, "Retrieved user should not be nil")
	assert.Equal(t, createResp.User.Id, getResp.User.Id)
	assert.Equal(t, createResp.User.Name, getResp.User.Name)
	assert.Equal(t, createResp.User.Email, getResp.User.Email)

	// Cleanup
	srv.Shutdown()
}
