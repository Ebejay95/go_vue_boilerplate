package storage

import (
	"database/sql"
	"fmt"
	"strings"

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

// Basic CRUD Operations

func (s *PostgresNotificationStore) GetNotification(id int32) (*models.Notification, bool) {
	query := `
		SELECT id, message, type, user_id, read, persistent, created_at, updated_at
		FROM notifications
		WHERE id = $1
	`

	notification := &models.Notification{}
	err := s.db.QueryRow(query, id).Scan(
		&notification.ID,
		&notification.Message,
		&notification.Type,
		&notification.UserID,
		&notification.Read,
		&notification.Persistent,
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
	// Only create persistent notifications in database
	if !params.Persistent {
		return nil, fmt.Errorf("non-persistent notifications should not be stored in database")
	}

	query := `
		INSERT INTO notifications (message, type, user_id, read, persistent)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, message, type, user_id, read, persistent, created_at, updated_at
	`

	notification := &models.Notification{}
	err := s.db.QueryRow(query, params.Message, params.Type, params.UserID, false, params.Persistent).Scan(
		&notification.ID,
		&notification.Message,
		&notification.Type,
		&notification.UserID,
		&notification.Read,
		&notification.Persistent,
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
		SET message = $2, type = $3, read = $4, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, message, type, user_id, read, persistent, created_at, updated_at
	`

	notification := &models.Notification{}
	err := s.db.QueryRow(query, params.ID, params.Message, params.Type, params.Read).Scan(
		&notification.ID,
		&notification.Message,
		&notification.Type,
		&notification.UserID,
		&notification.Read,
		&notification.Persistent,
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

// Advanced Listing Operations

func (s *PostgresNotificationStore) ListNotifications(params *models.ListNotificationsParams) ([]*models.Notification, int32, error) {
	// Build WHERE clause dynamically
	var conditions []string
	var args []interface{}
	argCount := 0

	if params.UserID != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", argCount))
		args = append(args, *params.UserID)
	}

	if params.Read != nil {
		argCount++
		conditions = append(conditions, fmt.Sprintf("read = $%d", argCount))
		args = append(args, *params.Read)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// Get total count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM notifications %s", whereClause)
	var total int32
	err := s.db.QueryRow(countQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count notifications: %w", err)
	}

	// Default pagination
	limit := params.Limit
	if limit <= 0 {
		limit = 50
	}
	offset := params.Offset
	if offset < 0 {
		offset = 0
	}

	// Add pagination to args
	argCount++
	limitArg := fmt.Sprintf("$%d", argCount)
	argCount++
	offsetArg := fmt.Sprintf("$%d", argCount)
	args = append(args, limit, offset)

	// Get notifications with pagination
	query := fmt.Sprintf(`
		SELECT id, message, type, user_id, read, persistent, created_at, updated_at
		FROM notifications
		%s
		ORDER BY created_at DESC
		LIMIT %s OFFSET %s
	`, whereClause, limitArg, offsetArg)

	rows, err := s.db.Query(query, args...)
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
			&notification.UserID,
			&notification.Read,
			&notification.Persistent,
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

func (s *PostgresNotificationStore) ListNotificationsByUser(userID int32, params *models.ListNotificationsParams) ([]*models.Notification, int32, error) {
	// Force userID in params
	params.UserID = &userID
	return s.ListNotifications(params)
}

// Read/Unread Operations

func (s *PostgresNotificationStore) MarkAsRead(id int32, userID int32) error {
	query := `
		UPDATE notifications
		SET read = true, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
	`

	result, err := s.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as read: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification with ID %d not found for user %d", id, userID)
	}

	return nil
}

func (s *PostgresNotificationStore) MarkAsUnread(id int32, userID int32) error {
	query := `
		UPDATE notifications
		SET read = false, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1 AND user_id = $2
	`

	result, err := s.db.Exec(query, id, userID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as unread: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("notification with ID %d not found for user %d", id, userID)
	}

	return nil
}

func (s *PostgresNotificationStore) MarkAllAsRead(userID int32) error {
	query := `
		UPDATE notifications
		SET read = true, updated_at = CURRENT_TIMESTAMP
		WHERE user_id = $1 AND read = false
	`

	_, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to mark all notifications as read: %w", err)
	}

	return nil
}

// Batch Operations

func (s *PostgresNotificationStore) DeleteReadNotifications(userID int32) error {
	query := `DELETE FROM notifications WHERE user_id = $1 AND read = true`

	_, err := s.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete read notifications: %w", err)
	}

	return nil
}

func (s *PostgresNotificationStore) DeleteOldNotifications(olderThanDays int) error {
	query := `
		DELETE FROM notifications
		WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '%d days'
	`

	_, err := s.db.Exec(fmt.Sprintf(query, olderThanDays))
	if err != nil {
		return fmt.Errorf("failed to delete old notifications: %w", err)
	}

	return nil
}

// Statistics

func (s *PostgresNotificationStore) GetUnreadCount(userID int32) (int32, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false`

	var count int32
	err := s.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread count: %w", err)
	}

	return count, nil
}

func (s *PostgresNotificationStore) GetNotificationStats(userID int32) (*NotificationStats, error) {
	// Get total counts
	totalQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1`
	unreadQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = false`
	readQuery := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND read = true`

	stats := &NotificationStats{
		ByType: make(map[string]int32),
	}

	err := s.db.QueryRow(totalQuery, userID).Scan(&stats.Total)
	if err != nil {
		return nil, fmt.Errorf("failed to get total count: %w", err)
	}

	err = s.db.QueryRow(unreadQuery, userID).Scan(&stats.Unread)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread count: %w", err)
	}

	err = s.db.QueryRow(readQuery, userID).Scan(&stats.Read)
	if err != nil {
		return nil, fmt.Errorf("failed to get read count: %w", err)
	}

	// Get counts by type
	typeQuery := `
		SELECT type, COUNT(*)
		FROM notifications
		WHERE user_id = $1
		GROUP BY type
	`

	rows, err := s.db.Query(typeQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get type counts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var notificationType string
		var count int32
		err := rows.Scan(&notificationType, &count)
		if err != nil {
			return nil, fmt.Errorf("failed to scan type count: %w", err)
		}
		stats.ByType[notificationType] = count
	}

	return stats, nil
}
