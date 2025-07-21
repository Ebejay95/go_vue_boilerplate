// internal/handlers/user_handler.go
package handlers

import (
	"context"
	"fmt"

	"backend-grpc-server/internal/models"
	"backend-grpc-server/internal/storage"
	pb "backend-grpc-server/pb"
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

// CreateUser creates a new user and broadcasts the update via socket
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
		"message":    fmt.Sprintf("Neuer Benutzer %s wurde erstellt", user.Name),
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

// UpdateUser updates an existing user and broadcasts the update via socket
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

	// Send socket event about user update
	h.socketHandler.EmitToAll("user_updated", map[string]interface{}{
		"id":    user.ID,
		"name":  user.Name,
		"email": user.Email,
		"age":   user.Age,
		"role":  user.Role,
	})

	// Send notification about the user update (to specific user and admins)
	h.socketHandler.EmitToUser(user.ID, "notification", map[string]interface{}{
		"id":         fmt.Sprintf("user_update_%d", user.ID),
		"message":    fmt.Sprintf("Ihr Profil wurde aktualisiert"),
		"type":       "success",
		"persistent": true,
		"createdAt":  user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"data": map[string]interface{}{
			"action":  "user_updated",
			"user_id": user.ID,
		},
	})

	// Also notify admins about the update
	h.socketHandler.EmitToGroup("admins", "notification", map[string]interface{}{
		"id":         fmt.Sprintf("admin_user_update_%d", user.ID),
		"message":    fmt.Sprintf("Benutzer %s wurde aktualisiert", user.Name),
		"type":       "info",
		"persistent": false,
		"createdAt":  user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		"data": map[string]interface{}{
			"action":  "user_updated",
			"user_id": user.ID,
			"name":    user.Name,
		},
	})

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

// DeleteUser deletes a user by ID and broadcasts the update via socket
func (h *UserHandler) DeleteUser(ctx context.Context, req *pb.DeleteUserRequest) (*pb.DeleteUserResponse, error) {
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
			"message":    fmt.Sprintf("Benutzer %s wurde gelöscht", userName),
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

// Additional helper methods for the socket system

// BulkDeleteUsers deletes multiple users and sends appropriate notifications
func (h *UserHandler) BulkDeleteUsers(ctx context.Context, userIDs []int32) error {
	var deletedUsers []string

	for _, userID := range userIDs {
		user, exists := h.store.GetUser(userID)
		if !exists {
			continue
		}

		err := h.store.DeleteUser(userID)
		if err != nil {
			// Send error notification for this specific user
			h.socketHandler.EmitToAll("notification", map[string]interface{}{
				"id":         fmt.Sprintf("bulk_delete_error_%d", userID),
				"message":    fmt.Sprintf("Fehler beim Löschen von Benutzer %s", user.Name),
				"type":       "error",
				"persistent": true,
				"createdAt":  "time.Now().Format(\"2006-01-02T15:04:05Z07:00\")",
			})
			continue
		}

		deletedUsers = append(deletedUsers, user.Name)

		// Send individual delete event
		h.socketHandler.EmitToAll("user_deleted", map[string]interface{}{
			"id":   userID,
			"name": user.Name,
		})
	}

	// Send summary notification
	if len(deletedUsers) > 0 {
		h.socketHandler.EmitToAll("notification", map[string]interface{}{
			"id":         fmt.Sprintf("bulk_delete_success_%d", len(deletedUsers)),
			"message":    fmt.Sprintf("%d Benutzer erfolgreich gelöscht", len(deletedUsers)),
			"type":       "success",
			"persistent": false,
			"createdAt":  "time.Now().Format(\"2006-01-02T15:04:05Z07:00\")",
			"data": map[string]interface{}{
				"action":        "bulk_delete",
				"deleted_count": len(deletedUsers),
				"deleted_users": deletedUsers,
			},
		})
	}

	return nil
}

// NotifyUserWelcome sends a welcome notification to a newly created user
func (h *UserHandler) NotifyUserWelcome(userID int32, userName string) {
	h.socketHandler.EmitToUser(userID, "notification", map[string]interface{}{
		"id":         fmt.Sprintf("welcome_%d", userID),
		"message":    fmt.Sprintf("Willkommen %s! Ihr Account wurde erfolgreich erstellt.", userName),
		"type":       "success",
		"persistent": true,
		"createdAt":  "time.Now().Format(\"2006-01-02T15:04:05Z07:00\")",
		"data": map[string]interface{}{
			"action":      "welcome",
			"user_id":     userID,
			"first_login": true,
		},
	})
}

// NotifySystemMaintenance sends maintenance notifications to all users
func (h *UserHandler) NotifySystemMaintenance(message string, maintenanceTime string) {
	h.socketHandler.EmitToAll("notification", map[string]interface{}{
		"id":         fmt.Sprintf("maintenance_%s", maintenanceTime),
		"message":    message,
		"type":       "warning",
		"persistent": true,
		"createdAt":  "time.Now().Format(\"2006-01-02T15:04:05Z07:00\")",
		"data": map[string]interface{}{
			"action":           "maintenance_notice",
			"maintenance_time": maintenanceTime,
			"priority":         "high",
		},
	})
}

// NotifyAdmins sends notifications specifically to admin users
func (h *UserHandler) NotifyAdmins(message string, notificationType string, persistent bool, data map[string]interface{}) {
	h.socketHandler.EmitToGroup("admins", "notification", map[string]interface{}{
		"id":         fmt.Sprintf("admin_%s", "time.Now().Unix()"),
		"message":    message,
		"type":       notificationType,
		"persistent": persistent,
		"createdAt":  "time.Now().Format(\"2006-01-02T15:04:05Z07:00\")",
		"data":       data,
	})
}

// SendCustomUserEvent sends custom events for specific user actions
func (h *UserHandler) SendCustomUserEvent(eventType string, userID int32, data map[string]interface{}) {
	// Send to all clients
	h.socketHandler.EmitToAll(eventType, map[string]interface{}{
		"user_id":   userID,
		"timestamp": "time.Now().Format(\"2006-01-02T15:04:05Z07:00\")",
		"data":      data,
	})

	// Also send to the specific user
	h.socketHandler.EmitToUser(userID, eventType, map[string]interface{}{
		"user_id":   userID,
		"timestamp": "time.Now().Format(\"2006-01-02T15:04:05Z07:00\")",
		"data":      data,
	})
}
