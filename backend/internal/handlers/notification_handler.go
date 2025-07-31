package handlers

import (
	"context"
	"fmt"
	"time"

	"backend-grpc-server/internal/models"
	"backend-grpc-server/internal/storage"
	"backend-grpc-server/internal/validation"
	pb "backend-grpc-server/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// NotificationHandler with socket integration and validation
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

// Socket Event Handlers with validation

func (h *NotificationHandler) handleCreateNotificationEvent(client *SocketClient, data interface{}) {
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

	// Validate input
	params := &models.CreateNotificationParams{
		Message:    message,
		Type:       notificationType,
		Persistent: persistent,
	}

	if err := validation.ValidateStruct(params); err != nil {
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": fmt.Sprintf("Validation failed: %v", err),
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
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Invalid data format",
		})
		return
	}

	notificationID, ok := dataMap["id"].(float64)
	if !ok || notificationID <= 0 {
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Invalid notification ID",
		})
		return
	}

	userID := int32(1) // TODO: Get from client.UserID
	if client.UserID != nil {
		userID = *client.UserID
	}

	err := h.store.MarkAsRead(int32(notificationID), userID)
	if err != nil {
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Failed to mark as read: " + err.Error(),
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
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Invalid data format",
		})
		return
	}

	notificationID, ok := dataMap["id"].(float64)
	if !ok || notificationID <= 0 {
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Invalid notification ID",
		})
		return
	}

	err := h.store.DeleteNotification(int32(notificationID))
	if err != nil {
		h.socketHandler.EmitToClient(client.ID, "error", map[string]interface{}{
			"message": "Failed to delete notification: " + err.Error(),
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
	// Validate input
	params := &models.CreateNotificationParams{
		Message:    message,
		Type:       notificationType,
		UserID:     targetID,
		Persistent: persistent,
	}

	if err := validation.ValidateStruct(params); err != nil {
		return fmt.Errorf("validation failed: %v", err)
	}

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

// gRPC Methods with validation

func (h *NotificationHandler) CreateNotification(ctx context.Context, req *pb.CreateNotificationRequest) (*pb.CreateNotificationResponse, error) {
	params := &models.CreateNotificationParams{
		Message:    req.Message,
		Type:       req.Type,
		UserID:     convertToInt32Pointer(req.UserId),
		Persistent: req.Persistent,
	}

	// Validate input
	if err := validation.ValidateStruct(params); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
	}

	notification, err := h.store.CreateNotification(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create notification: %v", err)
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
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "notification ID must be greater than 0")
	}

	notification, exists := h.store.GetNotification(req.Id)
	if !exists {
		return nil, status.Errorf(codes.NotFound, "notification with ID %d not found", req.Id)
	}

	return &pb.GetNotificationResponse{
		Notification: h.convertToProtoNotification(notification),
	}, nil
}

// ListNotifications returns notifications with filtering and validation
func (h *NotificationHandler) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.ListNotificationsResponse, error) {
	// Validate pagination parameters
	if req.Limit < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "limit cannot be negative")
	}
	if req.Offset < 0 {
		return nil, status.Errorf(codes.InvalidArgument, "offset cannot be negative")
	}
	if req.Limit > 1000 {
		return nil, status.Errorf(codes.InvalidArgument, "limit cannot exceed 1000")
	}

	params := &models.ListNotificationsParams{
		Limit:  req.Limit,
		Offset: req.Offset,
		UserID: convertToInt32Pointer(req.UserId),
		Read:   convertToBoolPointer(req.Read, req.HasReadFilter),
	}

	notifications, total, err := h.store.ListNotifications(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list notifications: %v", err)
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
	if req.Id <= 0 {
		return &pb.MarkNotificationAsReadResponse{
			Success: false,
			Message: "notification ID must be greater than 0",
		}, nil
	}

	if req.UserId <= 0 {
		return &pb.MarkNotificationAsReadResponse{
			Success: false,
			Message: "user ID must be greater than 0",
		}, nil
	}

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

// MarkNotificationAsUnread marks a notification as unread
func (h *NotificationHandler) MarkNotificationAsUnread(ctx context.Context, req *pb.MarkNotificationAsUnreadRequest) (*pb.MarkNotificationAsUnreadResponse, error) {
	if req.Id <= 0 {
		return &pb.MarkNotificationAsUnreadResponse{
			Success: false,
			Message: "notification ID must be greater than 0",
		}, nil
	}

	if req.UserId <= 0 {
		return &pb.MarkNotificationAsUnreadResponse{
			Success: false,
			Message: "user ID must be greater than 0",
		}, nil
	}

	err := h.store.MarkAsUnread(req.Id, req.UserId)
	if err != nil {
		return &pb.MarkNotificationAsUnreadResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Send socket update
	h.socketHandler.EmitToUser(req.UserId, "notification_updated", map[string]interface{}{
		"id":   req.Id,
		"read": false,
	})

	return &pb.MarkNotificationAsUnreadResponse{
		Success: true,
		Message: "Notification marked as unread",
	}, nil
}

// MarkAllNotificationsAsRead marks all notifications as read for a user
func (h *NotificationHandler) MarkAllNotificationsAsRead(ctx context.Context, req *pb.MarkAllNotificationsAsReadRequest) (*pb.MarkAllNotificationsAsReadResponse, error) {
	if req.UserId <= 0 {
		return &pb.MarkAllNotificationsAsReadResponse{
			Success: false,
			Message: "user ID must be greater than 0",
		}, nil
	}

	err := h.store.MarkAllAsRead(req.UserId)
	if err != nil {
		return &pb.MarkAllNotificationsAsReadResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Send socket update
	h.socketHandler.EmitToUser(req.UserId, "all_notifications_read", map[string]interface{}{
		"user_id": req.UserId,
	})

	return &pb.MarkAllNotificationsAsReadResponse{
		Success: true,
		Message: "All notifications marked as read",
	}, nil
}

// UpdateNotification updates a notification
func (h *NotificationHandler) UpdateNotification(ctx context.Context, req *pb.UpdateNotificationRequest) (*pb.UpdateNotificationResponse, error) {
	if req.Id <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "notification ID must be greater than 0")
	}

	params := &models.UpdateNotificationParams{
		ID:      req.Id,
		Message: req.Message,
		Type:    req.Type,
		Read:    req.Read,
	}

	// Validate the update parameters
	if req.Message != "" {
		testParams := &models.CreateNotificationParams{
			Message: req.Message,
			Type:    req.Type,
		}
		if err := validation.ValidateStruct(testParams); err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "validation failed: %v", err)
		}
	}

	notification, err := h.store.UpdateNotification(params)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update notification: %v", err)
	}

	// Send socket update
	h.socketHandler.EmitToAll("notification_updated", map[string]interface{}{
		"id":      notification.ID,
		"message": notification.Message,
		"type":    notification.Type,
		"read":    notification.Read,
	})

	return &pb.UpdateNotificationResponse{
		Notification: h.convertToProtoNotification(notification),
	}, nil
}

