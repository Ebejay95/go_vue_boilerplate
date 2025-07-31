package testutil

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"backend-grpc-server/internal/database"
	"backend-grpc-server/internal/models"
	_ "github.com/lib/pq"
)

// findMigrationsDir finds the migrations directory from the current working directory
func findMigrationsDir() string {
	// Try different possible paths
	paths := []string{
		"./internal/database/migrations",         // From backend root
		"../database/migrations",                 // From testutil directory
		"../../database/migrations",              // From deeper nested tests
		"./backend/internal/database/migrations", // From project root
	}

	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			absPath, _ := filepath.Abs(path)
			log.Printf("Found migrations directory at: %s", absPath)
			return path
		}
	}

	// If not found, try to find it relative to this file
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)
	migrationPath := filepath.Join(dir, "..", "database", "migrations")
	if _, err := os.Stat(migrationPath); err == nil {
		return migrationPath
	}

	log.Printf("Warning: Could not find migrations directory. Tried paths: %v", paths)
	return "./internal/database/migrations" // fallback
}

// SetupTestDB creates a test database connection and runs migrations
func SetupTestDB(t *testing.T) *database.DB {
	// Use environment variables or default test database
	db, err := database.NewConnectionWithoutMigrations()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Get the correct migrations path
	migrationsPath := findMigrationsDir()
	log.Printf("Using migrations path: %s", migrationsPath)

	// Run migrations for test with correct path
	if err := db.RunMigrationsFromPath(migrationsPath); err != nil {
		t.Fatalf("Failed to run migrations from %s: %v", migrationsPath, err)
	}

	return db
}

// SetupTestDBWithPath sets up test DB with explicit migrations path
func SetupTestDBWithPath(t *testing.T, migrationsPath string) *database.DB {
	db, err := database.NewConnectionWithoutMigrations()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	if err := db.RunMigrationsFromPath(migrationsPath); err != nil {
		t.Fatalf("Failed to run migrations from %s: %v", migrationsPath, err)
	}

	return db
}

// SetupTestDBSimple creates a test database connection without migrations
func SetupTestDBSimple(t *testing.T) *database.DB {
	db, err := database.NewConnectionWithoutMigrations()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}
	return db
}

// CleanupTestDB cleans up test database
func CleanupTestDB(t *testing.T, db *database.DB) {
	// Clean up tables in reverse order due to foreign keys
	tables := []string{"notifications", "users"}
	for _, table := range tables {
		_, err := db.Exec("TRUNCATE TABLE " + table + " CASCADE")
		if err != nil {
			// If TRUNCATE fails, try DELETE
			_, err = db.Exec("DELETE FROM " + table)
			if err != nil {
				log.Printf("Failed to clean table %s: %v", table, err)
			}
		}
	}
	db.Close()
}

// CreateTestUser creates a test user
func CreateTestUser(t *testing.T) *models.User {
	return &models.User{
		ID:    1,
		Name:  "Test User",
		Email: "test@example.com",
		Age:   25,
		Role:  "user",
	}
}

// CreateTestNotification creates a test notification
func CreateTestNotification(t *testing.T) *models.Notification {
	return &models.Notification{
		ID:         1,
		Message:    "Test notification",
		Type:       "info",
		UserID:     nil,
		Read:       false,
		Persistent: true,
	}
}

// SetupMinimalDB creates a test database with only essential tables for unit tests
func SetupMinimalDB(t *testing.T) *database.DB {
	db, err := database.NewConnectionWithoutMigrations()
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Create minimal tables for testing
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			age INTEGER NOT NULL CHECK (age >= 0),
			role VARCHAR(100) NOT NULL DEFAULT 'user',
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS notifications (
			id SERIAL PRIMARY KEY,
			message TEXT NOT NULL,
			type VARCHAR(50) NOT NULL DEFAULT 'info',
			user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
			read BOOLEAN DEFAULT FALSE,
			persistent BOOLEAN DEFAULT TRUE,
			created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			t.Fatalf("Failed to create test tables: %v", err)
		}
	}

	return db
}
