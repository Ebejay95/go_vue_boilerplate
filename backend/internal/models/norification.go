package models

import (
	"github.com/go-playground/validator/v10"
	"time"
)

// Persistent Notification - stored in database
type Notification struct {
	ID         int32     `json:"id" db:"id"`
	Message    string    `json:"message" db:"message" validate:"required,min=1,max=1000"`
	Type       string    `json:"type" db:"type" validate:"required,oneof=info warning error success"`
	UserID     *int32    `json:"user_id" db:"user_id" validate:"omitempty,min=1"`
	Read       bool      `json:"read" db:"read"`
	Persistent bool      `json:"persistent" db:"persistent"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

// Real-time Notification - only sent via WebSocket, not stored
type RealtimeNotification struct {
	ID        string                 `json:"id"` // UUID for tracking
	Message   string                 `json:"message"`
	Type      string                 `json:"type"` // "info", "warning", "error", "success"
	UserID    *int32                 `json:"user_id,omitempty"` // Optional: for targeted notifications
	Data      map[string]interface{} `json:"data,omitempty"` // Additional payload
	Timestamp time.Time              `json:"timestamp"`
}

// CRUD Parameters for persistent notifications
type CreateNotificationParams struct {
	Message    string `json:"message" validate:"required,min=1,max=1000"`
	Type       string `json:"type" validate:"required,oneof=info warning error success"`
	UserID     *int32 `json:"user_id,omitempty" validate:"omitempty,min=1"`
	Persistent bool   `json:"persistent"`
}

type UpdateNotificationParams struct {
	ID      int32  `json:"id"`
	Message string `json:"message"`
	Type    string `json:"type"`
	Read    bool   `json:"read"`
}

type ListNotificationsParams struct {
	Limit  int32  `json:"limit"`
	Offset int32  `json:"offset"`
	UserID *int32 `json:"user_id,omitempty"` // Filter by user
	Read   *bool  `json:"read,omitempty"`    // Filter by read status
}

type MarkNotificationReadParams struct {
	ID     int32 `json:"id"`
	UserID int32 `json:"user_id"` // Ensure user can only mark their own notifications
}

// Real-time notification creation parameters
type CreateRealtimeNotificationParams struct {
	Message string                 `json:"message"`
	Type    string                 `json:"type"`
	UserID  *int32                 `json:"user_id,omitempty"`
	Data    map[string]interface{} `json:"data,omitempty"`
}

func (n *CreateNotificationParams) Validate() error {
	validate := validator.New()
	return validate.Struct(n)
}
