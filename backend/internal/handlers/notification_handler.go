package handlers

import (
	"context"
	"fmt"
	"time"

	"backend-grpc-server/internal/models"
	"backend-grpc-server/internal/storage"
	pb "backend-grpc-server/pb"
)

// NotificationHandler with socket integration
type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
	store         storage.NotificationStore
	socketHandler *SocketHandler
}

// NewNotificationHandler creates a new notification handler
func NewNotificationHandler(store storage.NotificationStore, socketHandler *SocketHandler) *NotificationHandler {
	handler := &NotificationHandler{
		store:         store,
		socketHandler: socketHandler,
	}

	// Register socket event handlers for notifications
	handler.setupSocketHandlers()

	return handler
}

// setupSocketHandlers registers event handlers for notification-related socket events
func (h *NotificationHandler) setupSocketHandlers() {
	// Handle notification creation requests from frontend
	h.socketHandler.OnEvent("create_notification", func(client *SocketClient, data interface{}) {
		h.handleCreateNotificationEvent(client, data)
	})

	// Handle mark as read requests
	h.socketHandler.OnEvent("mark_as_read", func(client *SocketClient, data interface{}) {
		h.handleMarkAsReadEvent(client, data)
	})

	// Handle delete notification requests
	h.socketHandler.OnEvent("delete_notification", func(client *SocketClient, data interface{}) {
		h.handleDeleteNotificationEvent(client, data)
	})
}

// Socket Event Handlers

func (h *NotificationHandler) handleCreateNotificationEvent(client *SocketClient, data interface{}) {
	// Parse the notification data
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Invalid notification data format",
		})
		return
	}

	message, _ := dataMap["message"].(string)
	notificationType, _ := dataMap["type"].(string)
	persistent, _ := dataMap["persistent"].(bool)

	if message == "" {
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Message is required",
		})
		return
	}

	// Create the notification
	err := h.SendNotification(message, notificationType, "all", nil, persistent, nil)
	if err != nil {
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Failed to create notification: " + err.Error(),
		})
		return
	}

	// Send success response
	h.socketHandler.EmitToClient(client.ID, "notification_created", map[string]interface{}{
		"message": "Notification created successfully",
	})
}

func (h *NotificationHandler) handleMarkAsReadEvent(client *SocketClient, data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	notificationID, ok := dataMap["id"].(float64)
	if !ok {
		return
	}

	userID := int32(1) // TODO: Get from client.UserID
	if client.UserID != nil {
		userID = *client.UserID
	}

	err := h.store.MarkAsRead(int32(notificationID), userID)
	if err != nil {
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Failed to mark as read",
		})
		return
	}

	// Notify the user that notification was marked as read
	h.socketHandler.EmitToUser(userID, "notification_updated", map[string]interface{}{
		"id":   int32(notificationID),
		"read": true,
	})
}

func (h *NotificationHandler) handleDeleteNotificationEvent(client *SocketClient, data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		return
	}

	notificationID, ok := dataMap["id"].(float64)
	if !ok {
		return
	}

	err := h.store.DeleteNotification(int32(notificationID))
	if err != nil {
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Failed to delete notification",
		})
		return
	}

	// Notify all relevant users about deletion
	h.socketHandler.EmitToAll("notification_deleted", map[string]interface{}{
		"id": int32(notificationID),
	})
}

// Backend Notification Methods (callable from anywhere in your backend)

// SendNotification is the main method for sending notifications from backend
func (h *NotificationHandler) SendNotification(
	message string,
	notificationType string,
	targetType string, // "all" or "user"
	targetID *int32, // user ID if targetType is "user"
	persistent bool, // true = save to DB, false = real-time only
	data map[string]interface{}, // optional extra data
) error {
	// Create notification data
	notificationData := map[string]interface{}{
		"id":         fmt.Sprintf("notif_%d", time.Now().UnixNano()),
		"message":    message,
		"type":       notificationType,
		"persistent": persistent,
		"createdAt":  time.Now().Format(time.RFC3339),
		"data":       data,
	}

	// If persistent, save to database first
	if persistent {
		params := &models.CreateNotificationParams{
			Message:    message,
			Type:       notificationType,
			UserID:     targetID,
			Persistent: true,
		}

		dbNotification, err := h.store.CreateNotification(params)
		if err != nil {
			return fmt.Errorf("failed to save notification to database: %w", err)
		}

		// Update notification data with database ID
		notificationData["id"] = dbNotification.ID
		notificationData["createdAt"] = dbNotification.CreatedAt.Format(time.RFC3339)
	}

	// Send via socket based on target type
	switch targetType {
	case "all":
		h.socketHandler.EmitToAll("notification", notificationData)
	case "user":
		if targetID != nil {
			h.socketHandler.EmitToUser(*targetID, "notification", notificationData)
		}
	default:
		return fmt.Errorf("invalid target type: %s", targetType)
	}

	return nil
}