// DeleteNotification deletes a notification and sends socket update
func (h *NotificationHandler) DeleteNotification(ctx context.Context, req *pb.DeleteNotificationRequest) (*pb.DeleteNotificationResponse, error) {
	if req.Id <= 0 {
		return &pb.DeleteNotificationResponse{
			Success: false,
			Message: "notification ID must be greater than 0",
		}, nil
	}

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

// DeleteReadNotifications deletes all read notifications for a user
func (h *NotificationHandler) DeleteReadNotifications(ctx context.Context, req *pb.DeleteReadNotificationsRequest) (*pb.DeleteReadNotificationsResponse, error) {
	if req.UserId <= 0 {
		return &pb.DeleteReadNotificationsResponse{
			Success: false,
			Message: "user ID must be greater than 0",
		}, nil
	}

	err := h.store.DeleteReadNotifications(req.UserId)
	if err != nil {
		return &pb.DeleteReadNotificationsResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	// Send socket update
	h.socketHandler.EmitToUser(req.UserId, "read_notifications_deleted", map[string]interface{}{
		"user_id": req.UserId,
	})

	return &pb.DeleteReadNotificationsResponse{
		Success: true,
		Message: "Read notifications deleted successfully",
	}, nil
}

// GetNotificationStats returns notification statistics
func (h *NotificationHandler) GetNotificationStats(ctx context.Context, req *pb.GetNotificationStatsRequest) (*pb.GetNotificationStatsResponse, error) {
	if req.UserId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user ID must be greater than 0")
	}

	stats, err := h.store.GetNotificationStats(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get notification stats: %v", err)
	}

	return &pb.GetNotificationStatsResponse{
		Total:  stats.Total,
		Unread: stats.Unread,
		Read:   stats.Read,
		ByType: stats.ByType,
	}, nil
}

// SendRealtimeNotification sends a real-time only notification
func (h *NotificationHandler) SendRealtimeNotification(ctx context.Context, req *pb.SendRealtimeNotificationRequest) (*pb.SendRealtimeNotificationResponse, error) {
	// Validate input
	params := &models.CreateNotificationParams{
		Message:    req.Message,
		Type:       req.Type,
		UserID:     convertToInt32Pointer(req.UserId),
		Persistent: false, // Real-time notifications are never persistent
	}

	if err := validation.ValidateStruct(params); err != nil {
		return &pb.SendRealtimeNotificationResponse{
			Success: false,
			Message: fmt.Sprintf("Validation failed: %v", err),
		}, nil
	}

	// Convert data map
	data := make(map[string]interface{})
	for k, v := range req.Data {
		data[k] = v
	}

	// Send the notification
	targetType := "all"
	var targetID *int32
	if req.UserId != 0 {
		targetType = "user"
		targetID = &req.UserId
	}

	err := h.SendNotification(req.Message, req.Type, targetType, targetID, false, data)
	if err != nil {
		return &pb.SendRealtimeNotificationResponse{
			Success: false,
			Message: err.Error(),
		}, nil
	}

	return &pb.SendRealtimeNotificationResponse{
		Success: true,
		Message: "Realtime notification sent successfully",
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
