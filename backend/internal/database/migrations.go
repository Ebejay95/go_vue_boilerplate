// internal/database/migrations.go
package database

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// Migration represents a database migration
type Migration struct {
	Version     string
	Description string
	Up          string
	Down        string
	FilePath    string
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db *DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *DB) *MigrationManager {
	return &MigrationManager{db: db}
}

// createMigrationsTable creates the migrations tracking table
func (m *MigrationManager) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`
	_, err := m.db.Exec(query)
	return err
}

// isApplied checks if a migration has been applied
func (m *MigrationManager) isApplied(version string) (bool, error) {
	var count int
	err := m.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = $1", version).Scan(&count)
	return count > 0, err
}

// markApplied marks a migration as applied
func (m *MigrationManager) markApplied(version string) error {
	_, err := m.db.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version)
	return err
}

// markUnapplied removes a migration from applied list
func (m *MigrationManager) markUnapplied(version string) error {
	_, err := m.db.Exec("DELETE FROM schema_migrations WHERE version = $1", version)
	return err
}

// loadMigrationsFromFiles loads migrations from .sql files in the filesystem
func (m *MigrationManager) loadMigrationsFromFiles() ([]Migration, error) {
	migrationsDir := "./internal/database/migrations"

	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrations []Migration

	for _, file := range files {
		// Only accept .sql files
		if !strings.HasSuffix(file.Name(), ".sql") {
			log.Printf("Skipping non-SQL file: %s", file.Name())
			continue
		}

		// Parse filename: 2507211130_setup.sql -> version: 2507211130, description: setup
		name := strings.TrimSuffix(file.Name(), ".sql")
		parts := strings.SplitN(name, "_", 2)
		if len(parts) != 2 {
			log.Printf("Skipping migration file with invalid name format: %s (expected: YYYYMMDDHHNN_description)", file.Name())
			continue
		}

		version := parts[0]
		description := strings.ReplaceAll(parts[1], "_", " ")

		filePath := filepath.Join(migrationsDir, file.Name())
		content, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, fmt.Errorf("failed to read migration file %s: %w", file.Name(), err)
		}

		migration := Migration{
			Version:     version,
			Description: description,
			Up:          string(content),
			Down:        "", // Rollback not supported for file-based migrations
			FilePath:    filePath,
		}

		migrations = append(migrations, migration)
		log.Printf("Loaded migration: %s (%s)", version, description)
	}

	return migrations, nil
}

// RunMigrations runs all pending migrations from files
func (m *MigrationManager) RunMigrations() error {
	log.Println("Starting file-based migrations...")

	// Create migrations table if it doesn't exist
	if err := m.createMigrationsTable(); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	migrations, err := m.loadMigrationsFromFiles()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	if len(migrations) == 0 {
		log.Println("No migration files found in ./internal/database/migrations/")
		return nil
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	log.Printf("Found %d migration files", len(migrations))

	for _, migration := range migrations {
		applied, err := m.isApplied(migration.Version)
		if err != nil {
			return fmt.Errorf("failed to check if migration %s is applied: %w", migration.Version, err)
		}

		if applied {
			log.Printf("Migration %s (%s) already applied, skipping", migration.Version, migration.Description)
			continue
		}

		log.Printf("Applying migration %s: %s", migration.Version, migration.Description)

		// Start transaction
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("failed to start transaction for migration %s: %w", migration.Version, err)
		}

		// Execute migration as a single script
		if err := m.executeMigrationScript(tx, migration); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
		}

		// Mark as applied
		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", migration.Version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to mark migration %s as applied: %w", migration.Version, err)
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", migration.Version, err)
		}

		log.Printf("Successfully applied migration %s", migration.Version)
	}

	log.Println("All file-based migrations completed successfully")
	return nil
}

// executeMigrationScript executes the entire migration file as a single script
func (m *MigrationManager) executeMigrationScript(tx *sql.Tx, migration Migration) error {
	// Clean up the SQL content
	sqlContent := strings.TrimSpace(migration.Up)

	// Remove comments that start with -- (but keep SQL inside functions)
	lines := strings.Split(sqlContent, "\n")
	var cleanLines []string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// Skip empty lines and standalone comments
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}
		cleanLines = append(cleanLines, line)
	}

	if len(cleanLines) == 0 {
		log.Printf("No SQL content found in migration %s", migration.Version)
		return nil
	}

	// Join the lines back together
	cleanSQL := strings.Join(cleanLines, "\n")

	log.Printf("Executing migration script: %s...", truncateString(cleanSQL, 100))

	// Execute the entire script as one statement
	if _, err := tx.Exec(cleanSQL); err != nil {
		return fmt.Errorf("failed to execute migration script: %w\nSQL Content: %s", err, cleanSQL)
	}

	return nil
}

// truncateString truncates a string to maxLen characters
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Rollback rolls back the last migration (limited support for file-based migrations)
func (m *MigrationManager) Rollback() error {
	log.Println("Rollback not supported for file-based migrations")
	log.Println("To rollback, create a new migration file with reverse operations")
	return fmt.Errorf("rollback not implemented for file-based migrations")
}

// ListAppliedMigrations shows which migrations have been applied
func (m *MigrationManager) ListAppliedMigrations() error {
	rows, err := m.db.Query("SELECT version, applied_at FROM schema_migrations ORDER BY version")
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	log.Println("Applied migrations:")
	log.Println("Version\t\tApplied At")
	log.Println("-------\t\t----------")

	for rows.Next() {
		var version, appliedAt string
		if err := rows.Scan(&version, &appliedAt); err != nil {
			return fmt.Errorf("failed to scan migration row: %w", err)
		}
		log.Printf("%s\t\t%s", version, appliedAt)
	}

	return rows.Err()
}
