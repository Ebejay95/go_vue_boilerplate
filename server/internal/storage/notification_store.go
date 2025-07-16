package storage

import (
	"database/sql"
	"fmt"

	"backend-grpc-server/internal/database"
	"backend-grpc-server/internal/models"
)

type PostgresNotificationStore struct {
	db *database.DB
}

func NewPostgresNotificationStore(db *database.DB) NotificationStore {
	return &PostgresNotificationStore{
		db: db,
	}
}

func (s *PostgresNotificationStore) GetNotification(id int32) (*models.Notification, bool) {
	query := `
		SELECT id, message, type, created_at, updated_at
		FROM notifications
		WHERE id = $1
	`

	notification := &models.Notification{}
	err := s.db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.Message,
		&notification.Type,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		fmt.Printf("Error getting notification: %v\n", err)
		return nil, false
	}

	return notification, true
}

func (s *PostgresNotificationStore) CreateNotification(params *models.CreateNotificationParams) (*models.Notification, error) {
	query := `
		INSERT INTO notifications (message, type)
		VALUES ($1, $2)
		RETURNING id, message, type, created_at, updated_at
	`

	notification := &models.Notification{}
	err := s.db.QueryRow(query, params.Message, params.Type).Scan(
		&notification.ID,
		&notification.Message,
		&notification.Type,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	return notification, nil
}

func (s *PostgresNotificationStore) UpdateNotification(params *models.UpdateNotificationParams) (*models.Notification, error) {
	query := `
		UPDATE notifications
		SET message = $2, type = $3, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, message, type, created_at, updated_at
	`

	notification := &models.Notification{}
	err := s.db.QueryRow(query, params.ID, params.Message, params.Type).Scan(
		&notification.ID,
		&notification.Message,
		&notification.Type,
		&notification.CreatedAt,
		&notification.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("notification with ID %d not found", params.ID)
		}
		return nil, fmt.Errorf("failed to update notification: %w", err)
	}

	return notification, nil
}

func (s *PostgresNotificationStore) DeleteNotification(id int32) error {
	query := `DELETE FROM notifications WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification with ID %d not found", id)
	}

	return nil
}

func (s *PostgresNotificationStore) ListNotifications(params *models.ListNotificationsParams) ([]*models.Notification, int32, error) {
	// Default values
	limit := params.Limit
	if limit <= 0 {
		limit = 50 // Default limit
	}
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	// Get total count
	countQuery := `SELECT COUNT(*) FROM notifications`
	var total int32
	err := s.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	// Get notifications with pagination
	query := `
		SELECT id, message, type, created_at, updated_at
		FROM notifications
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list notifications: %w", err)
	}
	defer rows.Close()

	var notifications []*models.Notification
	for rows.Next() {
		notification := &models.Notification{}
		err := rows.Scan(
			&notification.ID,
			&notification.Message,
			&notification.Type,
			&notification.CreatedAt,
			&notification.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan notification: %w", err)
		}
		notifications = append(notifications, notification)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating notifications: %w", err)
	}

	return notifications, total, nil
}
