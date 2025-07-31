// internal/database/connection.go
package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewConnection() (*DB, error) {
	// Get database configuration from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "grpc_server_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL database: %s", dbname)

	// Create DB wrapper
	dbWrapper := &DB{db}

	// Run migrations automatically if enabled
	if getEnv("AUTO_MIGRATE", "true") == "true" {
		log.Println("Running database migrations...")
		migrationManager := NewMigrationManager(dbWrapper)
		if err := migrationManager.RunMigrations(); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}
	}

	return dbWrapper, nil
}

func NewConnectionWithoutMigrations() (*DB, error) {
	// Get database configuration from environment variables
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "postgres")
	dbname := getEnv("DB_NAME", "grpc_server_db")
	sslmode := getEnv("DB_SSLMODE", "disable")

	// Build connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Successfully connected to PostgreSQL database: %s", dbname)

	return &DB{db}, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func (db *DB) Close() error {
	return db.DB.Close()
}

// Health check method
func (db *DB) HealthCheck() error {
	return db.Ping()
}

// Migration helper methods
func (db *DB) RunMigrations() error {
	migrationManager := NewMigrationManager(db)
	return migrationManager.RunMigrations()
}

// Add the missing method for tests
func (db *DB) RunMigrationsFromPath(migrationsPath string) error {
	// Create a temporary migration manager with custom path
	migrationManager := &MigrationManager{db: db}

	// Load migrations from the specified path
	migrations, err := migrationManager.loadMigrationsFromPath(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to load migrations from %s: %w", migrationsPath, err)
	}

	// Create migrations table if it doesn't exist
	if err := migrationManager.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Run migrations
	return migrationManager.runMigrationsFromList(migrations)
}

func (db *DB) RollbackLastMigration() error {
	migrationManager := NewMigrationManager(db)
	return migrationManager.Rollback()
}
