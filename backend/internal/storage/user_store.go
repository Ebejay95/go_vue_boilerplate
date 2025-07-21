package storage

import (
	"database/sql"
	"fmt"

	"backend-grpc-server/internal/database"
	"backend-grpc-server/internal/models"
)

type PostgresUserStore struct {
	db *database.DB
}

func NewPostgresUserStore(db *database.DB) UserStore {
	return &PostgresUserStore{
		db: db,
	}
}

func (s *PostgresUserStore) GetUser(id int32) (*models.User, bool) {
	query := `
		SELECT id, name, email, age, role, created_at, updated_at
		FROM users
		WHERE id = $1
	`

	user := &models.User{}
	err := s.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false
		}
		fmt.Printf("Error getting user: %v\n", err)
		return nil, false
	}

	return user, true
}

func (s *PostgresUserStore) CreateUser(params *models.CreateUserParams) (*models.User, error) {
	query := `
		INSERT INTO users (name, email, age, role)
		VALUES ($1, $2, $3, $4)
		RETURNING id, name, email, age, role, created_at, updated_at
	`

	user := &models.User{}
	err := s.db.QueryRow(query, params.Name, params.Email, params.Age, params.Role).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (s *PostgresUserStore) UpdateUser(params *models.UpdateUserParams) (*models.User, error) {
	query := `
		UPDATE users
		SET name = $2, email = $3, age = $4, role = $5, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
		RETURNING id, name, email, age, role, created_at, updated_at
	`

	user := &models.User{}
	err := s.db.QueryRow(query, params.ID, params.Name, params.Email, params.Age, params.Role).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Age,
		&user.Role,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user with ID %d not found", params.ID)
		}
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

func (s *PostgresUserStore) DeleteUser(id int32) error {
	query := `DELETE FROM users WHERE id = $1`

	result, err := s.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user with ID %d not found", id)
	}

	return nil
}

func (s *PostgresUserStore) ListUsers(params *models.ListUsersParams) ([]*models.User, int32, error) {
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
	countQuery := `SELECT COUNT(*) FROM users`
	var total int32
	err := s.db.QueryRow(countQuery).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count users: %w", err)
	}

	// Get users with pagination
	query := `
		SELECT id, name, email, age, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := s.db.Query(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.Age,
			&user.Role,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("error iterating users: %w", err)
	}

	return users, total, nil
}
