package handlers

import (
	"context"
	"testing"

	"backend-grpc-server/internal/storage"
	"backend-grpc-server/internal/testutil"
	pb "backend-grpc-server/pb"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNotificationHandler_CreateNotification(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := storage.NewPostgresNotificationStore(db)
	socketHandler := NewSocketHandler()
	handler := NewNotificationHandler(store, socketHandler)

	req := &pb.CreateNotificationRequest{
		Message:    "Test notification",
		Type:       "info",
		Persistent: true,
	}

	resp, err := handler.CreateNotification(context.Background(), req)

	require.NoError(t, err)
	assert.NotNil(t, resp.Notification)
	assert.NotZero(t, resp.Notification.Id)
	assert.Equal(t, req.Message, resp.Notification.Message)
	assert.Equal(t, req.Type, resp.Notification.Type)
	assert.Equal(t, req.Persistent, resp.Notification.Persistent)
}

func TestNotificationHandler_CreateNotification_ValidationError(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := storage.NewPostgresNotificationStore(db)
	socketHandler := NewSocketHandler()
	handler := NewNotificationHandler(store, socketHandler)

	tests := []struct {
		name string
		req  *pb.CreateNotificationRequest
	}{
		{
			name: "empty message",
			req: &pb.CreateNotificationRequest{
				Message:    "",
				Type:       "info",
				Persistent: true,
			},
		},
		{
			name: "invalid type",
			req: &pb.CreateNotificationRequest{
				Message:    "Test message",
				Type:       "invalid_type",
				Persistent: true,
			},
		},
		{
			name: "message too long",
			req: &pb.CreateNotificationRequest{
				Message:    string(make([]byte, 1001)), // Over 1000 chars
				Type:       "info",
				Persistent: true,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.CreateNotification(context.Background(), tt.req)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), "validation failed")
		})
	}
}

func TestNotificationHandler_GetNotification(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := storage.NewPostgresNotificationStore(db)
	socketHandler := NewSocketHandler()
	handler := NewNotificationHandler(store, socketHandler)

	// Create a notification first
	createReq := &pb.CreateNotificationRequest{
		Message:    "Get test notification",
		Type:       "warning",
		Persistent: true,
	}

	createResp, err := handler.CreateNotification(context.Background(), createReq)
	require.NoError(t, err)

	// Get the notification
	getReq := &pb.GetNotificationRequest{
		Id: createResp.Notification.Id,
	}

	getResp, err := handler.GetNotification(context.Background(), getReq)

	require.NoError(t, err)
	assert.Equal(t, createResp.Notification.Id, getResp.Notification.Id)
	assert.Equal(t, createResp.Notification.Message, getResp.Notification.Message)
	assert.Equal(t, createResp.Notification.Type, getResp.Notification.Type)
}

func TestNotificationHandler_GetNotification_InvalidID(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := storage.NewPostgresNotificationStore(db)
	socketHandler := NewSocketHandler()
	handler := NewNotificationHandler(store, socketHandler)

	req := &pb.GetNotificationRequest{
		Id: 0, // Invalid ID
	}

	_, err := handler.GetNotification(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "must be greater than 0")
}

func TestNotificationHandler_ListNotifications_ValidationError(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := storage.NewPostgresNotificationStore(db)
	socketHandler := NewSocketHandler()
	handler := NewNotificationHandler(store, socketHandler)

	tests := []struct {
		name string
		req  *pb.ListNotificationsRequest
	}{
		{
			name: "negative limit",
			req: &pb.ListNotificationsRequest{
				Limit:  -1,
				Offset: 0,
			},
		},
		{
			name: "negative offset",
			req: &pb.ListNotificationsRequest{
				Limit:  10,
				Offset: -1,
			},
		},
		{
			name: "limit too large",
			req: &pb.ListNotificationsRequest{
				Limit:  1001,
				Offset: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := handler.ListNotifications(context.Background(), tt.req)
			assert.Error(t, err)
		})
	}
}

func TestNotificationHandler_SendNotification(t *testing.T) {
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	store := storage.NewPostgresNotificationStore(db)
	socketHandler := NewSocketHandler()
	handler := NewNotificationHandler(store, socketHandler)

	tests := []struct {
		name        string
		message     string
		nType       string
		targetType  string
		targetID    *int32
		persistent  bool
		expectError bool
	}{
		{
			name:        "valid global notification",
			message:     "Test global notification",
			nType:       "info",
			targetType:  "all",
			targetID:    nil,
			persistent:  true,
			expectError: false,
		},
		{
			name:        "valid user notification",
			message:     "Test user notification",
			nType:       "success",
			targetType:  "user",
			targetID:    func() *int32 { id := int32(1); return &id }(),
			persistent:  false,
			expectError: false,
		},
		{
			name:        "invalid message",
			message:     "",
			nType:       "info",
			targetType:  "all",
			targetID:    nil,
			persistent:  true,
			expectError: true,
		},
		{
			name:        "invalid type",
			message:     "Test message",
			nType:       "invalid_type",
			targetType:  "all",
			targetID:    nil,
			persistent:  true,
			expectError: true,
		},
		{
			name:        "invalid target type",
			message:     "Test message",
			nType:       "info",
			targetType:  "invalid_target",
			targetID:    nil,
			persistent:  true,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := handler.SendNotification(
				tt.message,
				tt.nType,
				tt.targetType,
				tt.targetID,
				tt.persistent,
				nil,
			)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
