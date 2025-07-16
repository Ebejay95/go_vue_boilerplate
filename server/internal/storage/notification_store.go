package storage

import (
	"sync"
	"backend-grpc-server/internal/models"
)

// NotificationStore interface defines the contract for notification storage operations
type NotificationStore interface {
	GetNotification(id int32) (*models.Notification, bool)
	CreateNotification(notification *models.Notification) *models.Notification
	ListNotifications() []*models.Notification
}

// InMemoryNotificationStore implements NotificationStore using in-memory storage
type InMemoryNotificationStore struct {
	notifications map[int32]*models.Notification
	mutex         sync.RWMutex
	nextID        int32
}

// NewInMemoryNotificationStore creates a new in-memory notification store
func NewInMemoryNotificationStore() *InMemoryNotificationStore {
	return &InMemoryNotificationStore{
		notifications: make(map[int32]*models.Notification),
		nextID:        1,
	}
}

// GetNotification retrieves a notification by ID
func (s *InMemoryNotificationStore) GetNotification(id int32) (*models.Notification, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	notification, exists := s.notifications[id]
	return notification, exists
}

// CreateNotification creates a new notification
func (s *InMemoryNotificationStore) CreateNotification(notification *models.Notification) *models.Notification {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	notification.ID = s.nextID
	s.notifications[s.nextID] = notification
	s.nextID++

	return notification
}

// ListNotifications returns all notifications
func (s *InMemoryNotificationStore) ListNotifications() []*models.Notification {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	notifications := make([]*models.Notification, 0, len(s.notifications))
	for _, notification := range s.notifications {
		notifications = append(notifications, notification)
	}

	return notifications
}