// Convenience methods for different notification types

func (h *NotificationHandler) NotifyAll(message, notificationType string, persistent bool) error {
	return h.SendNotification(message, notificationType, "all", nil, persistent, nil)
}

func (h *NotificationHandler) NotifyUser(userID int32, message, notificationType string, persistent bool) error {
	return h.SendNotification(message, notificationType, "user", &userID, persistent, nil)
}

func (h *NotificationHandler) NotifyUserWithData(userID int32, message, notificationType string, persistent bool, data map[string]interface{}) error {
	return h.SendNotification(message, notificationType, "user", &userID, persistent, data)
}

// gRPC Methods (existing, but enhanced with socket integration)

func (h *NotificationHandler) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.CreateNotificationResponse, error) {
	// Create via gRPC as before, but also send via socket
	params := &models.CreateNotificationParams{
		Message:    req.Message,
		Type:       req.Type,
		UserID:     convertToInt32Pointer(req.UserId),
		Persistent: req.Persistent,
	}

	notification, err := h.store.CreateNotification(params)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Send via socket
	targetType := "all"
	var targetID *int32
	if req.UserId != 0 {
		targetType = "user"
		targetID = &req.UserId
	}

	h.SendNotification(req.Message, req.Type, targetType, targetID, req.Persistent, nil)

	return &pb.CreateNotificationResponse{
		Notification: h.convertToProtoNotification(notification),
	}, nil
}

// GetNotification retrieves a notification by ID
func (h *NotificationHandler) GetNotification(ctx context.Context, req *pb.GetNotificationRequest) (*pb.GetNotificationResponse, error) {
	notification, exists := h.store.GetNotification(req.Id)
	if !exists {
		return nil, fmt.Errorf("notification with ID %d not found", req.Id)
	}

	return &pb.GetNotificationResponse{
		Notification: h.convertToProtoNotification(notification),
	}, nil
}

// ListNotifications returns notifications with filtering
func (h *NotificationHandler) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.ListNotificationsResponse, error) {
	params := &models.ListNotificationsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
		UserID: convertToInt32Pointer(req.UserId),
		Read:   convertToBoolPointer(req.Read, req.HasReadFilter),
	}

	notifications, total, err := h.store.ListNotifications(params)
	if err != nil {
		return nil, fmt.Errorf("failed to list notifications: %w", err)
	}

	var pbNotifications []*pb.Notification
	for _, notification := range notifications {
		pbNotifications = append(pbNotifications, h.convertToProtoNotification(notification))
	}

	return &pb.ListNotificationsResponse{
		Notifications: pbNotifications,
		Total:         total,
	}, nil
}

// MarkNotificationAsRead marks a notification as read and sends socket update
func (h *NotificationHandler) MarkNotificationAsRead(ctx context.Context, req *pb.MarkNotificationAsReadRequest) (*pb.MarkNotificationAsReadResponse, error) {
	err := h.store.MarkAsRead(req.Id, req.UserId)
	if err != nil {
		return &pb.MarkNotificationAsReadResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Send socket update
	h.socketHandler.EmitToUser(req.UserId, "notification_updated", map[string]interface{}{
		"id":   req.Id,
		"read": true,
	})

	return &pb.MarkNotificationAsReadResponse{
		Success: true,
		Message: "Notification marked as read",
	}, nil
}

// DeleteNotification deletes a notification and sends socket update
func (h *NotificationHandler) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.DeleteNotificationResponse, error) {
	err := h.store.DeleteNotification(req.Id)
	if err != nil {
		return &pb.DeleteNotificationResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Send socket update
	h.socketHandler.EmitToAll("notification_deleted", map[string]interface{}{
		"id": req.Id,
	})

	return &pb.DeleteNotificationResponse{
		Success: true,
		Message: fmt.Sprintf("Notification with ID %d successfully deleted", req.Id),
	}, nil
}

// Helper methods
func (h *NotificationHandler) convertToProtoNotification(notification *models.Notification) *pb.Notification {
	var userID int32
	if notification.UserID != nil {
		userID = *notification.UserID
	}

	return &pb.Notification{
		Id:         notification.ID,
		Message:    notification.Message,
		Type:       notification.Type,
		UserId:     userID,
		Read:       notification.Read,
		Persistent: notification.Persistent,
		CreatedAt:  notification.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  notification.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

func convertToInt32Pointer(val int32) *int32 {
	if val == 0 {
		return nil
	}
	return &val
}

func convertToBoolPointer(val bool, hasFilter bool) *bool {
	if !hasFilter {
		return nil
	}
	return &val
}
