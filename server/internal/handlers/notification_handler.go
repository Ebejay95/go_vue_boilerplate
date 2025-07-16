package handlers

import (
	"context"
	"fmt"

	"backend-grpc-server/internal/models"
	"backend-grpc-server/internal/storage"
	pb "backend-grpc-server/pb"
)

// NotificationHandler handles notification-related gRPC requests
type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	store storage.NotificationStore
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(store storage.NotificationStore) *NotificationHandler {
	return &NotificationHandler{
		store: store,
	}
}

// GetNotification retrieves a notification by ID
func (h *NotificationHandler) GetNotification(ctx context.Context, req *pb.GetNotificationRequest) (*pb.GetNotificationResponse, error) {
	notification, exists := h.store.GetNotification(req.Id)
	if !exists {
		return nil, fmt.Errorf("notification with ID %d not found", req.Id)
	}

	return &pb.GetNotificationResponse{
		Notification: &pb.Notification{
			Id:        notification.ID,
			Message:   notification.Message,
			Type:      notification.Type,
			CreatedAt: notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: notification.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// CreateNotification creates a new notification
func (h *NotificationHandler) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.CreateNotificationResponse, error) {
	params := &models.CreateNotificationParams{
		Message: req.Message,
		Type:    req.Type,
	}

	notification, err := h.store.CreateNotification(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return &pb.CreateNotificationResponse{
		Notification: &pb.Notification{
			Id:        notification.ID,
			Message:   notification.Message,
			Type:      notification.Type,
			CreatedAt: notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: notification.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		},
	}, nil
}

// ListNotifications returns all notifications (ohne Pagination da Proto es nicht unterstützt)
func (h *NotificationHandler) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.ListNotificationsResponse, error) {
	// Da das Proto keine Pagination hat, verwenden wir Standard-Parameter
	params := &models.ListNotificationsParams{
		Limit:  100, // Standard-Limit
		Offset: 0,   // Keine Pagination
	}

	notifications, _, err := h.store.ListNotifications(params)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	var pbNotifications []*pb.Notification
	for _, notification := range notifications {
		pbNotifications = append(pbNotifications, &pb.Notification{
			Id:        notification.ID,
			Message:   notification.Message,
			Type:      notification.Type,
			CreatedAt: notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt: notification.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		})
	}

	return &pb.ListNotificationsResponse{
		Notifications: pbNotifications,
	}, nil
}

// DeleteNotification deletes a notification by ID
func (h *NotificationHandler) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.DeleteNotificationResponse, error) {
	err := h.store.DeleteNotification(req.Id)
	if err != nil {
		return &pb.DeleteNotificationResponse{
			Success: false,
		}, nil
	}

	return &pb.DeleteNotificationResponse{
		Success: true,
	}, nil
}
