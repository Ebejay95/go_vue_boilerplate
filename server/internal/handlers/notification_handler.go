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
			Id:      notification.ID,
			Message: notification.Message,
			Type:    notification.Type,
		},
	}, nil
}

// CreateNotification creates a new notification
func (h *NotificationHandler) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.CreateNotificationResponse, error) {
	notification := &models.Notification{
		Message: req.Message,
		Type:    req.Type,
	}

	createdNotification := h.store.CreateNotification(notification)

	return &pb.CreateNotificationResponse{
		Notification: &pb.Notification{
			Id:      createdNotification.ID,
			Message: createdNotification.Message,
			Type:    createdNotification.Type,
		},
	}, nil
}

// ListNotifications returns all notifications
func (h *NotificationHandler) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.ListNotificationsResponse, error) {
	notifications := h.store.ListNotifications()

	var pbNotifications []*pb.Notification
	for _, notification := range notifications {
		pbNotifications = append(pbNotifications, &pb.Notification{
			Id:      notification.ID,
			Message: notification.Message,
			Type:    notification.Type,
		})
	}

	return &pb.ListNotificationsResponse{
		Notifications: pbNotifications,
	}, nil
}