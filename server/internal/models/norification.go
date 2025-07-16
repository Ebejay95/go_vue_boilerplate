package models

import "time"

type Notification struct {
	ID        int32     `json:"id" db:"id"`
	Message   string    `json:"message" db:"message"`
	Type      string    `json:"type" db:"type"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateNotificationParams struct {
	Message string `json:"message"`
	Type    string `json:"type"`
}

type UpdateNotificationParams struct {
	ID      int32  `json:"id"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type ListNotificationsParams struct {
	Limit  int32 `json:"limit"`
	Offset int32 `json:"offset"`
}
