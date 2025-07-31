package storage

import (
	"testing"

	"backend-grpc-server/internal/models"
	"backend-grpc-server/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPostgresNotificationStore_CreateNotification(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := NewPostgresNotificationStore(db)

	params := &models.CreateNotificationParams{
		Message:    "Test notification",
		Type:       "info",
		Persistent: true,
	}

	notification, err := store.CreateNotification(params)

	require.NoError(t, err)
	assert.NotZero(t, notification.ID)
	assert.Equal(t, params.Message, notification.Message)
	assert.Equal(t, params.Type, notification.Type)
	assert.Equal(t, params.Persistent, notification.Persistent)
	assert.False(t, notification.Read)
	assert.NotZero(t, notification.CreatedAt)
	assert.NotZero(t, notification.UpdatedAt)
}

func TestPostgresNotificationStore_CreateNotification_WithUser(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create a user first
	userStore := NewPostgresUserStore(db)
	user, err := userStore.CreateUser(&models.CreateUserParams{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
		Role:  "user",
	})
	require.NoError(t, err)

	notificationStore := NewPostgresNotificationStore(db)

	params := &models.CreateNotificationParams{
		Message:    "User-specific notification",
		Type:       "success",
		UserID:     &user.ID,
		Persistent: true,
	}

	notification, err := notificationStore.CreateNotification(params)

	require.NoError(t, err)
	assert.NotZero(t, notification.ID)
	assert.Equal(t, params.Message, notification.Message)
	assert.Equal(t, params.Type, notification.Type)
	assert.NotNil(t, notification.UserID)
	assert.Equal(t, user.ID, *notification.UserID)
}

func TestPostgresNotificationStore_GetNotification(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := NewPostgresNotificationStore(db)

	// Create a notification first
	params := &models.CreateNotificationParams{
		Message:    "Get test notification",
		Type:       "warning",
		Persistent: true,
	}

	createdNotification, err := store.CreateNotification(params)
	require.NoError(t, err)

	// Get the notification
	notification, exists := store.GetNotification(createdNotification.ID)

	assert.True(t, exists)
	assert.Equal(t, createdNotification.ID, notification.ID)
	assert.Equal(t, createdNotification.Message, notification.Message)
	assert.Equal(t, createdNotification.Type, notification.Type)
}

func TestPostgresNotificationStore_GetNotification_NotFound(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := NewPostgresNotificationStore(db)

	notification, exists := store.GetNotification(99999)

	assert.False(t, exists)
	assert.Nil(t, notification)
}

func TestPostgresNotificationStore_MarkAsRead(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Create user and notification stores
	userStore := NewPostgresUserStore(db)
	notificationStore := NewPostgresNotificationStore(db)

	// Create a user
	user, err := userStore.CreateUser(&models.CreateUserParams{
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
		Role:  "user",
	})
	require.NoError(t, err)

	// Create a notification for the user
	notification, err := notificationStore.CreateNotification(&models.CreateNotificationParams{
		Message:    "Mark as read test",
		Type:       "info",
		UserID:     &user.ID,
		Persistent: true,
	})
	require.NoError(t, err)
	assert.False(t, notification.Read)

	// Mark as read
	err = notificationStore.MarkAsRead(notification.ID, user.ID)
	require.NoError(t, err)

	// Verify it's marked as read
	updatedNotification, exists := notificationStore.GetNotification(notification.ID)
	require.True(t, exists)
	assert.True(t, updatedNotification.Read)
}

func TestPostgresNotificationStore_ListNotifications(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := NewPostgresNotificationStore(db)

	// Create multiple notifications
	notifications := []models.CreateNotificationParams{
		{Message: "Notification 1", Type: "info", Persistent: true},
		{Message: "Notification 2", Type: "warning", Persistent: true},
		{Message: "Notification 3", Type: "error", Persistent: true},
	}

	for _, params := range notifications {
		_, err := store.CreateNotification(&params)
		require.NoError(t, err)
	}

	// Test listing notifications with pagination
	params := &models.ListNotificationsParams{
		Limit:  2,
		Offset: 0,
	}

	notificationList, total, err := store.ListNotifications(params)

	require.NoError(t, err)
	assert.Equal(t, int32(3), total)
	assert.Len(t, notificationList, 2)
}

func TestPostgresNotificationStore_GetNotificationStats(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	userStore := NewPostgresUserStore(db)
	notificationStore := NewPostgresNotificationStore(db)

	// Create a user
	user, err := userStore.CreateUser(&models.CreateUserParams{
		Name:  "Stats Test User",
		Email: "stats@example.com",
		Age:   30,
		Role:  "user",
	})
	require.NoError(t, err)

	// Create notifications with different types and read states
	notifications := []struct {
		message string
		nType   string
		read    bool
	}{
		{"Info 1", "info", false},
		{"Info 2", "info", true},
		{"Warning 1", "warning", false},
		{"Error 1", "error", false},
	}

	var createdNotifications []*models.Notification
	for _, n := range notifications {
		notification, err := notificationStore.CreateNotification(&models.CreateNotificationParams{
			Message:    n.message,
			Type:       n.nType,
			UserID:     &user.ID,
			Persistent: true,
		})
		require.NoError(t, err)
		createdNotifications = append(createdNotifications, notification)

		// Mark as read if needed
		if n.read {
			err = notificationStore.MarkAsRead(notification.ID, user.ID)
			require.NoError(t, err)
		}
	}

	// Get stats
	stats, err := notificationStore.GetNotificationStats(user.ID)
	require.NoError(t, err)

	assert.Equal(t, int32(4), stats.Total)
	assert.Equal(t, int32(3), stats.Unread)
	assert.Equal(t, int32(1), stats.Read)
	assert.Equal(t, int32(2), stats.ByType["info"])
	assert.Equal(t, int32(1), stats.ByType["warning"])
	assert.Equal(t, int32(1), stats.ByType["error"])
}
