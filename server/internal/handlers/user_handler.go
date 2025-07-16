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
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Age:       user.Age,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// CreateUser creates a new user
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	params := &models.CreateUserParams{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
		Role:  req.Role,
	}

	user, err := h.store.CreateUser(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &pb.CreateUserResponse{
		User: &pb.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Age:       user.Age,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// UpdateUser updates an existing user
func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	params := &models.UpdateUserParams{
		ID:    req.Id,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
		Role:  req.Role,
	}

	user, err := h.store.UpdateUser(params)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &pb.UpdateUserResponse{
		User: &pb.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Age:       user.Age,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// DeleteUser deletes a user by ID
func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	err := h.store.DeleteUser(req.Id)
	if err != nil {
		return &pb.DeleteUserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.DeleteUserResponse{
		Success: true,
		Message: fmt.Sprintf("User with ID %d successfully deleted", req.Id),
	}, nil
}

// ListUsers returns all users with pagination
func (h *UserHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	params := &models.ListUsersParams{
		Limit:  req.Limit,
		Offset: req.Offset,
	}

	users, total, err := h.store.ListUsers(params)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, &pb.User{
			Id:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Age:       user.Age,
			Role:      user.Role,
			CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
		Total: total,
	}, nil
}
