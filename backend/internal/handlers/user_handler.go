package handlers

import (
	"context"
	"fmt"

	"backend-grpc-server/internal/models"
	"backend-grpc-server/internal/storage"
	"backend-grpc-server/internal/validation"
	pb "backend-grpc-server/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// UserHandler handles user-related gRPC requests with Socket broadcasting
type UserHandler struct {
	pb.UnimplementedUserServiceServer
	store         storage.UserStore
	socketHandler *SocketHandler
}

// NewUserHandler creates a new user handler with Socket support
func NewUserHandler(store storage.UserStore, socketHandler *SocketHandler) *UserHandler {
	return &UserHandler{
		store:         store,
		socketHandler: socketHandler,
	}
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user ID must be greater than 0")
	}

	user, exists := h.store.GetUser(req.Id)
	if !exists {
		return nil, status.Errorf(codes.NotFound, "user with ID %d not found", req.Id)
	}

	return &pb.GetUserResponse{
		User: h.convertToProtoUser(user),
	}, nil
}

// CreateUser creates a new user and broadcasts the update via socket
func (h *UserHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	params := &models.CreateUserParams{
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
		Role:  req.Role,
	}

	// Validate input
	if err := validation.ValidateStruct(params); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	user, err := h.store.CreateUser(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	// Send socket event about user creation
	h.socketHandler.EmitToAll("user_created", map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
		"role":  user.Role,
	})

	// Send notification about the new user
	h.socketHandler.EmitToAll("notification", map[string]interface{}{
		"id":         fmt.Sprintf("user_create_%d", user.ID),
		"message":    fmt.Sprintf("New user %s was created", user.Name),
		"type":       "info",
		"persistent": false,
		"createdAt":  user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"data": map[string]interface{}{
			"action":  "user_created",
			"user_id": user.ID,
			"name":    user.Name,
		},
	})

	return &pb.CreateUserResponse{
		User: h.convertToProtoUser(user),
	}, nil
}

// UpdateUser updates an existing user and broadcasts the update via socket
func (h *UserHandler) UpdateUser(ctx context.Context, req *pb.UpdateUserRequest) (*pb.UpdateUserResponse, error) {
	params := &models.UpdateUserParams{
		ID:    req.Id,
		Name:  req.Name,
		Email: req.Email,
		Age:   req.Age,
		Role:  req.Role,
	}

	// Validate input
	if err := validation.ValidateStruct(params); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	user, err := h.store.UpdateUser(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	// Send socket event about user update
	h.socketHandler.EmitToAll("user_updated", map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
		"role":  user.Role,
	})

	// Send notification to specific user and admins
	h.socketHandler.EmitToUser(user.ID, "notification", map[string]interface{}{
		"id":         fmt.Sprintf("user_update_%d", user.ID),
		"message":    "Your profile has been updated",
		"type":       "success",
		"persistent": true,
		"createdAt":  user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"data": map[string]interface{}{
			"action":  "user_updated",
			"user_id": user.ID,
		},
	})

	return &pb.UpdateUserResponse{
		User: h.convertToProtoUser(user),
	}, nil
}

// DeleteUser deletes a user by ID and broadcasts the update via socket
func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
	if req.Id <= 0 {
		return &pb.DeleteUserResponse{
			Success: false,
			Message: "user ID must be greater than 0",
		}, nil
	}

	// Get user info before deletion for notification
	user, exists := h.store.GetUser(req.Id)
	var userName string
	if exists {
		userName = user.Name
	}

	err := h.store.DeleteUser(req.Id)
	if err != nil {
		return &pb.DeleteUserResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Send socket event about user deletion
	h.socketHandler.EmitToAll("user_deleted", map[string]interface{}{
		"id":   req.Id,
		"name": userName,
	})

	// Send notification about the user deletion
	if userName != "" {
		h.socketHandler.EmitToAll("notification", map[string]interface{}{
			"id":         fmt.Sprintf("user_delete_%d", req.Id),
			"message":    fmt.Sprintf("User %s was deleted", userName),
			"type":       "warning",
			"persistent": false,
			"createdAt":  "time.Now().Format(\"2006-01-02T15:04:05Z07:00\")",
			"data": map[string]interface{}{
				"action":  "user_deleted",
				"user_id": req.Id,
				"name":    userName,
			},
		})
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

	// Validate pagination params
	if params.Limit < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "limit cannot be negative")
	}
	if params.Offset < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "offset cannot be negative")
	}
	if params.Limit > 1000 {
		return nil, status.Errorf(codes.InvalidArgument, "limit cannot exceed 1000")
	}

	users, total, err := h.store.ListUsers(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list users: %v", err)
	}

	var pbUsers []*pb.User
	for _, user := range users {
		pbUsers = append(pbUsers, h.convertToProtoUser(user))
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
		Total: total,
	}, nil
}

// Helper method to convert model user to proto user
func (h *UserHandler) convertToProtoUser(user *models.User) *pb.User {
	return &pb.User{
		Id:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Age:       user.Age,
		Role:      user.Role,
		CreatedAt: user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
