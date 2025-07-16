package handlers

import (
	"context"
	"fmt"

	"backend-grpc-server/internal/models"
	"backend-grpc-server/internal/storage"
	pb "backend-grpc-server/pb"
)

// UserHandler handles user-related gRPC requests
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	store storage.UserStore
}

// NewUserHandler creates a new user handler
func NewUserHandler(store storage.UserStore) *UserHandler {
	return &UserHandler{
		store: store,
	}
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, exists := h.store.GetUser(req.Id)
	if !exists {
		return nil, fmt.Errorf("user with ID %d not found", req.Id)
	}

	return &pb.GetUserResponse{
		User: &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
			Role:  user.Role,
		},
	}, nil
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	user := &models.User{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
		Role:  req.Role,
	}

	createdUser := h.store.CreateUser(user)

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:    createdUser.ID,
			Name:  createdUser.Name,
			Email: createdUser.Email,
			Age:   createdUser.Age,
			Role:  createdUser.Role,
		},
	}, nil
}

// ListUsers returns all users
func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	users := h.store.ListUsers()

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Age:   user.Age,
			Role:  user.Role,
		})
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
	}, nil
}