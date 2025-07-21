package storage

import "backend-grpc-server/internal/models"

// UserStore interface - bleibt gleich
type UserStore interface {
	GetUser(id int32) (*models.User, bool)
	CreateUser(params *models.CreateUserParams) (*models.User, error)
	UpdateUser(params *models.UpdateUserParams) (*models.User, error)
	DeleteUser(id int32) error
	ListUsers(params *models.ListUsersParams) ([]*models.User, int32, error)
}

// Enhanced NotificationStore interface mit vollständigen CRUD Operations
type NotificationStore interface {
	// Basic CRUD
	GetNotification(id int32) (*models.Notification, bool)
	CreateNotification(params *models.CreateNotificationParams) (*models.Notification, error)
	UpdateNotification(params *models.UpdateNotificationParams) (*models.Notification, error)
	DeleteNotification(id int32) error

	// Advanced listing with filters
	ListNotifications(params *models.ListNotificationsParams) ([]*models.Notification, int32, error)
	ListNotificationsByUser(userID int32, params *models.ListNotificationsParams) ([]*models.Notification, int32, error)

	// Mark as read/unread functionality
	MarkAsRead(id int32, userID int32) error
	MarkAsUnread(id int32, userID int32) error
	MarkAllAsRead(userID int32) error

	// Batch operations
	DeleteReadNotifications(userID int32) error
	DeleteOldNotifications(olderThanDays int) error

	// Statistics
	GetUnreadCount(userID int32) (int32, error)
	GetNotificationStats(userID int32) (*NotificationStats, error)
}

// NotificationStats provides statistics about notifications
type NotificationStats struct {
	Total    int32 `json:"total"`
	Unread   int32 `json:"unread"`
	Read     int32 `json:"read"`
	ByType   map[string]int32 `json:"by_type"`
}
