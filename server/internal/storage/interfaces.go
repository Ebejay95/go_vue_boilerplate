package storage

import "backend-grpc-server/internal/models"

// UserStore interface defines the contract for user storage operations
type UserStore interface {
	GetUser(id int32) (*models.User, bool)
	CreateUser(params *models.CreateUserParams) (*models.User, error)
	UpdateUser(params *models.UpdateUserParams) (*models.User, error)
	DeleteUser(id int32) error
	ListUsers(params *models.ListUsersParams) ([]*models.User, int32, error)
}

// NotificationStore interface defines the contract for notification storage operations
type NotificationStore interface {
	GetNotification(id int32) (*models.Notification, bool)
	CreateNotification(params *models.CreateNotificationParams) (*models.Notification, error)
	UpdateNotification(params *models.UpdateNotificationParams) (*models.Notification, error)
	DeleteNotification(id int32) error
	ListNotifications(params *models.ListNotificationsParams) ([]*models.Notification, int32, error)
}
