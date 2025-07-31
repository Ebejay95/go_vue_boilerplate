package handlers

import (
	"context"
	"testing"

	"backend-grpc-server/internal/storage"
	"backend-grpc-server/internal/testutil"
	pb "backend-grpc-server/pb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserHandler_CreateUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := storage.NewPostgresUserStore(db)
	socketHandler := NewSocketHandler()
	handler := NewUserHandler(store, socketHandler)

	req := &pb.CreateUserRequest{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
		Role:  "user",
	}

	resp, err := handler.CreateUser(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, resp.User)
	assert.NotZero(t, resp.User.Id)
	assert.Equal(t, req.Name, resp.User.Name)
	assert.Equal(t, req.Email, resp.User.Email)
	assert.Equal(t, req.Age, resp.User.Age)
	assert.Equal(t, req.Role, resp.User.Role)
}

func TestUserHandler_GetUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := storage.NewPostgresUserStore(db)
	socketHandler := NewSocketHandler()
	handler := NewUserHandler(store, socketHandler)

	// Create a user first
	createReq := &pb.CreateUserRequest{
		Name:  "Get Test User",
		Email: "gettest@example.com",
		Age:   30,
		Role:  "admin",
	}

	createResp, err := handler.CreateUser(context.Background(), createReq)
	require.NoError(t, err)

	// Get the user
	getReq := &pb.GetUserRequest{
		Id: createResp.User.Id,
	}

	getResp, err := handler.GetUser(context.Background(), getReq)

	require.NoError(t, err)
	assert.Equal(t, createResp.User.Id, getResp.User.Id)
	assert.Equal(t, createResp.User.Name, getResp.User.Name)
	assert.Equal(t, createResp.User.Email, getResp.User.Email)
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := storage.NewPostgresUserStore(db)
	socketHandler := NewSocketHandler()
	handler := NewUserHandler(store, socketHandler)

	req := &pb.GetUserRequest{
		Id: 99999,
	}

	_, err := handler.GetUser(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}
